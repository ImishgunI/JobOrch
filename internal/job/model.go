package job

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

type Job interface {
	Execute(ctx context.Context) error
}

type JobRecord struct {
	ID        uuid.UUID `db:"id"`
	Status    Status    `db:"status"`
	Payload   []byte    `db:"payload"`
	Error     *string   `db:"error"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (j *JobRecord) Execute(ctx context.Context) error {
	log.Printf("Do work")
	return nil
}
