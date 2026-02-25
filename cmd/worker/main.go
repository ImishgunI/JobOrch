package main

import (
	"JobOrch/adapters/postgres"
	"JobOrch/internal/platform/workerpool"
	"context"
	"log"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	pool := workerpool.NewPool(ctx, 10)
	repo, err := postgres.NewPool(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Printf("%+v\n", err)
		cancel()
	}
	for {
		jobs, err := repo.GetPending(ctx, 10)
		if err != nil {
			log.Printf("%+v\n", err)
			break
		}
		if len(jobs) == 0 {
			break
		}
		for i := range jobs {
			pool.Submit(ctx, &jobs[i])
		}
	}
	cancel()
	pool.ShutDown()
}
