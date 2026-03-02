package job

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Job interface {
	Execute(ctx context.Context) chan JobResult
}

type JobRecord struct {
	ID        uuid.UUID `db:"id"`
	Status    Status    `db:"status"`
	Payload   []byte    `db:"payload"`
	Error     *string   `db:"error"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type JobResult struct {
	JobID uuid.UUID
	Err   error
}

func (j *JobRecord) Execute(ctx context.Context) chan JobResult
