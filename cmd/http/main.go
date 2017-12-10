package main

import (
	"log"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/siklol/blockchain"
	"github.com/siklol/blockchain/handler"
)

func main() {
	// router
	router := gin.Default()

	// we assume that the genesis block is created equally on all our handlers. Else we would have different blockchains
	genesisMsg := []byte("")
	genesisTimestamp := time.Date(1985, 2, 20, 4, 59, 0, 0, time.UTC)
	chain := blockchain.NewBlockchain(blockchain.Sha256, blockchain.Hashcash, genesisMsg, genesisTimestamp)

	// handlers
	handler.NewDefaultHandler().Register(router)
	handler.NewBlockHandler(chain).Register(router)

	if err := router.Run(); err != nil {
		log.Fatalf("fatal error occured: %s", err.Error())
	}
}
