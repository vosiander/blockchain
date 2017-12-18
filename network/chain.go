package network

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/siklol/blockchain"
)

type ChainNetwork struct {
	c                *blockchain.Blockchain
	peerNotification <-chan *Peer
	mu               *sync.Mutex
	ticker           *time.Ticker
	peers            []*Peer
}

func NewChainNetwork(c *blockchain.Blockchain, cc <-chan *Peer) *ChainNetwork {
	cn := &ChainNetwork{
		ticker:           time.NewTicker(time.Second * 30),
		c:                c,
		peerNotification: cc,
		mu:               &sync.Mutex{},
	}

	go func() {
		for t := range cn.ticker.C {
			fmt.Println("timed synchronization at", t)

			for _, p := range cn.peers {
				log.Printf("Sync peer %s:%s\n", p.IP, p.Port)
				// TODO if peer is faulty, remove or deprecate it
				cn.synchronizePeer(p)
			}
		}
	}()

	return cn
}

func (cn *ChainNetwork) Listen() {
	for p := range cn.peerNotification {
		log.Printf("New peer received %s:%s\n", p.IP, p.Port)
		cn.AddPeer(p)
	}
}

func (cn *ChainNetwork) AddPeer(p *Peer) {
	cn.peers = append(cn.peers, p)
}

func (cn *ChainNetwork) synchronizePeer(p *Peer) {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if !cn.checkVersion(p) {
		return
	}

	if !cn.hasCompatibleGenesisBlock(p) {
		return
	}

	cn.broadcastPeer(p)

	peerTip, err := Tip(p)
	if err != nil {
		log.Println("error requesting tip block: " + err.Error())
		return
	}

	log.Println(fmt.Sprintf("Tip block: %#v", peerTip))

	if cn.c.Exists(peerTip) {
		log.Println("peer tip exists in chain")
		return
	}

	if cn.c.Tip().Index > peerTip.Index {
		log.Println("peer tip is behind.")
		// TODO push overwrite command?
		return
	}

	index := int64(0) // FIXME this is not efficient! Use merkel tree or something else. We need to find last common index
	for true {
		log.Printf("old: index %d", index)

		index++

		b, err := BlockAtIndex(p, index)
		if err != nil {
			log.Println(fmt.Sprintf("error requesting block at index %d : ", index) + err.Error())
			return
		}

		if cn.c.Exists(b) {
			log.Println(fmt.Sprintf("skipping block %s at index %d", b.Hash, index))
			continue
		}

		break
	}

	firstBranchIndex := index

	var newChainPart []*blockchain.Block
	for true {
		b, err := BlockAtIndex(p, index)
		if err != nil {
			log.Println(fmt.Sprintf("error requesting block at index %d : ", index) + err.Error())
			return
		}

		if b == nil {
			log.Println(fmt.Sprintf("no block at index %d : ", index) + err.Error())
			return
		}

		newChainPart = append(newChainPart, b)
		index++

		if b.Index == peerTip.Index {
			break
		}
	}

	if err := cn.c.ReplaceChainFromIndex(firstBranchIndex, newChainPart); err != nil {
		log.Println("error replacing chain: " + err.Error())
		return
	}

}

func (cn *ChainNetwork) checkVersion(p *Peer) bool {
	v, err := Version(p)
	if err != nil {
		log.Println("error requesting version number: " + err.Error())
		return false
	}

	if !blockchain.IsCompatibleWithCurrent(v) {
		log.Println("incompatible versions " + v)
		return false
	}

	return true
}

func (cn *ChainNetwork) hasCompatibleGenesisBlock(p *Peer) bool {
	genesisBlock, err := GenesisBlock(p)
	if err != nil {
		log.Println("error requesting genesis block: " + err.Error())
		return false
	}

	if !cn.c.Exists(genesisBlock) {
		log.Println(fmt.Sprintf("Genesis block is incompatible with block chain: %#v", genesisBlock))
		return false
	}

	return true
}

func (cn *ChainNetwork) broadcastPeer(p *Peer) {
	err := AddPeer(p, &Peer{
		IP:   net.ParseIP("127.0.0.1"), // FIXME change to dynamic ip
		Port: "8080",
	})

	if err != nil {
		log.Println("error pushing host peer ip to node: " + err.Error())
	}
}
