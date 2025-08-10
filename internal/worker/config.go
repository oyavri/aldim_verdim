package worker

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DbConnectionString string
	ConsumerGroupId    string
	Broker             string
	BrokerTopic        string
	MaxGoroutineCount  int
}

func LoadConfig() (Config, error) {
	// Disabled due to containerized server

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// 	return Config{}, err
	// }

	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	brokerHostname := os.Getenv("KAFKA_HOSTNAME")
	brokerPort := os.Getenv("KAFKA_PORT")
	broker := fmt.Sprintf("%s:%s", brokerHostname, brokerPort)
	brokerTopic := os.Getenv("KAFKA_TOPIC")
	groupId := os.Getenv("CONSUMER_GROUP_ID")
	maxGoRoutines, err := strconv.Atoi(os.Getenv("MAX_GOROUTINE_COUNT"))

	if err != nil {
		return Config{}, err
	}

	return Config{
		DbConnectionString: dbConnStr,
		ConsumerGroupId:    groupId,
		Broker:             broker,
		BrokerTopic:        brokerTopic,
		MaxGoroutineCount:  maxGoRoutines,
	}, nil
}
