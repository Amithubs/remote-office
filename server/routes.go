package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/tejashwikalptaru/remote-office/handler"
	"github.com/tejashwikalptaru/remote-office/middleware"
	"net/http"
)

type Server struct {
	chi.Router
}

// SetupRoutes provides all the routes that can be used
func SetupRoutes () *Server {
	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.CommonMiddlewares()...)

		// health endpoint
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		// unauthorized routes, ie. public routes
		r.Post("/signup", handler.SignUp)
		r.Post("/login", handler.LoginUser)

		r.Route("/user",  func(user chi.Router) {
			user.Group(userRoutes)
		})
		r.Route("/leave",  func(leave chi.Router) {
			leave.Group(leaveRoutes)
		})
		r.Route("/standup",  func(standup chi.Router) {
			standup.Group(standupRoutes)
		})
		r.Route("/admin",  func(admin chi.Router) {
			admin.Group(adminRoutes)
		})
	})
	return &Server{Router: router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
