package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/siklol/blockchain/handler"
)

func main() {

	// router
	router := gin.Default()

	// handlers
	handler.NewDefaultHandler().Register(router)

	if err := router.Run(); err != nil {
		log.Fatalf("fatal error occured: %s", err.Error())
	}
}
