package restaurants

import (
	"github.com/WilliamJohnathonLea/restaurants-api/types"
	"github.com/gocraft/dbr/v2"
)

type RestaurantsRepo interface {
	CreateRestaurant(types.Restaurant) error
	GetRestaurantByID(string) (types.Restaurant, error)
	GetRestaurantMenusByID(string) (types.RestaurantMenus, error)
	CreateMenu(string, types.Menu) error
}

type SqlRestaurantsRepo struct {
	db *dbr.Session
}

func NewRepo(db *dbr.Session) RestaurantsRepo {
	return &SqlRestaurantsRepo{db}
}

// CreateMenu implements RestaurantsRepo.
func (rr *SqlRestaurantsRepo) CreateMenu(restaurantID string, menu types.Menu) error {
	panic("unimplemented")
}

// CreateRestaurant implements RestaurantsRepo.
func (rr *SqlRestaurantsRepo) CreateRestaurant(restaurant types.Restaurant) error {
	panic("unimplemented")
}

// GetRestaurantByID implements RestaurantsRepo.
func (rr *SqlRestaurantsRepo) GetRestaurantByID(id string) (types.Restaurant, error) {
	panic("unimplemented")
}

// GetRestaurantMenusByID implements RestaurantsRepo.
func (rr *SqlRestaurantsRepo) GetRestaurantMenusByID(id string) (types.RestaurantMenus, error) {
	panic("unimplemented")
}
