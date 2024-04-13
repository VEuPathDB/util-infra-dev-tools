package env

import (
	"strings"
)

func BuildNewEnv(hosts map[string]string) string {
	sb := strings.Builder{}
	sb.Grow(2048)

	sb.WriteString("\n# Tunnel Host Connection Info\n")
	sb.WriteString(TunnelHost)
	sb.WriteString("=\n")
	sb.WriteString(TunnelPort)
	sb.WriteString("=\n")
	sb.WriteString(TunnelUser)
	sb.WriteString("=\n")

	sb.WriteString("\n# Container Hosts\n")
	for key, host := range hosts {
		sb.WriteString(key)
		sb.WriteByte('=')
		sb.WriteString(host)
		sb.WriteByte('\n')
	}

	sb.WriteByte('\n')

	return sb.String()
}
