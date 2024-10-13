package orders_test

import (
	"testing"
	"time"

	"github.com/WilliamJohnathonLea/restaurants-api/db"
	"github.com/WilliamJohnathonLea/restaurants-api/services/orders"
	"github.com/WilliamJohnathonLea/restaurants-api/types"
	"github.com/gocraft/dbr/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type OrdersRepoSuite struct {
	suite.Suite
	db *dbr.Session
	repo orders.OrdersRepo

	order types.Order
}

const dbUrl = "postgres://postgres:postgres@localhost:5432/restaurants?sslmode=disable"
const migrationsUrl = "file://../../db/migrations"

// SetupTest is called before each test
func (s *OrdersRepoSuite) SetupTest() {
	s.db.DeleteFrom("orders").Exec()
	s.db.DeleteFrom("line_items").Exec()

	orderID := uuid.NewString()
	s.order = types.Order{
		ID:           orderID,
		RestaurantID: uuid.NewString(),
		UserID:       uuid.NewString(),
		Items: []types.LineItem{
			{
				ID:       uuid.NewString(),
				OrderID:  orderID,
				ItemID:   uuid.NewString(),
				Name:     "test",
				Price:    1.00,
				Quantity: 1,
			},
		},
		CreatedAt: time.UnixMilli(0).UTC(),
	}

	s.db.InsertInto("restaurants").
		Columns("id", "name").
		Values(s.order.RestaurantID, "Test").
		Exec()

	s.db.InsertInto("orders").
		Columns("id", "restaurant_id", "user_id", "created_at").
		Values(s.order.ID, s.order.RestaurantID, s.order.UserID, s.order.CreatedAt).
		Exec()

	lq := s.db.InsertInto("line_items").
		Columns("id", "order_id", "item_id", "name", "price", "quantity")

	for _, i := range s.order.Items {
		lq.Record(i)
	}
	lq.Exec()
}

func (s *OrdersRepoSuite) TestGetOrderByID() {

	s.Run("order exists", func() {
		dbOrder, err := s.repo.GetOrderByID(s.order.ID)
		s.Nil(err)
		s.Equal(s.order, dbOrder)
	})

	s.Run("order does not exist", func() {
		_, err := s.repo.GetOrderByID(uuid.NewString())
		s.NotNil(err)
	})
}

func TestRunOrdersRepoSuite(t *testing.T) {
	repoSuite := &OrdersRepoSuite{}
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
	repoSuite.repo = orders.NewRepo(repoSuite.db)

	suite.Run(t, repoSuite)
}
