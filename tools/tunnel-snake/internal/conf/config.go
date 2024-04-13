package conf

import "tunnel-snake/internal/xnet"

type Config struct {
	Username string
	Remote   xnet.Server
	Ports    []uint16
}
