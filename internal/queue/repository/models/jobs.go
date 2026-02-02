package models

import (
	"encoding/json"
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
	Id          int
	Name        string
	Status      Status
	Payload     json.RawMessage
	Attempts    int
	MaxAttempts int
	AvailableAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
