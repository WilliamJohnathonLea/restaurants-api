package server

import (
	"encoding/json"
	"net/http"

	"github.com/WilliamJohnathonLea/restaurants-api/notifier"
	"github.com/WilliamJohnathonLea/restaurants-api/types"
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
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		ctx.JSON(http.StatusOK, order)
	}
}

func PostNewOrder(sa *ServerApp) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var order types.Order
		err := ctx.BindJSON(&order)
		if err != nil {
			setErrorResponse(ctx, http.StatusBadRequest, err)
			return
		}

		orderBytes, err := json.Marshal(order)
		if err != nil {
			setErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}

		err = sa.notifier.Notify(notifier.RabbitNotification{
			Exchange:   "", // default exchange
			RoutingKey: "new_orders",
			Mandatory:  true,
			Body:       orderBytes,
		})
		if err != nil {
			setErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusAccepted, order)
	}
}

func setErrorResponse(ctx *gin.Context, httpCode int, err error) {
	ctx.JSON(httpCode, gin.H{
		"error": err.Error(),
	})
}
