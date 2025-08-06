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
	kafkaConsumer := NewKafkaConsumer(cfg.ConsumerGroupId, []string{cfg.Broker}, cfg.BrokerTopic)

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
		if ctx.Err() != nil {
			log.Println("Context cancelled, stopping worker")
			return
		}

		events, err := kafkaConsumer.ConsumeBatch(ctx, 50)
		if err != nil {
			log.Printf("Error consuming events: %v", err)
			continue
		}

		e := make(map[string][]entity.Event)
		processedWalletIds := make(chan string)

		for _, event := range events {
			walletTransactions, ok := e[event.WalletId]
			if !ok {
				walletTransactions := []entity.Event{}
				walletTransactions = append(walletTransactions, event)
				e[event.WalletId] = walletTransactions
				continue
			}

			walletTransactions = append(walletTransactions, event)
			e[event.WalletId] = walletTransactions
		}

		for walletId, transactions := range e {
			go func(walletId string, transactions []entity.Event) {
				err := service.HandleEvents(ctx, transactions)
				if err != nil {
					log.Printf("Failed to handle events with the error: %v", err)
				}
				processedWalletIds <- walletId
			}(walletId, transactions)
		}

		if len(processedWalletIds) > 0 {
			for walletId := range processedWalletIds {
				delete(e, walletId)
			}
		}
	}
}
