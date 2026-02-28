package scheduler

import (
	"JobOrch/internal/platform/workerpool"
	"JobOrch/internal/storage"
	"context"
	"log"
	"time"
)

type Scheduler struct {
	repo         storage.JobRepository
	pool         *workerpool.Pool
	batchSize    int
	pollInterval time.Duration
}

func NewScheduler(repo storage.JobRepository, pool *workerpool.Pool, batchSize int, pollInterval time.Duration) *Scheduler {
	return &Scheduler{
		repo:         repo,
		pool:         pool,
		batchSize:    batchSize,
		pollInterval: pollInterval,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			jobs, err := s.repo.GetPending(ctx, s.batchSize)
			if err != nil {
				log.Printf("%+v\n", err)
				continue
			}
			if len(jobs) == 0 {
				continue
			}
			for i := range jobs {
				s.pool.Submit(ctx, &jobs[i])
			}
		}
	}
}
