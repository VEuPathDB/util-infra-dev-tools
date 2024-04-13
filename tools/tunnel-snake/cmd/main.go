package main

import (
	"time"

	"tunnel-snake/internal/conf"
	"tunnel-snake/internal/ssh"
)

func main() {
	config := conf.LoadFromEnv()

	server, err := ssh.NewForwarder(ssh.ForwarderConfig{
		Username: config.Username,
		Server:   config.Remote,
	})
	if err != nil {
		panic(err)
	}

	for _, port := range config.Ports {
		if err = server.NewEndpoint(port); err != nil {
			panic(err)
		}
	}

	for {
		time.Sleep(5 * time.Second)
	}
}
