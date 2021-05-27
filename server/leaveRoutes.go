package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/tejashwikalptaru/remote-office/handler"
	"github.com/tejashwikalptaru/remote-office/middleware"
)

func leaveRoutes(r chi.Router) {
	r.Group(func(leave chi.Router) {
		leave.Use(middleware.AuthMiddleware)
		leave.Post("/", handler.CreateLeave)
		leave.Put("/{id}", handler.ModifyLeave)
		leave.Get("/", handler.GetLeave)
		leave.Get("/", handler.GetLeaveStats)
	})
}
