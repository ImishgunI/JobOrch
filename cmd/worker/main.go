package main

import (
	"JobOrch/internal/platform/workerpool"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	pool := workerpool.NewPool(ctx, 10)
	
	cancel()
	pool.ShutDown()
}
