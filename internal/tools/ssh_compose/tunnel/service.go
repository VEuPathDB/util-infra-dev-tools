package tunnel

import (
	"strings"

	"vpdb-dev-tool/internal/tools/ssh_compose/compose"
	"vpdb-dev-tool/internal/tools/ssh_compose/env"
	"vpdb-dev-tool/internal/tools/ssh_compose/hosts"
)

func buildService(name, image, hostString string, host hosts.Host, volumes []compose.Volume) compose.Service {
	envName := strings.ToUpper(name)
	reqEnvName := envToRequiredInjectVar(envName)

	return compose.Service{
		Image:      image,
		Entrypoint: makeSSHCommand(host, reqEnvName, hostString),
		Volumes:    volumes,
		Environment: map[string]string{
			env.SSHSockDef: injectionRef[env.SSHSockTgt],
		},
		Networks: map[string]compose.Network{
		        "default": {},
			"internal": {
				Aliases: []string{reqEnvName},
			},
		},
	}
}
