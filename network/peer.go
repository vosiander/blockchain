package network

import "net"

type Peer struct {
	IP   net.IP `json:"ip"`
	Port string `json:"port"`
}
