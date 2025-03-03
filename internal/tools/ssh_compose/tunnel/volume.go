package tunnel

import (
	"path/filepath"
	"vpdb-dev-tool/internal/lib/xstrings"
	"vpdb-dev-tool/internal/tools/ssh_compose/compose"
	"vpdb-dev-tool/internal/tools/ssh_compose/env"
)

func makeVolumes(sshHome string) []compose.Volume {
	return []compose.Volume{
		{
			Type:   "bind",
			Source: injectionRef[env.SSHSockSrc],
			Target: injectionRef[env.SSHSockTgt],
		},
		{
			Type:   "bind",
			Source: filepath.Join(xstrings.Coalesce(sshHome, "$HOME/.ssh"), "known_hosts"),
			Target: "/root/.ssh/known_hosts",
		},
	}
}
