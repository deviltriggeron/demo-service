package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	e "demo-service/internal/entity"
)

func LoadConfigDB() *e.ConfigDB {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &e.ConfigDB{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		ServerPort:       os.Getenv("SERVER_PORT"),
	}
}

func LoadConfigKafka() *e.ConfigBroker {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &e.ConfigBroker{
		Broker:  os.Getenv("BROKER"),
		GroupID: os.Getenv("GROUP_ID"),
	}
}

func LoadConfigAddr() string {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return os.Getenv("ADDR")
}
