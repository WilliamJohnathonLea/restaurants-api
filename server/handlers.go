package server

import (
	"github.com/gin-gonic/gin"
)

func Health(sa *ServerApp) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"health": "okay",
		})
	}
}

func GetOrders(sa *ServerApp) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
