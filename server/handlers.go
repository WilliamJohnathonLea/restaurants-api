package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
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
		ctx.Status(http.StatusNotImplemented)
	}
}

func GetOrderByID(sa *ServerApp) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderID := ctx.Param("id")

		order, err := sa.ordersRepo.GetOrderByID(orderID)
		// Handle not found
		if err == dbr.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "order not found",
				"order": orderID,
			})
			return
		}

		// Handle other error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"error": "internal server error",
			})
			return
		}

		ctx.JSON(http.StatusOK, order)
	}
}
