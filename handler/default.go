package handler

import "github.com/gin-gonic/gin"

type defaultHandler struct {
}

func NewDefaultHandler() *defaultHandler {
	return &defaultHandler{}
}

func (h *defaultHandler) Register(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) { h.GetDefault(c) })
	router.GET("/version", func(c *gin.Context) { h.GetVersion(c) })
}

func (h *defaultHandler) GetDefault(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world!",
	})
}

func (h *defaultHandler) GetVersion(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": 1, // TODO do not hardcode this
	})
}
