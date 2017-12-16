package network

import (
	"fmt"
	"log"
	"sync"

	"github.com/siklol/blockchain"
)

// This could also be done via kafka or some other messaging infrastructure
type ChainNetwork struct {
	c                *blockchain.Blockchain
	peerNotification <-chan *Peer
	mu               *sync.Mutex
}

func NewChainNetwork(c *blockchain.Blockchain, cc <-chan *Peer) *ChainNetwork {
	return &ChainNetwork{
		c:                c,
		peerNotification: cc,
		mu:               &sync.Mutex{},
	}
}

func (cn *ChainNetwork) Listen() {
	for p := range cn.peerNotification {
		log.Printf("New peer received %s:%s\n", p.IP, p.Port)

		// TODO add all peers from peer
		cn.synchronizePeer(p)
	}
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

	// TODO download all blocks which are not currently in the chain
	currentTip := cn.c.Tip()
	index := currentTip.Index
	for true {
		index++

		b, err := BlockAtIndex(p, index)
		if err != nil {
			log.Println(fmt.Sprintf("error requesting block at index %d : ", index) + err.Error())
			return
		}

		currentTip = b

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
			return
		}
	}
}
