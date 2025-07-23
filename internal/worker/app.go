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

	ctx := context.Background()
	dbPool, err := db.NewPostgresPool(ctx, cfg.DbConnectionString)
	if err != nil {
		log.Fatal("Error creating database pool")
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	signal.Stop(signalChan)

	log.Println("Database connection is closing")
	dbPool.Close()
}
