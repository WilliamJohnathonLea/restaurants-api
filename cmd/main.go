package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/WilliamJohnathonLea/restaurants-api/server"
	"github.com/gocraft/dbr/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load environment")
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")

	amqpUsername := os.Getenv("AMQP_USERNAME")
	amqpPassword := os.Getenv("AMQP_PASSWORD")
	amqpHost := os.Getenv("AMQP_HOST")

	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))

	// Open DB connection
	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s/restaurants?sslmode=disable",
		dbUsername,
		dbPassword,
		dbHost,
	)
	conn, err := dbr.Open("postgres", dbUrl, nil)
	if err != nil {
		log.Fatalf("error opening db connection %+v", err)
	}
	defer conn.Close()

	sess := conn.NewSession(nil)
	defer sess.Close()

	// Set up server
	server := server.New(
		server.WithPort(port),
		server.WithDbSession(sess),
	)

	server.Run()

}
