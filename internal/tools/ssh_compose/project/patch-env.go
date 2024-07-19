package project

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"vpdb-dev-tool/internal/lib/xio"
	"vpdb-dev-tool/internal/tools/ssh_compose/env"
)

func patchEnvFile(file xio.ReqRWFile, hosts map[string]string) {
	fullLen := len(hosts) + 3 // host keys + 3 tunnel keys (host, port, username)

	matches := make(map[string]bool, fullLen)
	keys := make([]string, 0, fullLen)

	for key := range hosts {
		keys = append(keys, key)
	}

	keys = append(keys, env.TunnelHost, env.TunnelPort, env.TunnelUser)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		for _, key := range keys {
			if strings.HasPrefix(line, key+"=") {
				matches[key] = true
			}
		}
	}

	if scanner.Err() != nil {
		log.Fatalf("encountered error scanning .env file: %s\n", scanner.Err())
	}

	// If all the keys are already accounted for.
	if len(matches) >= fullLen {
		return
	}

	_, _ = file.Seek(0, io.SeekEnd)

	_reqWriteEnvLine(file, fmt.Sprintf("\n\n# Generated @ %s", time.Now().Format(time.DateTime)))

	if !matches[env.TunnelUser] {
		_reqWriteEnvPair(file, env.TunnelUser, "")
	}

	if !matches[env.TunnelHost] {
		_reqWriteEnvPair(file, env.TunnelHost, "")
	}

	if !matches[env.TunnelPort] {
		_reqWriteEnvPair(file, env.TunnelPort, "")
	}

	for key, val := range hosts {
		if !matches[key] {
			_reqWriteEnvPair(file, key, val)
		}
	}
}

func _reqWriteEnvPair(file xio.ReqRWFile, key, value string) {
	_reqWriteEnvLine(file, fmt.Sprintf("%s=%s", key, value))
}

func _reqWriteEnvLine(file xio.ReqRWFile, line string) {
	_, _ = file.WriteString(line)
	_, _ = file.WriteString("\n")
}
