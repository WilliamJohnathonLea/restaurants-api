package restaurants_test

import (
	"testing"

	"github.com/WilliamJohnathonLea/restaurants-api/db"
	"github.com/WilliamJohnathonLea/restaurants-api/services/restaurants"
	"github.com/WilliamJohnathonLea/restaurants-api/types"
	"github.com/gocraft/dbr/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type RestaurantRepoSuite struct {
	suite.Suite
	db   *dbr.Session
	repo restaurants.RestaurantsRepo
}

const dbUrl = "postgres://postgres:postgres@localhost:5432/restaurants?sslmode=disable"
const migrationsUrl = "file://../../db/migrations"

func (s *RestaurantRepoSuite) SetupTest() {
	s.db.DeleteFrom("restaurants").Exec()
	s.db.DeleteFrom("menus").Exec()
	s.db.DeleteFrom("menu_items").Exec()
}

func (s *RestaurantRepoSuite) TestCreateRestaurant() {
	err := s.repo.CreateRestaurant(types.Restaurant{
		ID:   uuid.NewString(),
		Name: "Test",
	})
	s.Nil(err)
}

func (s *RestaurantRepoSuite) TestGetRestaurantByID() {
	rID := uuid.NewString()
	expected := types.Restaurant{
		ID:   rID,
		Name: "Test",
	}
	err := s.repo.CreateRestaurant(expected)
	if err != nil {
		s.FailNow("restaurant should be inserted", "err: %s", err.Error())
	}

	s.Run("finds existing restaurant", func() {
		actual, err := s.repo.GetRestaurantByID(rID)
		s.Nil(err)
		s.Equal(expected, actual)
	})

	s.Run("cannot find restaurant", func() {
		_, err := s.repo.GetRestaurantByID(uuid.NewString())
		s.NotNil(err)
		s.Equal(dbr.ErrNotFound, err)
	})
}

func (s *RestaurantRepoSuite) TestCreateMenu() {
	rID := uuid.NewString()
	expected := types.Restaurant{
		ID:   rID,
		Name: "Test",
	}
	err := s.repo.CreateRestaurant(expected)
	s.Nil(err, "a restaurant is required to attach a menu to")

	err = s.repo.CreateMenu(types.Menu{
		ID:           uuid.NewString(),
		RestaurantID: rID,
		Name:         "Starters",
	})
	s.Nil(err)
}

func (s *RestaurantRepoSuite) TestGetRestaurantMenusByID() {
	s.Run("find a restaurant and its menus", func() {
		rID := uuid.NewString()
		mID := uuid.NewString()
		iID := uuid.NewString()

		err := s.repo.CreateRestaurant(types.Restaurant{
			ID:   rID,
			Name: "Test",
		})
		s.Nil(err, "a restaurant is required to attach a menu to")

		err = s.repo.CreateMenu(types.Menu{
			ID:           mID,
			RestaurantID: rID,
			Name:         "Starters",
			Items: []types.MenuItem{
				{ID: iID, MenuID: mID, Name: "Item", Price: 1.99},
			},
		})
		s.Nil(err)

		rm, err := s.repo.GetRestaurantMenusByID(rID)
		s.Nil(err)
		s.Equal(rID, rm.Restaurant.ID)
		s.Equal(mID, rm.Menus[0].ID)
		s.Equal(iID, rm.Menus[0].Items[0].ID)
	})

	s.Run("find a restaurant with no menu", func() {
		rID := uuid.NewString()

		err := s.repo.CreateRestaurant(types.Restaurant{
			ID:   rID,
			Name: "Test",
		})
		s.Nil(err, "a restaurant is required to attach a menu to")

		rm, err := s.repo.GetRestaurantMenusByID(rID)
		s.Nil(err)
		s.Equal(rID, rm.Restaurant.ID)
	})

	s.Run("find a restaurant with a menu with no items", func() {
		rID := uuid.NewString()
		mID := uuid.NewString()

		err := s.repo.CreateRestaurant(types.Restaurant{
			ID:   rID,
			Name: "Test",
		})
		s.Nil(err, "a restaurant is required to attach a menu to")

		err = s.repo.CreateMenu(types.Menu{
			ID:           mID,
			RestaurantID: rID,
			Name:         "Starters",
		})
		s.Nil(err)

		rm, err := s.repo.GetRestaurantMenusByID(rID)
		s.Nil(err)
		s.Equal(rID, rm.Restaurant.ID)
		s.Equal(mID, rm.Menus[0].ID)
	})

	s.Run("error on non-existent restaurant", func() {
		_, err := s.repo.GetRestaurantMenusByID(uuid.NewString())
		s.Equal(dbr.ErrNotFound, err)
	})
}

func TestRunSuite(t *testing.T) {
	repoSuite := &RestaurantRepoSuite{}
	m, err := db.NewMigrator(dbUrl, migrationsUrl)
	if err != nil {
		t.Fatalf("error setting up migrator %s", err.Error())
	}
	defer m.Close()
	err = m.Run()
	if err != nil {
		t.Fatalf("error running migration %s", err.Error())
	}

	conn, err := dbr.Open("postgres", dbUrl, nil)
	if err != nil {
		t.Fatalf("error setting up db %s", err.Error())
	}
	defer conn.Close()
	repoSuite.db = conn.NewSession(nil)
	repoSuite.repo = restaurants.NewRepo(repoSuite.db)

	suite.Run(t, repoSuite)
}
