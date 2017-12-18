package main

import (
	"log"

	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/siklol/blockchain"
	"github.com/siklol/blockchain/handler"
	"github.com/siklol/blockchain/network"
)

func main() {
	advertisedHost := os.Getenv("ADVERTISED_HOST")
	if advertisedHost == "" {
		advertisedHost = "127.0.0.1"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// router
	router := gin.Default()

	// we assume that the genesis block is created equally on all our handlers. Else we would have different blockchains
	genesisMsg := []byte("")
	genesisTimestamp := time.Date(1985, 2, 20, 4, 59, 0, 0, time.UTC)
	chain := blockchain.NewBlockchain(blockchain.Sha256, blockchain.Hashcash, genesisMsg, genesisTimestamp)

	// networking
	pool := network.NewPool(advertisedHost, port)
	cn := network.NewChainNetwork(chain, pool)
	go cn.Listen()

	// handlers
	handler.NewDefaultHandler().Register(router)
	handler.NewBlockHandler(chain).Register(router)
	handler.NewPeersHandler(pool).Register(router)

	if err := router.Run(); err != nil {
		log.Fatalf("fatal error occured: %s", err.Error())
	}
}
