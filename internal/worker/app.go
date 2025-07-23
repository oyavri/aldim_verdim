package worker

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oyavri/aldim_verdim/pkg/db"
)

func Run() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("Error loading config")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := db.NewPostgresPool(ctx, cfg.DbConnectionString)
	if err != nil {
		log.Fatal("Error creating database pool")
	}

	// Since I am using single broker, I will wrap the only broker into a slice
	kafkaConsumer := NewKafkaConsumer([]string{cfg.Broker}, cfg.BrokerTopic)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	signal.Stop(signalChan)

	log.Println("Consumer is closing")
	kafkaConsumer.Close()

	log.Println("Database connection is closing")
	dbPool.Close()
}
