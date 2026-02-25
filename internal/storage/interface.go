package storage

import (
	j "JobOrch/internal/job"
	"context"
)

type JobRepository interface {
	Create(ctx context.Context, job j.Job) error
	UpdateStatus(ctx context.Context, id string, status j.Status) error
	SetError(ctx context.Context, id, err string) error
	GetPending(ctx context.Context, limit int) ([]j.Job, error) // Возвращает задачи со статусом pending при падении приложения
}
