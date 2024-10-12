package orders

import (
	"github.com/WilliamJohnathonLea/restaurants-api/types"
	"github.com/gocraft/dbr/v2"
)

type OrdersRepo interface {
	GetOrders() ([]types.Order, error)
	GetOrderByID(string) (types.Order, error)
}

type SqlOrdersRepo struct {
	db *dbr.Session
}

func NewRepo(db *dbr.Session) OrdersRepo {
	return &SqlOrdersRepo{db}
}

func (or *SqlOrdersRepo) GetOrders() ([]types.Order, error) {
	// TODO Implement after introducing users db
	return nil, nil
}

func (or *SqlOrdersRepo) GetOrderByID(orderID string) (types.Order, error) {
	var order types.Order
	var items []types.LineItem

	err := or.db.Select("id", "restaurant_id", "user_id", "created_at", "completed_at").
		From("orders").
		Where("id = ?", orderID).
		LoadOne(&order)

	if err != nil {
		return order, err
	}

	// adjust timezone to UTC
	order.CreatedAt = order.CreatedAt.UTC()
	if order.CompletedAt != nil {
		*order.CompletedAt = order.CompletedAt.UTC()
	}

	_, err = or.db.Select("id", "order_id", "item_id", "name", "price", "quantity").
		From("line_items").
		Where("order_id = ?", order.ID).
		Load(&items)

	order.Items = items

	return order, err
}
