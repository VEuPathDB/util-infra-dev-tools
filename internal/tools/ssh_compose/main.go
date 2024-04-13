package ssh_compose

import (
	"log"

	"vpdb-dev-tool/internal/tools/ssh_compose/hosts"
	"vpdb-dev-tool/internal/tools/ssh_compose/project"
	"vpdb-dev-tool/internal/tools/ssh_compose/tunnel"
)

func main(hostsFile, image string) {
	hostList := hosts.ReadHostsFile(hostsFile)
	if len(hostList) == 0 {
		log.Fatalln("input contained no valid host entries")
	}

	tunnelConfigs := tunnel.BuildTunnelConfigs(tunnel.Config{
		ComposeVersion: "3.5",
		DockerImage:    image,
		Entries:        hostList,
	})

	project.WriteOutConfigs(tunnelConfigs, hostsFile)
}
