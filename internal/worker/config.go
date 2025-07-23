package worker

import (
	"fmt"
	"os"
)

type Config struct {
	DbConnectionString string
	ConsumerGroupId    string
	Broker             string
	BrokerTopic        string
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

	return Config{
		DbConnectionString: dbConnStr,
		ConsumerGroupId:    groupId,
		Broker:             broker,
		BrokerTopic:        brokerTopic,
	}, nil
}
