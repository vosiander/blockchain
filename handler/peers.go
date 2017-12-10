package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/siklol/blockchain/network"
)

type peersHandler struct {
	pool *network.Pool
}

func NewPeersHandler(pool *network.Pool) *peersHandler {
	return &peersHandler{
		pool: pool,
	}
}

func (h *peersHandler) Register(router *gin.Engine) {
	router.GET("/peers", func(c *gin.Context) { h.GetPeers(c) })
	router.POST("/peers", func(c *gin.Context) { h.PostPeers(c) })
}

func (h *peersHandler) GetPeers(c *gin.Context) {
	c.JSON(200, struct {
		Peers []*network.Peer `json:"peers"`
	}{
		Peers: h.pool.GetPeers(),
	})
}

func (h *peersHandler) PostPeers(c *gin.Context) {
	var peer *network.Peer
	if err := c.Bind(&peer); err != nil {
		c.AbortWithError(400, err)
		return
	}

	h.pool.AddPeer(peer)

	c.JSON(200, peer)
}
