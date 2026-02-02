package response

import (
	"encoding/json"
	"task-queue/internal/queue/repository/models"
	"time"
)

type JobsDTO struct {
	Id          int             `json:"id"`
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Status      string          `json:"status"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`
	AvailableAt time.Time       `json:"available_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at"`
}

func NewJobsResponse(job models.Jobs) *JobsDTO {
	return &JobsDTO{
		Id:          job.Id,
		Type:        job.Type,
		Payload:     job.Payload,
		Status:      string(job.Status),
		Attempts:    job.Attempts,
		MaxAttempts: job.MaxAttempts,
		AvailableAt: job.AvailableAt,
		CreatedAt:   job.CreatedAt,
		UpdatedAt:   job.UpdatedAt,
	}
}
