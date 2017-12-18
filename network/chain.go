package network

import (
	"fmt"
	"log"
	"sync"
	"time"

	"errors"

	"github.com/siklol/blockchain"
)

type ChainNetwork struct {
	c                *blockchain.Blockchain
	peerNotification <-chan *Peer
	mu               *sync.Mutex
	ticker           *time.Ticker
	pool             *Pool
}

func NewChainNetwork(c *blockchain.Blockchain, pool *Pool) *ChainNetwork {
	cn := &ChainNetwork{
		ticker: time.NewTicker(time.Second * 30),
		c:      c,
		mu:     &sync.Mutex{},
		pool:   pool,
	}

	go func() {
		for range cn.ticker.C {
			for _, p := range cn.pool.Peers {
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
		cn.pool.AddPeer(p)
	}
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

	if cn.c.Exists(peerTip) {
		return
	}

	if cn.c.Tip().Index > peerTip.Index {
		// TODO push overwrite command?
		return
	}

	// FIXME this is not efficient! Use merkel tree or something else. We need to find last common index
	firstBranchIndex, err := cn.findFirstBranchIndex(p)
	if err != nil {
		log.Println(fmt.Sprintf("error finding first branch firstBranchIndex %d : ", firstBranchIndex) + err.Error())
		return
	}

	newChainPart, err := cn.receiveChainFromIndex(p, firstBranchIndex, peerTip)
	if err != nil {
		log.Println(fmt.Sprintf("error receiving chain %d : ", firstBranchIndex) + err.Error())
		return
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
	// FIXME this is also not efficient. There should be an efficient way to create and distribute peers
	if err := AddPeer(p, cn.pool.GetHost()); err != nil {
		log.Println("error pushing host peer ip to node: " + err.Error())
	}
}

func (cn *ChainNetwork) findFirstBranchIndex(p *Peer) (int64, error) {
	index := int64(0)
	for true {
		index++

		b, err := BlockAtIndex(p, index)
		if err != nil {
			log.Println(fmt.Sprintf("error requesting block at index %d : ", index) + err.Error())
			return 0, err
		}

		if cn.c.Exists(b) {
			log.Println(fmt.Sprintf("skipping block %s at index %d", b.Hash, index))
			continue
		}

		break
	}

	return index, nil
}

func (cn *ChainNetwork) receiveChainFromIndex(p *Peer, index int64, peerTip *blockchain.Block) ([]*blockchain.Block, error) {
	var newChainPart []*blockchain.Block
	for true {
		b, err := BlockAtIndex(p, index)
		if err != nil {
			log.Println(fmt.Sprintf("error requesting block at index %d : ", index) + err.Error())
			return nil, err
		}

		if b == nil {
			log.Println(fmt.Sprintf("no block at index %d : ", index))
			return nil, errors.New("no block found at index")
		}

		newChainPart = append(newChainPart, b)
		index++

		if b.Index == peerTip.Index {
			break
		}
	}

	return newChainPart, nil
}
