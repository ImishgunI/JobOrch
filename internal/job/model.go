package job

import (
	"context"
	"time"
)

type Job interface {
	Execute(ctx context.Context) error
}

type JobRecord struct {
	ID        string
	Status    Status
	Payload   []byte
	Error     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
