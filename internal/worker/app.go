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

	dbPool, err := db.NewPostgresPool(ctx, cfg.DbConnectionString)
	if err != nil {
		log.Fatal("Error creating database pool")
	}

	// Since I am using single broker, I will wrap the only broker into a slice
	kafkaConsumer := NewKafkaConsumer(cfg.ConsumerGroupId, []string{cfg.Broker}, cfg.BrokerTopic)

	repository := NewWalletRepository(dbPool)
	service := NewEventService(repository, cfg.MaxGoroutineCount)

	// Graceful shutdown
	go func(dbPool *pgxpool.Pool, kc *KafkaConsumer, cancelFunc context.CancelFunc) {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		signal.Stop(signalChan)

		cancelFunc()

		log.Println("Consumer is closing")
		kc.Close()

		log.Println("Database connection is closing")
		dbPool.Close()

		cancelFunc()
	}(dbPool, kafkaConsumer, cancel)

	log.Println("Worker started running")
	for {
		if ctx.Err() != nil {
			log.Println("Context cancelled, stopping worker")
			return
		}

		event, err := kafkaConsumer.Consume(ctx)
		if err != nil {
			log.Printf("Error consuming event: %v", err)
			continue
		}

		go func(ctx context.Context, e entity.Event) {
			event := e
			log.Printf("Handling event: %v", e)
			service.HandleEvent(ctx, event)
		}(ctx, event)
	}
}
