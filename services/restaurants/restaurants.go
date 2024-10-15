package restaurants

import (
	"database/sql"

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
	type rowData struct {
		RestaurantID   string          `dbr:"restaurant_id"`
		RestaurantName string          `dbr:"restaurant_name"`
		MenuID         sql.NullString  `dbr:"menu_id"`
		MenuName       sql.NullString  `dbr:"menu_name"`
		ItemID         sql.NullString  `dbr:"item_id"`
		ItemName       sql.NullString  `dbr:"item_name"`
		ItemPrice      sql.NullFloat64 `dbr:"item_price"`
	}

	var rm types.RestaurantMenus
	tx, err := rr.db.Begin()
	if err != nil {
		return rm, err
	}

	err = db.WithTx(tx, func() error {
		rows := []rowData{}
		count, err := tx.Select(
			"restaurants.id AS restaurant_id",
			"restaurants.name AS restaurant_name",
			"menus.id AS menu_id",
			"menus.name AS menu_name",
			"menu_items.id AS item_id",
			"menu_items.name AS item_name",
			"menu_items.price AS item_price",
		).From("restaurants").
			LeftJoin("menus", "restaurants.id = menus.restaurant_id").
			LeftJoin("menu_items", "menus.id = menu_items.menu_id").
			Where("restaurants.id = ?", id).
			Load(&rows)
		if err != nil {
			return err
		}
		if count < 1 {
			return dbr.ErrNotFound
		}

		var restaurant *types.Restaurant
		menuMap := make(map[string]*types.Menu)

		for _, row := range rows {
			if restaurant == nil {
				restaurant = &types.Restaurant{
					ID:   row.RestaurantID,
					Name: row.RestaurantName,
				}
				rm.Restaurant = *restaurant
			}

			// Add a menu if it's not already in the menuMap
			if _, ok := menuMap[row.MenuID.String]; !ok && row.MenuID.String != "" {
				menu := types.Menu{
					ID:           row.MenuID.String,
					RestaurantID: row.RestaurantID,
					Name:         row.MenuName.String,
				}
				rm.Menus = append(rm.Menus, menu)
				menuMap[row.MenuID.String] = &rm.Menus[len(rm.Menus)-1]
			}

			// Add an item to the menu if it exists
			if row.ItemID.String != "" {
				item := types.MenuItem{
					ID:    row.ItemID.String,
					Name:  row.ItemName.String,
					Price: row.ItemPrice.Float64,
				}
				if menu, exists := menuMap[row.MenuID.String]; exists {
					item.MenuID = menu.ID
					menu.Items = append(menu.Items, item)
				}
			}
		}

		return err
	})

	return rm, err
}
