package network

type Pool struct {
	Peers   []*Peer `json:"peers"`
	peerMap map[string]bool
}

func NewPool() *Pool {
	return &Pool{
		peerMap: make(map[string]bool),
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
}

func (pool *Pool) GetPeers() []*Peer {
	return pool.Peers
}
