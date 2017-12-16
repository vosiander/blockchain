package network

type Pool struct {
	Peers               []*Peer `json:"peers"`
	peerMap             map[string]bool
	NotificationChannel chan *Peer
}

func NewPool() *Pool {
	return &Pool{
		peerMap:             make(map[string]bool),
		NotificationChannel: make(chan *Peer, 1024),
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
