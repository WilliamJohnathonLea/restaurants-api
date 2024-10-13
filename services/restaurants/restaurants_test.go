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
	db *dbr.Session
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
		ID: uuid.NewString(),
		Name: "Test",
	})
	s.Nil(err)
}

func (s *RestaurantRepoSuite) TestGetRestaurantByID() {
	rID := uuid.NewString()
	expected := types.Restaurant{
		ID: rID,
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

func TestRunSuite(t *testing.T) {
	repoSuite := &RestaurantRepoSuite{}
	m, err := db.NewMigrator(dbUrl, migrationsUrl)
	if err != nil {
		t.Fatalf("error setting up migrator %s", err.Error())
	}
	defer m.Close()
	m.Run()

	conn, err := dbr.Open("postgres", dbUrl, nil)
	if err != nil {
		t.Fatalf("error setting up db %s", err.Error())
	}
	defer conn.Close()
	repoSuite.db = conn.NewSession(nil)
	repoSuite.repo = restaurants.NewRepo(repoSuite.db)

	suite.Run(t, repoSuite)
}
