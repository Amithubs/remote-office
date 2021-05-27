package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/tejashwikalptaru/remote-office/handler"
	"github.com/tejashwikalptaru/remote-office/middleware"
)

func standupRoutes(r chi.Router) {
	r.Group(func(standup chi.Router) {
		standup.Use(middleware.AuthMiddleware)
		standup.Post("/", handler.CreateStandUp)
		standup.Put("/{id}", handler.ModifyStandUp)
		standup.Get("/", handler.GetStandUp)
	})
}
