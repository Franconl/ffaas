package httpapi

import (
	"net/http"

	"github.com/Franconl/ffaas/internal/repo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(store repo.Flags) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("ok"))
	})

	handlerAdmin := NewAdminHandler(store)

	handlerSdk := NewSdkHandler(store)

	r.Post("/flags", handlerAdmin.Create)

	r.Get("/flags", handlerAdmin.List)

	r.Get("/flags/{id}", handlerAdmin.GetByID)

	r.Get("/flags/key/{key}", handlerAdmin.GetByKey)

	r.Delete("/flags/{id}", handlerAdmin.DeleteByID)

	r.Put("/flags/{id}", handlerAdmin.Update)

	r.Get("/sdk/flags", handlerSdk.List)

	r.Get("/sdk/eval", handlerSdk.Eval)

	return r
}
