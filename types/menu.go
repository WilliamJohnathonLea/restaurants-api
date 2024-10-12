package types

type Menu struct {
	ID           string     `json:"id,omitempty"`
	RestaurantID string     `json:"restaurantId"`
	Name         string     `json:"name"`
	Items        []MenuItem `json:"items"`
}

type MenuItem struct {
	ID     string `json:"id,omitempty"`
	MenuID string `json:"menuId"`
	Name   string `json:"name"`
}
