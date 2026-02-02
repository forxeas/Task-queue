package repository

import (
	"context"
	"task-queue/internal/db"
	"task-queue/internal/queue/repository/models"
	"task-queue/internal/queue/transport/dto/request"
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

// Нужно добавить селект для того чтобы диспатчер клал их в канал для воркеров
