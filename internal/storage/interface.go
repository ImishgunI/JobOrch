package storage

import (
	"JobOrch/internal/job"
	"context"

	"github.com/google/uuid"
)

type JobRepository interface {
	Create(ctx context.Context, record *job.JobRecord) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status job.Status) error
	SetError(ctx context.Context, id uuid.UUID, pg_err string) error
	GetPending(ctx context.Context, limit int) ([]job.JobRecord, error) // Возвращает задачи со статусом pending при падении приложения
}
