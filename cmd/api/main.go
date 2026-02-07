package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"task-queue/internal/core"
	"task-queue/internal/core/server"
	"task-queue/internal/db"
	"task-queue/internal/queue/repository"
	"task-queue/internal/queue/repository/models"
	"task-queue/internal/queue/service"
	"task-queue/internal/queue/transport/httpQueue"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ch := make(chan models.Jobs, 100)
	database, err := db.Connection(ctx)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	repo := repository.NewRepository(database)
	handlerQueue := httpQueue.NewHandler(ctx, repo)

	dispatcher := service.NewDispatcher(ch, repo)
	worker := service.NewWorker(repo, ch)

	dispatcher.Start(ctx)
	worker.Start(ctx)

	router := core.NewRouter(handlerQueue)

	if err := server.Start(ctx, os.Getenv("PORT"), *router); err != nil {
		slog.Error(err.Error())
		return
	}
}
