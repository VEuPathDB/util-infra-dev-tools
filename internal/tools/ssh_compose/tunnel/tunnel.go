package tunnel

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"vpdb-dev-tool/internal/tools/ssh_compose/compose"
	"vpdb-dev-tool/internal/tools/ssh_compose/env"
	"vpdb-dev-tool/internal/tools/ssh_compose/hosts"
)

const (
	prefixTunnel = "tunnel_"
)

var (
	rgxHostToName = regexp.MustCompile(`\W+`)

	injectionRef = make(map[string]string, 6)
)

func init() {
	injectionRef[env.TunnelUser] = envToRequiredInjectVar(env.TunnelUser)
	injectionRef[env.TunnelHost] = envToRequiredInjectVar(env.TunnelHost)
	injectionRef[env.TunnelPort] = envToRequiredInjectVar(env.TunnelPort)
	injectionRef[env.SSHSockSrc] = envToOptionalInjectVar(env.SSHSockSrc, "$"+env.SSHSockDef)
	injectionRef[env.SSHSockTgt] = envToOptionalInjectVar(env.SSHSockTgt, "$"+env.SSHSockDef)
}

func makeServiceName(host hosts.Host) string {
	return prefixTunnel + strings.ToLower(rgxHostToName.ReplaceAllString(host.Address, "_")) + "_" + strconv.Itoa(int(host.Port))
}

func envToRequiredInjectVar(envKey string) string {
	return "${" + envKey + ":?}"
}

func envToOptionalInjectVar(envKey, alternative string) string {
	return fmt.Sprintf("${%s:-%s}", envKey, alternative)
}

func makeTunnelDef(host hosts.Host, reqEnvVar string) string {
	return fmt.Sprintf("%s:%d:%s:%d", reqEnvVar, host.Port, reqEnvVar, host.Port)
}

func makeHostRef() string {
	return fmt.Sprintf("%s@%s", injectionRef[env.TunnelUser], injectionRef[env.TunnelHost])
}

func makeSSHCommand(host hosts.Host, reqEnvVar, hostString string) string {
	return fmt.Sprintf(
		"ssh -tNn -p %s -o ServerAliveInterval=60 -L %s:%d:%s:%d %s",
		injectionRef[env.TunnelPort],
		reqEnvVar,
		host.Port,
		reqEnvVar,
		host.Port,
		hostString,
	)
}

type ServiceBlock struct {
	Service compose.Service
	Name    string
	EnvVar  string
}
