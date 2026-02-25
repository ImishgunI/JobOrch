package storage

import (
	j "JobOrch/internal/job"
	"context"
)

type JobRepository interface {
	Create(ctx context.Context, record j.JobRecord) error
	UpdateStatus(ctx context.Context, id string, status j.Status) error
	SetError(ctx context.Context, id, pg_err string) error
	GetPending(ctx context.Context, limit int) ([]j.JobRecord, error) // Возвращает задачи со статусом pending при падении приложения
}
