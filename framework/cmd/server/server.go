package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Nesrux/api-enconder/application/services"
	"github.com/Nesrux/api-enconder/framework/database"
	"github.com/Nesrux/api-enconder/framework/queue"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var db database.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	autoMigrate, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("error parsing key: AUTO_MIGRATE_DB to boolean")
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("error parsing key: DEBUG to boolean")
	}

	db.AutoMigrateDb = autoMigrate
	db.Debug = debug
	db.DsnTest = os.Getenv("DSN_TEST")
	db.Dsn = os.Getenv("DSN")
	db.DbTypeTest = os.ExpandEnv("DB_TYPE_TEST")
	db.DbType = os.ExpandEnv("DB_TYPE")
	db.Env = os.ExpandEnv("ENV")

}

func main() {
	messageChanel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalf("error connecting to DataBase")
	}
	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMq()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChanel)

	jobManager := services.NewJobManager(dbConnection, rabbitMQ,
		jobReturnChannel, messageChanel)

	jobManager.Start(ch)
}
