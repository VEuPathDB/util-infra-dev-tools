package ssh

import (
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type endpoint struct {
	tunnel *ssh.Client
}

func (e endpoint) Open(port uint16) {
	localAddr := fmt.Sprintf("0.0.0.0:%d", port)
	channel, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("failed to create channel %s", localAddr)
		panic(nil) // unreachable
	}
	defer channel.Close()

	for {
		connection, err := channel.Accept()

		if err != nil {
			log.Errorf("error while attempting to accept connection: %s", err)
			continue
		}

		log.Infof("accepted connection %s", connection.LocalAddr())

		go e.forward(connection)
	}
}

func (e endpoint) forward(incoming net.Conn) {
	outgoing, err := e.tunnel.Dial("tcp", incoming.LocalAddr().String())
	if err != nil {
		log.Errorf("failed to establish tunneled connection to %s: %s", incoming.LocalAddr(), err)
		return
	}
	log.Infof("created new tunnel to remote address %s through %s", incoming.LocalAddr(), e.tunnel.RemoteAddr())

	go pipe(incoming, outgoing)
	go pipe(outgoing, incoming)
}

func pipe(to, from net.Conn) {
	_, err := io.Copy(to, from)
	if err != nil {
		log.Errorf("tunnel error: %s", err)
	}
}
