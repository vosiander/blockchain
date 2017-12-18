package network

import "net"

type Pool struct {
	Peers               []*Peer `json:"peers"`
	peerMap             map[string]bool
	NotificationChannel chan *Peer

	advertisedHost string
	port           string
}

func NewPool(advertisedHost, port string) *Pool {
	return &Pool{
		peerMap:             make(map[string]bool),
		NotificationChannel: make(chan *Peer, 1024),
		advertisedHost:      advertisedHost,
		port:                port,
	}
}

func (p *Pool) AddPeer(pe *Peer) {
	ipPort := pe.IP.String() + ":" + pe.Port

	_, ok := p.peerMap[ipPort]
	if ok {
		return
	}

	p.Peers = append(p.Peers, pe)
	p.peerMap[ipPort] = true

	p.NotificationChannel <- pe
}

func (pool *Pool) GetPeers() []*Peer {
	return pool.Peers
}

func (pool *Pool) GetHost() *Peer {
	return &Peer{
		IP:   net.ParseIP(pool.advertisedHost),
		Port: pool.port,
	}
}
