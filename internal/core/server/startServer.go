package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Start(ctx context.Context, addr string, router mux.Router) error {
	server := &http.Server{
		Addr: addr, Handler: &router,
	}

	go func() {
		<-ctx.Done()
		ctxShotDown, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		_ = server.Shutdown(ctxShotDown)
	}()

	return server.ListenAndServe()
}
