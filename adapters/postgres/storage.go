package postgres

import (
	j "JobOrch/internal/job"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	pool *pgxpool.Pool
}

func NewPool(ctx context.Context, connection string) (*postgresRepo, error) {
	conn, err := pgxpool.New(ctx, connection)
	if err != nil {
		return nil, err
	}
	return &postgresRepo{
		pool: conn,
	}, nil
}

func (r *postgresRepo) Create(ctx context.Context, record *j.JobRecord) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO jobs (id, status, payload, error) VALUES ($1, $2, $3, $4)`,
		record.ID, record.Status, record.Payload, record.Error)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status j.Status) error {
	_, err := r.pool.Exec(ctx, `
	UPDATE jobs
	SET status = $1
	WHERE id = $2
	`, status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresRepo) SetError(ctx context.Context, id uuid.UUID, pg_err string) error {
	_, err := r.pool.Exec(ctx, `
	UPDATE jobs
	SET error = $1
	WHERE id = $2
	`, pg_err, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresRepo) GetPending(ctx context.Context, limit int) ([]j.JobRecord, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
	SELECT id, status, payload, error, created_at, updated_at FROM jobs
	WHERE status = $1
	ORDER BY created_at
	FOR UPDATE SKIP LOCKED
	LIMIT $2`, j.PENDING, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var record []j.JobRecord
	var ids []uuid.UUID
	for rows.Next() {
		var rec j.JobRecord
		err := rows.Scan(
			&rec.ID,
			&rec.Status,
			&rec.Payload,
			&rec.Error,
			&rec.CreatedAt,
			&rec.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		record = append(record, rec)
		ids = append(ids, rec.ID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		_, err := tx.Exec(ctx, `
		UPDATE jobs
		SET status = $1
		WHERE id = ANY($2)`, j.PROCESSING, ids)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return record, nil
}
