package main

import (
	"JobOrch/adapters/postgres"
	"JobOrch/internal/platform/workerpool"
	"JobOrch/internal/scheduler"
	"context"
	"log"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	pool := workerpool.NewPool(ctx, 10)
	repo, err := postgres.NewPool(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Printf("%+v\n", err)
		cancel()
	}
	duration := time.Duration(500) * time.Millisecond
	s := scheduler.NewScheduler(repo, pool, 10, duration)
	go func() {
		for result := range pool.Results() {
			postgres.Updater(ctx, repo, result)
		}
	}()
	s.Run(ctx)
	cancel()
	pool.ShutDown()
}
