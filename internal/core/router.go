package core

import (
	"task-queue/internal/queue/transport/httpQueue"

	"github.com/gorilla/mux"
)

func NewRouter(hQueue *httpQueue.Handler) *mux.Router {
	router := mux.NewRouter()

	hQueue.RegisterRoute(router)

	return router
}
