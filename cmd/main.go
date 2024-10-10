package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/WilliamJohnathonLea/restaurants-api/db"
	"github.com/WilliamJohnathonLea/restaurants-api/notifier"
	"github.com/WilliamJohnathonLea/restaurants-api/server"
	"github.com/gocraft/dbr/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	fatalOnError(err, "failed to load environment")

	dbAdminUsername := os.Getenv("DB_ADMIN_USERNAME")
	dbAdminPassword := os.Getenv("DB_ADMIN_PASSWORD")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")

	amqpUsername := os.Getenv("AMQP_USERNAME")
	amqpPassword := os.Getenv("AMQP_PASSWORD")
	amqpHost := os.Getenv("AMQP_HOST")

	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	tokenDecoderKey := os.Getenv("AUTH_TOKEN_DECODER_KEY")

	// Run migrations
	migrationDbUrl := restaurantsDbUrl(dbAdminUsername, dbAdminPassword, dbHost)
	m, err := db.NewMigrator(migrationDbUrl, "file://db/migrations")
	fatalOnError(err, "error setting up migrator")
	defer m.Close()

	err = m.Run()
	fatalOnError(err, "error running migration")
	m.Close() // Close the migrator's connection after success

	// Open DB connection
	dbUrl := restaurantsDbUrl(dbUsername, dbPassword, dbHost)
	conn, err := dbr.Open("postgres", dbUrl, nil)
	fatalOnError(err, "error opening db connection")
	defer conn.Close()

	sess := conn.NewSession(nil)
	defer sess.Close()

	// Set up AMQP
	amqpUrl := fmt.Sprintf(
		"amqp://%s:%s@%s/",
		amqpUsername,
		amqpPassword,
		amqpHost,
	)
	rn, err := notifier.NewRabbitNotifier(
		notifier.WithURL(amqpUrl),
	)
	fatalOnError(err, "failed to connect to rabbitmq")
	defer rn.Close()

	// Set up server
	server := server.New(
		server.WithPort(port),
		server.WithTokenKey(tokenDecoderKey),
		server.WithDbSession(sess),
		server.WithNotifier(rn),
		server.WithRoute("GET", "/healthcheck", server.Health),
		server.WithAuthRoute("POST", "/orders", server.PostNewOrder),
		server.WithAuthRoute("GET", "/orders", server.GetOrders),
		server.WithAuthRoute("GET", "/orders/:id", server.GetOrderByID),
	)

	server.Run()

}

func restaurantsDbUrl(username, password, host string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/restaurants?sslmode=disable",
		username,
		password,
		host,
	)
}

func fatalOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s %s", msg, err.Error())
	}
}
