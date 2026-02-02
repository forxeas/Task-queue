package service

import (
	"context"
	"task-queue/internal/queue/repository"
	"task-queue/internal/queue/repository/models"
)

type Worker struct {
	repo repository.Repository
	ch   <-chan models.Jobs
}

func (w Worker) Start(ctx context.Context) {

}

/*
Нужно чтобы воркер брал данные из канала и изменял статус на выполнено, сделать такой рандом чтобы она могла ломаться
и я мог проверять работает ли она корректно
*/

/*
Нужно сделать диспатчер который будет брать данные из бд и класть в канал. Там должен быть селект, который проверяет
контекст и дергает таймер. Таймер дергает бд и берет от туда данные.
*/
