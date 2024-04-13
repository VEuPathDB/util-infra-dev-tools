package tunnel

import (
	"strings"

	"ssh-compose/internal/compose"
	"ssh-compose/internal/env"
	"ssh-compose/internal/hosts"
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
			"internal": {
				Aliases: []string{reqEnvName},
			},
		},
	}
}
