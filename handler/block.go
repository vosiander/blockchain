package handler

import (
	"log"

	"strconv"

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
	router.GET("/blocks/genesis", func(c *gin.Context) { h.GetGenesis(c) })
	router.GET("/blocks/tip", func(c *gin.Context) { h.GetTip(c) })
	router.GET("/blocks/index/:number", func(c *gin.Context) { h.GetIndex(c) })
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

func (h *blockHandler) GetGenesis(c *gin.Context) {
	c.JSON(200, h.chain.Genesis())
}

func (h *blockHandler) GetTip(c *gin.Context) {
	c.JSON(200, h.chain.Tip())
}

func (h *blockHandler) GetIndex(c *gin.Context) {
	index := c.Param("number")

	if index == "" {
		c.JSON(400, gin.H{
			"error": "invalid index number",
		})
		return
	}

	i, err := strconv.Atoi(index)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid index number",
		})
		return
	}

	b := h.chain.BlockAtIndex(int64(i))
	if b == nil {
		c.JSON(404, gin.H{
			"error": "invalid index number",
		})
		return
	}

	c.JSON(200, b)
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
