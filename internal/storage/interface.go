package storage

import (
	j "JobOrch/internal/job"
	"context"

	"github.com/google/uuid"
)

type JobRepository interface {
	Create(ctx context.Context, record *j.JobRecord) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status j.Status) error
	SetError(ctx context.Context, id uuid.UUID, pg_err string) error
	GetPending(ctx context.Context, limit int) ([]j.JobRecord, error) // Возвращает задачи со статусом pending при падении приложения
}
