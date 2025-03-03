package tunnel

import (
	"strings"

	"vpdb-dev-tool/internal/tools/ssh_compose/compose"
	"vpdb-dev-tool/internal/tools/ssh_compose/hosts"
)

type Config struct {
	DockerImage string
	Entries     map[hosts.Host][]string
	SSHHome     string
}

type BuiltConfigs struct {
	Compose compose.Config
	Hosts   map[string]string
}

// BuildTunnelConfigs assembles a docker-compose configuration based on the
// given configuration.
//
// The entries in the given configuration will be converted into service entries
// in the built docker-compose configuration.
func BuildTunnelConfigs(config Config) BuiltConfigs {
	// Make the user@host reference value that will be used in all the generated
	// SSH commands.
	hostString := makeHostRef()

	// Create the bind volume list that will be used in all the SSH service
	// containers.
	volumes := makeVolumes(config.SSHHome)

	// Make a new map to hold all the docker-compose service entries.
	services := make(map[string]compose.Service, countServiceEntries(config.Entries))

	// Make a new map for env-var to host address pairs.
	hosts := make(map[string]string, len(config.Entries))

	// Make a new map to hold a mapping of dependent name to list of dependency
	// docker-compose service names.
	//
	// Example:
	// {
	//   rest-service: [ s3_some_host_9000, db_some_host_5432 ]
	// }
	//
	// This will later be translated into docker-compose service entries like:
	//
	//   rest-service:
	//     depends_on:
	//     - s3_some_host_9000
	//     - db_some_host_5432
	serviceDependents := make(map[string][]string, len(config.Entries))

	// Build the service entries for the new SSH container services.
	for host, dependents := range config.Entries {
		serviceName := makeServiceName(host)
		services[serviceName] = buildService(serviceName, config.DockerImage, hostString, host, volumes)

		hosts[strings.ToUpper(serviceName)] = host.Address

		// For each of the dependents for the new SSH service, add an entry to the
		// serviceDependents map so that we can use the new service name instead of
		// the host address.
		for _, dependent := range dependents {
			serviceDependents[dependent] = append(serviceDependents[dependent], serviceName)
		}
	}

	// Now that we've built the SSH containers, we can add service entries to link
	// the dependents to the new service entries.
	for dependent, dependencies := range serviceDependents {
		structs := make(map[string]compose.Dependency, len(dependencies))
		for _, name := range dependencies {
			structs[name] = compose.Dependency{Condition: compose.DependencyConditionHealthy}
		}

		services[dependent] = compose.Service{DependsOn: structs}
	}

	return BuiltConfigs{
		Compose: compose.Config{
			Services: services,
		},
		Hosts: hosts,
	}
}

func countServiceEntries(entries map[hosts.Host][]string) int {
	out := len(entries)

	for _, v := range entries {
		out += len(v)
	}

	return out
}
