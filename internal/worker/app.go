package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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

	// Pool size should be equal (or greater than) the maximum goroutine count
	connString := fmt.Sprintf("%v?pool_max_conns=%v", cfg.DbConnectionString, cfg.MaxGoroutineCount)

	dbPool, err := db.NewPostgresPool(ctx, connString)
	if err != nil {
		log.Fatal("Error creating database pool")
	}

	// Since I am using single broker, I will wrap the only broker into a slice
	kafkaConsumer := NewKafkaConsumer(cfg.ConsumerGroupId, []string{cfg.Broker}, cfg.BrokerTopic)

	repository := NewWalletRepository(dbPool)
	service := NewEventService(repository)

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
	}(dbPool, kafkaConsumer, cancel)

	var wg sync.WaitGroup
	concurrentGoroutineCount := make(chan struct{}, cfg.MaxGoroutineCount)

	log.Println("Worker started running")
	for {
		if ctx.Err() != nil {
			log.Println("Context cancelled, stopping worker after jobs are done.")
			wg.Wait()
			log.Println("All of the jobs are done, shutting worker down.")
			return
		}

		event, err := kafkaConsumer.Consume(ctx)
		if err != nil {
			// In case the error is due to context cancellation
			if ctx.Err() != nil {
				continue
			}

			log.Printf("Error consuming event: %v", err)
			continue
		}

		wg.Add(1)
		concurrentGoroutineCount <- struct{}{}

		go func(ctx context.Context, event entity.Event) {
			defer func() {
				<-concurrentGoroutineCount
				wg.Done()
			}()
			log.Printf("Handling event: %v", event)
			service.HandleEvent(ctx, event)
		}(ctx, event)
	}
}
