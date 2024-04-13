package tunnel

import (
	"ssh-compose/internal/compose"
	"ssh-compose/internal/env"
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
