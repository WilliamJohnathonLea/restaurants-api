package types

type Menu struct {
	ID           string     `json:"id,omitempty"`
	RestaurantID string     `json:"restaurantId"`
	Name         string     `json:"name"`
	Items        []MenuItem `json:"items"`
}

type MenuItem struct {
	ID     string  `json:"id,omitempty"`
	MenuID string  `json:"menuId,omitempty"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}
