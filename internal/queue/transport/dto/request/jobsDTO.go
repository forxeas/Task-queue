package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

type JobsDTO struct {
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	MaxAttempts *int            `json:"max_attempts"`
}

func NewJobs(
	typePayload string,
	payload json.RawMessage,
	maxAttempts *int,
) *JobsDTO {
	return &JobsDTO{
		Type:        typePayload,
		Payload:     payload,
		MaxAttempts: maxAttempts,
	}
}

func (j *JobsDTO) Validate() error {
	trimmedPayload := bytes.TrimSpace(j.Payload)

	if strings.Trim(j.Type, " ") == "" {
		return errors.New("job type is missing")
	}

	if len(trimmedPayload) == 0 {
		return errors.New("job payload is empty")
	}

	return nil
}
