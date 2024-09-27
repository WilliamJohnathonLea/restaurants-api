package main

import (
	"log"
	"os"
	"strconv"

	"github.com/WilliamJohnathonLea/restaurants-api/server"
	"github.com/joho/godotenv"
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

	server := server.New(
		server.WithPort(port),
	)

	server.Run()

}
