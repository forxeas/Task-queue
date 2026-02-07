package service

import (
	"context"
	"errors"
	"task-queue/internal/queue/repository"
	"task-queue/internal/queue/repository/models"
	"time"
)

type Worker struct {
	repo repository.Repository
	ch   <-chan models.Jobs
}

func NewWorker(repo repository.Repository, ch chan models.Jobs) *Worker {
	return &Worker{repo: repo, ch: ch}
}

func (w Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-w.ch:
			if err := w.HandleJob(int64(*job.Id)); err == nil {
				if err := w.repo.MarkJobSuccess(ctx, int64(*job.Id)); err != nil {
					panic(err)
				}
				return
			}

			attempts := job.Attempts + 1

			if attempts >= job.MaxAttempts {
				if err := w.repo.MarkJobFailed(ctx, int64(*job.Id)); err != nil {
					panic(err)
				}
			}

			retry := job.AvailableAt.Add(10 * time.Second)

			if err := w.repo.MarkJobRetry(ctx, int64(*job.Id), attempts, retry); err != nil {
				panic(err)
			}
		}
	}
}

func (w Worker) HandleJob(id int64) error {
	if id%3 == 0 {
		return errors.New("failed to job")
	}

	return nil
}

/*
Нужно чтобы воркер брал данные из канала и изменял статус на выполнено, сделать такой рандом чтобы она могла ломаться
и я мог проверять работает ли она корректно
*/

/*
Нужно сделать диспатчер который будет брать данные из бд и класть в канал. Там должен быть селект, который проверяет
контекст и дергает таймер. Таймер дергает бд и берет от туда данные.
*/
