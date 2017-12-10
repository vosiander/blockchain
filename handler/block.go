package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/siklol/blockchain"
)

type postBlockData struct {
	Data string `json:"data" binding:"required"`
}

type blockHandler struct {
	chain *blockchain.Blockchain
}

func NewBlockHandler(chain *blockchain.Blockchain) *blockHandler {
	return &blockHandler{
		chain: chain,
	}
}

func (h *blockHandler) Register(router *gin.Engine) {
	router.GET("/blocks", func(c *gin.Context) { h.GetBlocks(c) })
	router.GET("/blocks/tip", func(c *gin.Context) { h.GetTip(c) })
	router.POST("/blocks", func(c *gin.Context) { h.PostBlocks(c) })
}

func (h *blockHandler) GetBlocks(c *gin.Context) {
	blocks := struct {
		Blocks []*blockchain.Block `json:"blocks"`
	}{
		Blocks: h.chain.Blocks(),
	}

	c.JSON(200, blocks)
}

func (h *blockHandler) GetTip(c *gin.Context) {
	c.JSON(200, h.chain.Tip())
}

func (h *blockHandler) PostBlocks(c *gin.Context) {
	var jsonData postBlockData
	c.Bind(&jsonData)

	if err := h.chain.Mine([]byte(jsonData.Data)); err != nil {
		log.Println("error mining a block: " + err.Error())
		c.JSON(500, gin.H{
			"error": "problem occured mining a new block",
		})
		return
	}

	c.JSON(200, h.chain.Tip())
}
