package workerpool

import (
	j "JobOrch/internal/job"
	"context"
	"log"
	"runtime/debug"
	"sync"
)

type Pool struct {
	taskQueue chan j.Job
	wg        *sync.WaitGroup
	pool_size int
}

func NewPool(ctx context.Context, maxConcurrency int) *Pool {
	job := make(chan j.Job, 10)
	var wg sync.WaitGroup
	p := &Pool{
		taskQueue: job,
		wg:        &wg,
		pool_size: maxConcurrency,
	}
	p.wg.Add(p.pool_size)
	for range p.pool_size {
		go func() {
			defer p.wg.Done()
			for {
				select {
				case job, ok := <-p.taskQueue:
					if !ok {
						return
					}
					p.worker(ctx, job)
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	return p
}

func (p *Pool) Submit(ctx context.Context, job j.Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.taskQueue <- job:
		return nil
	}
}

func (p *Pool) worker(ctx context.Context, task j.Job) <-chan j.JobResult {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%+v\n", r)
			log.Println(string(debug.Stack()))
		}
	}()
	err := task.Execute(ctx)
	if err != nil {
		log.Printf("%+v", err)
	}
	rec := task.(*j.JobRecord)
	res := j.JobResult{
		JobID: rec.ID,
		Err:   err,
	}
	JobCh := make(chan j.JobResult, 1)
	JobCh <- res
	return JobCh
}

func (p *Pool) ShutDown() {
	close(p.taskQueue)
	p.wg.Wait()
}
