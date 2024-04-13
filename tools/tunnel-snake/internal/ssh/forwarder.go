package ssh

import (
	"errors"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"tunnel-snake/internal/xnet"
)

type Forwarder interface {
	NewEndpoint(port uint16) error
}

type ForwarderConfig struct {
	Username string
	Server   xnet.Server
}

func NewForwarder(config ForwarderConfig) (Forwarder, error) {
	socketPath := os.Getenv("SSH_AUTH_SOCK")
	if len(socketPath) == 0 {
		return nil, errors.New("SSH_AUTH_SOCK environment variable blank or absent")
	}

	socket, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("encountered error while attempting to open ssh auth socket: %s", err)
	}

	client := agent.NewClient(socket)

	sshConfig := ssh.ClientConfig{
		User:            config.Username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeysCallback(client.Signers)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return forwarder{
		agent:  client,
		config: sshConfig,
		server: config.Server.String(),
	}, nil
}

type forwarder struct {
	agent  agent.ExtendedAgent
	config ssh.ClientConfig
	server string
}

func (f forwarder) NewEndpoint(port uint16) error {
	sshCon, err := ssh.Dial("tcp", f.server, &f.config)
	if err != nil {
		return err
	}
	log.Infof("opened ssh connection to %s for port %d", f.server, port)

	ep := endpoint{sshCon}

	go ep.Open(port)

	return nil
}
