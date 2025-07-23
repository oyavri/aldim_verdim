package worker

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oyavri/aldim_verdim/pkg/db"
	"github.com/oyavri/aldim_verdim/pkg/entity"
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

	repository := NewWalletRepository(dbPool)
	service := NewEventService(repository)

	// Graceful shutdown
	go func(dbPool *pgxpool.Pool, kc *KafkaConsumer, cancelFunc context.CancelFunc) {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		signal.Stop(signalChan)

		log.Println("Consumer is closing")
		kc.Close()

		log.Println("Database connection is closing")
		dbPool.Close()

		cancelFunc()
	}(dbPool, kafkaConsumer, cancel)

	log.Println("Worker started running")
	for {
		e, err := kafkaConsumer.Consume(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Context cancelled, stopping worker")
				return
			}

			log.Printf("Error consuming event: %v", err)
			continue
			// Retry mechanism might be helpful here, if it is possible to save the event
		}

		go func(event entity.Event) {
			err := service.HandleEvent(ctx, event)
			if err != nil {
				log.Printf("Failed to handle event %v with the error: %v", event, err)
			}
		}(e)
	}
}
