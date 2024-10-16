package server

import (
	"fmt"

	"github.com/WilliamJohnathonLea/restaurants-api/notifier"
	"github.com/WilliamJohnathonLea/restaurants-api/services/orders"
	"github.com/WilliamJohnathonLea/restaurants-api/services/restaurants"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
)

type ServerOpt func(*ServerApp)

type RouteHandler func(*ServerApp) gin.HandlerFunc

type ServerApp struct {
	port            int
	tokenKey        string
	router          *gin.Engine
	dbSession       *dbr.Session
	notifier        *notifier.RabbitNotifer
	ordersRepo      orders.OrdersRepo
	restaurantsRepo restaurants.RestaurantsRepo
}

func New(opts ...ServerOpt) *ServerApp {
	app := &ServerApp{}
	router := gin.Default()
	port := 8080

	app.router = router
	app.port = port

	for _, opt := range opts {
		opt(app)
	}

	if app.dbSession != nil {
		app.ordersRepo = orders.NewRepo(app.dbSession)
		app.restaurantsRepo = restaurants.NewRepo(app.dbSession)
	}

	return app
}

func WithDbSession(sess *dbr.Session) ServerOpt {
	return func(sa *ServerApp) {
		sa.dbSession = sess
	}
}

func WithPort(port int) ServerOpt {
	return func(sa *ServerApp) {
		sa.port = port
	}
}

func WithNotifier(rn *notifier.RabbitNotifer) ServerOpt {
	return func(sa *ServerApp) {
		sa.notifier = rn
	}
}

func WithTokenKey(key string) ServerOpt {
	return func(sa *ServerApp) {
		sa.tokenKey = key
	}
}

func WithRoute(method, path string, handler RouteHandler) ServerOpt {
	return func(sa *ServerApp) {
		sa.router.Handle(method, path, handler(sa))
	}
}

func WithAuthRoute(method, path string, handler RouteHandler) ServerOpt {
	return func(sa *ServerApp) {
		sa.router.Handle(method, path, Authenticated(sa.tokenKey), handler(sa))
	}
}

func (sa *ServerApp) Run() error {
	addr := fmt.Sprintf(":%d", sa.port)
	return sa.router.Run(addr)
}
