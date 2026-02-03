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
        	RETURNING id, type, payload, attempts, max_attempts, available_at, created_at, updated_at`

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
			RETURNING id, type, payload, attempts, max_attempts, available_at, created_at, updated_at`

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
	sql := `SELECT * FROM jobs 
         	WHERE status = $1 AND available_at <= now()
         	ORDER BY created_at
         	LIMIT 100
         	FOR UPDATE SKIP LOCKED`

	rows, err := r.Db.Conn.Query(ctx, sql, models.StatusPending)

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

	return jobs, nil
}

func (r *Repository) MarkJobSuccess(ctx context.Context, id int64) error {
	sql := `UPDATE jobs SET status = $1, updated_at = now() WHERE id = $2`

	cmd, err := r.Db.Conn.Exec(ctx, sql, models.StatusDone, id)
	return checkErr(cmd, err)
}

func (r *Repository) MarkJobFailed(ctx context.Context, id int64) error {
	sql := `UPDATE jobs SET status = $1 WHERE id = $2`

	cmd, err := r.Db.Conn.Exec(ctx, sql, models.StatusFailed, id)
	return checkErr(cmd, err)
}

func (r *Repository) MarkJobRetry(ctx context.Context, id int64, attempts int, availableAt time.Time) error {
	sql := `UPDATE jobs SET attempts = $1, available_at = $2 WHERE id = $3`

	cmd, err := r.Db.Conn.Exec(ctx, sql, attempts, availableAt, id)
	return checkErr(cmd, err)
}
