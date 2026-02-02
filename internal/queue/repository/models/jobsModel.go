package models

import (
	"encoding/json"
	"task-queue/internal/queue/transport/dto/request"
	"time"
)

type Status string

const (
	StatusPending      Status = "pending"
	StatusFailed       Status = "failed"
	StatusDone         Status = "done"
	StatusInProgressed Status = "InProgressed"
)

type Jobs struct {
	Id          *int
	Type        string
	Status      Status
	Payload     json.RawMessage
	Attempts    int
	MaxAttempts int
	AvailableAt time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func NewJobs(jobs request.JobsDTO) *Jobs {
	return &Jobs{
		Id:          nil,
		Type:        jobs.Type,
		Status:      StatusPending,
		Payload:     jobs.Payload,
		Attempts:    0,
		MaxAttempts: *jobs.MaxAttempts,
		AvailableAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
}

func (j *Jobs) GetPayload(v *any) error {
	return json.Unmarshal(j.Payload, &v)

}
