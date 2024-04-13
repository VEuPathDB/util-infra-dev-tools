package hosts

import (
	"io"
	"log"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadHosts(rawStream io.Reader) map[Host][]string {
	dec := yaml.NewDecoder(rawStream)

	var config Config

	err := dec.Decode(&config)
	if err != nil {
		log.Fatalf("failed to load hosts config: %s\n", err)
	}

	out := make(map[Host][]string, len(config.Hosts))

	for k, v := range config.Hosts {
		pos := strings.LastIndexByte(k, ':')
		if pos < 1 {
			log.Fatalf("invalid host value: %s\n", k)
		}

		host := k[:pos]
		port, err := strconv.ParseUint(k[pos+1:], 10, 16)
		if err != nil {
			log.Fatalf("invalid host value: %s\n", k)
		}

		out[Host{host, uint16(port)}] = v
	}

	return out
}
