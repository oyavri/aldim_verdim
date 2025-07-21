package frontend

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	Hostname           string
	DbConnectionString string
	Brokers            []string
	Topic              string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return Config{}, err
	}

	hostname := os.Getenv("HOSTNAME")
	port := os.Getenv("PORT")
	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	broker := os.Getenv("BROKERS")
	brokers := make([]string, 1)
	brokers = append(brokers, broker) // NEED TO REFACTOR
	topic := os.Getenv("TOPIC")

	return Config{
		Hostname:           hostname,
		Port:               port,
		DbConnectionString: dbConnStr,
		Brokers:            brokers,
		Topic:              topic,
	}, nil
}
