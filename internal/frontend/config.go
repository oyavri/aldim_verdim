package frontend

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	Hostname           string
	DbConnectionString string
	Broker             string
	BrokerTopic        string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return Config{}, err
	}

	hostname := os.Getenv("FIBER_HOSTNAME")
	port := os.Getenv("FIBER_PORT")
	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	brokerHostname := os.Getenv("KAFKA_HOSTNAME")
	brokerPort := os.Getenv("KAFKA_PORT")
	broker := fmt.Sprintf("%s:%s", brokerHostname, brokerPort)
	brokerTopic := os.Getenv("KAFKA_TOPIC")

	return Config{
		Hostname:           hostname,
		Port:               port,
		DbConnectionString: dbConnStr,
		Broker:             broker,
		BrokerTopic:        brokerTopic,
	}, nil
}
