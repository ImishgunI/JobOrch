package workerpool

import (
	j "JobOrch/internal/job"
	"context"
	"log"
	"runtime/debug"
	"sync"
)

type Pool struct {
	taskQueue   chan j.Job
	resultQueue chan j.JobResult
	wg          *sync.WaitGroup
	pool_size   int
}

func NewPool(ctx context.Context, maxConcurrency int) *Pool {
	job := make(chan j.Job, 10)
	result := make(chan j.JobResult, 10)
	var wg sync.WaitGroup
	p := &Pool{
		taskQueue:   job,
		resultQueue: result,
		wg:          &wg,
		pool_size:   maxConcurrency,
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

func (p *Pool) worker(ctx context.Context, task j.Job) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%+v\n", r)
			log.Println(string(debug.Stack()))
		}
	}()
	result := task.Execute(ctx)
	p.resultQueue <- (<-result)
}

func (p *Pool) Results() <-chan j.JobResult {
	return p.resultQueue
}

func (p *Pool) ShutDown() {
	close(p.taskQueue)
	p.wg.Wait()
	close(p.resultQueue)
}
