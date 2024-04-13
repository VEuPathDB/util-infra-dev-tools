package conf

import (
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"tunnel-snake/internal/xnet"
)

func LoadFromEnv() Config {
	user := requireEnv("TUNNEL_USER")
	host, err := xnet.ParseServerAddress(requireEnv("TUNNEL_HOST"))
	if err != nil {
		log.Fatalf("failed to parse ssh host address: %s", err)
		panic(nil) // unreachable
	}

	portStrings := strings.Split(requireEnv("TUNNEL_PORTS"), ",")
	ports := make([]uint16, len(portStrings))

	for i, raw := range portStrings {
		port, err := strconv.ParseUint(raw, 10, 16)

		if err != nil {
			log.Fatalf("could not parse port string: %s", err)
			panic(nil) // unreachable
		}

		ports[i] = uint16(port)
	}

	return Config{user, host, ports}
}

func requireEnv(key string) string {
	out := os.Getenv(key)

	if len(out) == 0 {
		log.Fatalf("missing required environment variable %s", key)
		panic(nil) // unreachable
	}

	return out
}
