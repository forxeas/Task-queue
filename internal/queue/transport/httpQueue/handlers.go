package httpQueue

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"task-queue/internal/core/helper"
	"task-queue/internal/queue/repository"
	"task-queue/internal/queue/transport/dto/request"
	"task-queue/internal/queue/transport/dto/response"

	"github.com/gorilla/mux"
)

type Handler struct {
	Ctx  context.Context
	Repo repository.Repository
}

func NewHandler(ctx context.Context, Repo *repository.Repository) *Handler {
	return &Handler{Ctx: ctx, Repo: *Repo}
}

func (h *Handler) RegisterRoute(r *mux.Router) {
	jobs := r.PathPrefix("/jobs").Subrouter()

	jobs.HandleFunc("", h.AddNewTask)
	jobs.HandleFunc("/{id}", h.GetJob)
}

func (h *Handler) AddNewTask(w http.ResponseWriter, r *http.Request) {
	var job request.JobsDTO

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		slog.Warn(err.Error())
		return
	}

	reqJob, err := h.Repo.CreateJob(h.Ctx, job)

	if err != nil {
		slog.Warn(err.Error())
		return
	}

	resDTO := response.NewJobsResponse(*reqJob)

	if err := helper.WriteJson(w, resDTO, 201); err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		if err := helper.WriteJsonError(w, err, 500); err != nil {
			slog.Warn(err.Error())
			return
		}
	}

	job, err := h.Repo.SelectFromId(h.Ctx, int64(id))

	if err != nil {
		if err := helper.WriteJsonError(w, err, 500); err != nil {
			slog.Warn(err.Error())
			return
		}
	}

	if err := helper.WriteJson(w, job, 200); err != nil {
		slog.Warn(err.Error())
		return
	}
}
