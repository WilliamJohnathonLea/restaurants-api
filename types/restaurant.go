package types

type Restaurant struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type RestaurantMenus struct {
	Restaurant Restaurant `json:"restaurant"`
	Menus      []Menu     `json:"menus"`
}
