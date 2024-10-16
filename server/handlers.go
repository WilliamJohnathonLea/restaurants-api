package server

import (
	"encoding/json"
	"net/http"

	"github.com/WilliamJohnathonLea/restaurants-api/notifier"
	"github.com/WilliamJohnathonLea/restaurants-api/types"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
	"github.com/google/uuid"
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
		order, err := fromBody[types.Order](ctx)
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
			Headers:    map[string]interface{}{"user_id": order.UserID},
			Body:       orderBytes,
		})
		if err != nil {
			setErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusAccepted, order)
	}
}

func PostNewRestaurant(sa *ServerApp) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		restaurant, err := fromBody[types.Restaurant](ctx)
		if err != nil {
			setErrorResponse(ctx, http.StatusBadRequest, err)
			return
		}

		restaurant.ID = uuid.NewString()
		err = sa.restaurantsRepo.CreateRestaurant(restaurant)
		if err != nil {
			setErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusCreated, restaurant)
	}
}

func PostNewMenu(sa *ServerApp) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		menu, err := fromBody[types.Menu](ctx)
		if err != nil {
			setErrorResponse(ctx, http.StatusBadRequest, err)
			return
		}

		// Set new ID for the menu and it's items
		mID := uuid.NewString()
		menu.ID = mID
		for i := range menu.Items {
			menu.Items[i].ID = uuid.NewString()
			menu.Items[i].MenuID = mID
		}

		err = sa.restaurantsRepo.CreateMenu(menu)
		if err != nil {
			setErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusCreated, menu)
	}
}

func setErrorResponse(ctx *gin.Context, httpCode int, err error) {
	ctx.JSON(httpCode, gin.H{
		"error": err.Error(),
	})
}

func fromBody[T any](ctx *gin.Context) (T, error) {
	var result T

	err := ctx.BindJSON(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
