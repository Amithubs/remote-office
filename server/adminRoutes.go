package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/tejashwikalptaru/remote-office/handler"
	"github.com/tejashwikalptaru/remote-office/middleware"
)

func adminRoutes(r chi.Router) {
	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AuthMiddleware)
		admin.Get("/", handler.AdminRights)
	})
}
