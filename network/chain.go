package network

import (
	"fmt"
	"log"
	"sync"
	"time"

	"net"

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

	v, err := Version(p)
	if err != nil {
		log.Println("error requesting version number: " + err.Error())
		return
	}

	log.Println(fmt.Sprintf("Version number from peer: %d", v))

	genesisBlock, err := GenesisBlock(p)
	if err != nil {
		log.Println("error requesting genesis block: " + err.Error())
		return
	}

	log.Println(fmt.Sprintf("Genesis block: %#v", genesisBlock))

	if !cn.c.Exists(genesisBlock) {
		log.Println(fmt.Sprintf("Genesis block is incompatible with block chain: %#v", genesisBlock))
		return
	}

	tip, err := Tip(p)
	if err != nil {
		log.Println("error requesting tip block: " + err.Error())
		return
	}

	log.Println(fmt.Sprintf("Tip block: %#v", tip))

	if cn.c.Exists(tip) {
		// TODO push own ip and port to node
		log.Println("peer tip exists in chain")
		return
	}

	index := cn.c.Tip().Index
	for true {
		index++

		b, err := BlockAtIndex(p, index)
		if err != nil {
			log.Println(fmt.Sprintf("error requesting block at index %d : ", index) + err.Error())
			return
		}

		// TODO if the other chain has longer chains, use their chain

		currentTip := b
		if b == nil {
			log.Println(fmt.Sprintf("no block at index %d : ", index) + err.Error())
			return
		}

		if err := cn.c.Append(b); err != nil {
			log.Println("Could not append block to current blockchain. error: " + err.Error())
			return
		}

		if currentTip.Hash == tip.Hash {
			log.Println("finished synchronization. tip equals latest block from peer")
			break
		}
	}

	// TODO announce peer
	err = AddPeer(p, &Peer{
		IP:   net.ParseIP("192.168.178.20"), // FIXME change to dynamic ip
		Port: "8080",
	})
	if err != nil {
		log.Println("error pushing host peer ip to node: " + err.Error())
		return
	}
}
