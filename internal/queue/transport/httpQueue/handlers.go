package httpQueue

import (
	"context"
	"task-queue/internal/queue/service"
)

type Handler struct {
	Worker service.Worker
}

func NewHandler(worker service.Worker) *Handler {
	return &Handler{Worker: worker}
}

func (h Handler) AddQueueHandler(ctx context.Context) {

}
