package types

import "time"

type Order struct {
	ID           string     `json:"id,omitempty"`
	RestaurantID string     `json:"restaurantId"`
	Items        []LineItem `json:"items"`
	CreatedAt    time.Time  `json:"createdAt"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
}

type LineItem struct {
	// ID is the unique ID for this line item entry in the database
	ID string `json:"id,omitempty"`
	// OrderID is the ID of the Order to which this LineItem belongs
	OrderID string `json:"orderId,omitempty"`
	// ItemID refers to the ID of the MenuItem
	ItemID string `json:"itemId"`
	// Name refers to the Name of the MenuItem
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity uint    `json:"quantity"`
}
