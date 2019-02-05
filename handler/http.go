package handler

import (
	"net/http"
	"time"

	"github.com/mochisuna/load-test-sample/domain/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/render"
	"gopkg.in/go-playground/validator.v9"
)

var rendering = render.New(render.Options{})
var validate = validator.New()

// Services is grouping application services structure
type Services struct {
	UserService service.UserService
}

// Server HTTP server
type Server struct {
	http.Server
	Services
	RedirectURL string //これイマイチ
}

// New inject to domain services
func New(addr string, services *Services, redirectURL string) *Server {
	return &Server{
		Server: http.Server{
			Addr: addr,
		},
		Services:    *services,
		RedirectURL: redirectURL,
	}
}

// ListenAndServe override http ListenAndServe
func (s *Server) ListenAndServe() error {
	r := chi.NewRouter()

	// cord option
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	// 公式提供のmiddleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// routings
	r.Route("/v1", func(r chi.Router) {
		// User
		r.Route("/users", func(r chi.Router) {
			r.Post("/", s.registerUser)
			r.Get("/{userID}", s.referUser)
		})
		// display
		r.Route("/display", func(r chi.Router) {
			r.Post("/", s.displayUser)
		})
	})

	s.Handler = r
	return s.Server.ListenAndServe()
}
