package repository

import (
	"context"
	"task-queue/internal/db"
	"task-queue/internal/queue/repository/models"
	"task-queue/internal/queue/transport/dto/request"
	"time"
)

type Repository struct {
	Db *db.Db
}

func NewRepository(db *db.Db) *Repository {
	return &Repository{Db: db}
}

func (r *Repository) CreateJob(ctx context.Context, dto request.JobsDTO) (*models.Jobs, error) {
	var model models.Jobs

	sql := `INSERT INTO jobs (type, payload, max_attempts) 
			VALUES ($1, $2, $3) 
        	RETURNING id, type, status, payload, attempts, max_attempts, available_at, created_at, updated_at`

	err := r.Db.Conn.QueryRow(ctx, sql, dto.Type, dto.Payload, dto.MaxAttempts).
		Scan(
			&model.Id,
			&model.Type,
			&model.Status,
			&model.Payload,
			&model.Attempts,
			&model.MaxAttempts,
			&model.AvailableAt,
			&model.CreatedAt,
			&model.UpdatedAt,
		)

	return &model, err
}

func (r *Repository) UpdateJob(ctx context.Context, model models.Jobs) (*models.Jobs, error) {
	sql := `UPDATE jobs 
			SET type = $1, 
			    status = $2,
			    payload = $3, 
			    attempts = $4, 
			    max_attempts = $5, 
			    available_at = $6, 
			    updated_at = $7
			WHERE id = $8
			RETURNING id, type, status, payload, attempts, max_attempts, available_at, created_at, updated_at`

	err := r.Db.Conn.QueryRow(
		ctx,
		sql,
		model.Type,
		model.Status,
		model.Payload,
		model.Attempts,
		model.MaxAttempts,
		model.AvailableAt,
	).Scan(
		&model.Id,
		&model.Type,
		&model.Status,
		&model.Payload,
		&model.Attempts,
		&model.MaxAttempts,
		&model.AvailableAt,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	return &model, err
}

func (r *Repository) SelectJobs(ctx context.Context) ([]*models.Jobs, error) {
	jobs := make([]*models.Jobs, 0)
	sql := `WITH picked AS (
				SELECT * FROM jobs 
         		WHERE status = $1 AND available_at <= now()
         		ORDER BY created_at
         		LIMIT 100
         		FOR UPDATE SKIP LOCKED
			)
			UPDATE jobs AS j FROM picked
			SET status = $2, updated_at = now()
			WHERE j.id = picked.id
			RETURNING 
				j.id, j.type, 
				j.status, j.payload,
				j.attempts, 
				j.max_attempts, 
				j.available_at, 
				j.created_at,
				j.updated_at`

	rows, err := r.Db.Conn.Query(ctx, sql, models.StatusPending, models.StatusInProgressed)

	if err != nil {
		return jobs, err
	}
	defer rows.Close()

	for rows.Next() {
		var job models.Jobs

		if err := rows.Scan(
			&job.Id,
			&job.Type,
			&job.Status,
			&job.Payload,
			&job.Attempts,
			&job.MaxAttempts,
			&job.AvailableAt,
			&job.CreatedAt,
			&job.UpdatedAt,
		); err != nil {
			return jobs, err
		}

		jobs = append(jobs, &job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *Repository) SelectFromId(ctx context.Context, id int64) (*models.Jobs, error) {
	var job models.Jobs

	sql := `SELECT id, type, status, payload, created_at, updated_at FROM jobs WHERE id = $1`

	if err := r.Db.Conn.QueryRow(ctx, sql, id).
		Scan(
			&job.Id,
			&job.Type,
			&job.Status,
			&job.Payload,
			&job.CreatedAt,
			&job.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *Repository) MarkJobSuccess(ctx context.Context, id int64) error {
	sql := `UPDATE jobs SET status = $1, updated_at = now() WHERE id = $2`

	cmd, err := r.Db.Conn.Exec(ctx, sql, models.StatusDone, id)
	return checkErr(cmd, err)
}

func (r *Repository) MarkJobFailed(ctx context.Context, id int64) error {
	sql := `UPDATE jobs SET status = $1, updated_at = now() WHERE id = $2`

	cmd, err := r.Db.Conn.Exec(ctx, sql, models.StatusFailed, id)
	return checkErr(cmd, err)
}

func (r *Repository) MarkJobRetry(ctx context.Context, id int64, attempts int, availableAt time.Time) error {
	sql := `UPDATE jobs SET attempts = $1, available_at = $2 WHERE id = $3`

	cmd, err := r.Db.Conn.Exec(ctx, sql, attempts, availableAt, id)
	return checkErr(cmd, err)
}
