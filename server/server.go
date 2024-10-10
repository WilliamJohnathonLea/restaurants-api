package server

import (
	"fmt"

	"github.com/WilliamJohnathonLea/restaurants-api/notifier"
	"github.com/WilliamJohnathonLea/restaurants-api/services/orders"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
)

type ServerOpt func(*ServerApp)

type RouteHandler func(*ServerApp) gin.HandlerFunc

type ServerApp struct {
	Port       int
	TokenKey   string
	Router     *gin.Engine
	DbSession  *dbr.Session
	Notifier   *notifier.RabbitNotifer
	ordersRepo orders.OrdersRepo
}

func New(opts ...ServerOpt) *ServerApp {
	app := &ServerApp{}
	router := gin.Default()
	port := 8080

	app.Router = router
	app.Port = port

	for _, opt := range opts {
		opt(app)
	}

	if app.DbSession != nil {
		app.ordersRepo = orders.NewRepo(app.DbSession)
	}

	return app
}

func WithDbSession(sess *dbr.Session) ServerOpt {
	return func(sa *ServerApp) {
		sa.DbSession = sess
	}
}

func WithPort(port int) ServerOpt {
	return func(sa *ServerApp) {
		sa.Port = port
	}
}

func WithNotifier(rn *notifier.RabbitNotifer) ServerOpt {
	return func(sa *ServerApp) {
		sa.Notifier = rn
	}
}

func WithTokenKey(key string) ServerOpt {
	return func(sa *ServerApp) {
		sa.TokenKey = key
	}
}

func WithRoute(method, path string, handler RouteHandler) ServerOpt {
	return func(sa *ServerApp) {
		sa.Router.Handle(method, path, handler(sa))
	}
}

func WithAuthRoute(method, path string, handler RouteHandler) ServerOpt {
	return func(sa *ServerApp) {
		sa.Router.Handle(method, path, Authenticated(sa.TokenKey), handler(sa))
	}
}

func (sa *ServerApp) Run() error {
	addr := fmt.Sprintf(":%d", sa.Port)
	return sa.Router.Run(addr)
}
