package service

import (
	"context"
	"log/slog"
	"task-queue/internal/queue/repository"
	"task-queue/internal/queue/repository/models"
	"time"
)

type Dispatcher struct {
	Ch   chan models.Jobs
	Repo *repository.Repository
}

func NewDispatcher(ch chan models.Jobs, repo *repository.Repository) *Dispatcher {
	return &Dispatcher{Ch: ch, Repo: repo}
}

func (d Dispatcher) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			jobs, err := d.Repo.SelectJobs(ctx)

			if err != nil {
				slog.Warn(err.Error())
				continue
			}

			for _, v := range jobs {
				select {
				case <-ctx.Done():
					return
				case d.Ch <- *v:
				}
			}
		}
	}
}
