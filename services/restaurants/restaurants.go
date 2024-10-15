package restaurants

import (
	"github.com/WilliamJohnathonLea/restaurants-api/db"
	"github.com/WilliamJohnathonLea/restaurants-api/types"
	"github.com/gocraft/dbr/v2"
)

type RestaurantsRepo interface {
	CreateRestaurant(types.Restaurant) error
	GetRestaurantByID(string) (types.Restaurant, error)
	GetRestaurantMenusByID(string) (types.RestaurantMenus, error)
	CreateMenu(types.Menu) error
}

type SqlRestaurantsRepo struct {
	db *dbr.Session
}

func NewRepo(db *dbr.Session) RestaurantsRepo {
	return &SqlRestaurantsRepo{db}
}

// CreateMenu implements RestaurantsRepo.
func (rr *SqlRestaurantsRepo) CreateMenu(menu types.Menu) error {
	tx, err := rr.db.Begin()
	if err != nil {
		return err
	}

	err = db.WithTx(tx, func() error {
		_, err := tx.InsertInto("menus").
			Columns("id", "restaurant_id", "name").
			Values(menu.ID, menu.RestaurantID, menu.Name).
			Exec()

		if len(menu.Items) > 0 {
			itemQuery := tx.InsertInto("menu_items").
				Columns("id", "menu_id", "name", "price")

			for _, item := range menu.Items {
				itemQuery.Values(item.ID, menu.ID, item.Name, item.Price)
			}
			_, err = itemQuery.Exec()
		}

		return err
	})

	return err
}

// CreateRestaurant implements RestaurantsRepo.
func (rr *SqlRestaurantsRepo) CreateRestaurant(restaurant types.Restaurant) error {
	tx, err := rr.db.Begin()
	if err != nil {
		return err
	}

	err = db.WithTx(tx, func() error {
		_, err := tx.InsertInto("restaurants").
			Columns("id", "name").
			Values(restaurant.ID, restaurant.Name).
			Exec()
		return err
	})

	return err
}

// GetRestaurantByID implements RestaurantsRepo.
func (rr *SqlRestaurantsRepo) GetRestaurantByID(id string) (types.Restaurant, error) {
	var restaurant types.Restaurant
	tx, err := rr.db.Begin()
	if err != nil {
		return restaurant, err
	}

	err = db.WithTx(tx, func() error {
		return tx.Select("id", "name").
			From("restaurants").
			Where("id = ?", id).
			LoadOne(&restaurant)
	})

	return restaurant, err
}

// GetRestaurantMenusByID implements RestaurantsRepo.
func (rr *SqlRestaurantsRepo) GetRestaurantMenusByID(id string) (types.RestaurantMenus, error) {
	panic("unimplemented")
}
