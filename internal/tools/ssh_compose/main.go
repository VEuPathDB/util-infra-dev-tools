package ssh_compose

import (
	"fmt"
	"log"

	"vpdb-dev-tool/internal/tools/ssh_compose/hosts"
	"vpdb-dev-tool/internal/tools/ssh_compose/project"
	"vpdb-dev-tool/internal/tools/ssh_compose/tunnel"
)

const followupSteps = "\n" +
	"Next steps:\n" +
	"  * Edit the project's `.env` file to fill in any newly generated variables.\n" +
	"  * Add `-f " + project.OutputComposeFileName + "` to any docker compose" +
	" up/down commands being used to include the newly generated ssh tunnel" +
	" containers.\n" +
	"  * Ensure that an ssh-agent instance is running which will provide the" +
	" value for the SSH_AUTH_SOCK environment variable used by the generated" +
	" compose file."

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

	fmt.Println(followupSteps)
}
