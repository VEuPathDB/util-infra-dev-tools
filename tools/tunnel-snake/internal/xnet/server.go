package xnet

import (
	"fmt"
	"strconv"
	"strings"
)

type Server struct {
	Host string
	Port uint16
}

func (s Server) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func ParseServerAddress(addr string) (Server, error) {
	c := strings.LastIndexByte(addr, ':')
	if c < 1 {
		return Server{}, fmt.Errorf("invalid server address: %s", addr)
	}

	port, err := strconv.ParseUint(addr[c+1:], 10, 16)
	if err != nil {
		return Server{}, fmt.Errorf("invalid port for address %s", addr)
	}

	return Server{addr[:c], uint16(port)}, nil
}
