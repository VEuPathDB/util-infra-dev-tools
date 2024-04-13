package tunnel

import (
	"vpdb-dev-tool/internal/tools/ssh_compose/compose"
	"vpdb-dev-tool/internal/tools/ssh_compose/env"
)

func makeVolumes() []compose.Volume {
	return []compose.Volume{
		{
			Type:   "bind",
			Source: injectionRef[env.SSHSockSrc],
			Target: injectionRef[env.SSHSockTgt],
		},
		{
			Type:   "bind",
			Source: "$HOME/.ssh/known_hosts",
			Target: "/root/.ssh/known_hosts",
		},
	}
}
