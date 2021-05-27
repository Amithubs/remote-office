package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/tejashwikalptaru/remote-office/handler"
	"github.com/tejashwikalptaru/remote-office/middleware"
)

// userRoutes as routes related to users with needs authorized token
func userRoutes(r chi.Router) {
	r.Group(func(user chi.Router) {
		user.Use(middleware.AuthMiddleware)
		user.Get("/", handler.GetUserInfo)
		user.Put("/", handler.UpdateUser)
		user.Put("/logout",handler.LogoutUser)
		user.Post("/upload",handler.UploadProfileImage)
	})
}
