package server

import (
	"hexagonal-arch-go/domain/ticket"
	"net/http"

	"github.com/go-chi/chi"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	Ticket ticket.Service

	router chi.Router
}

// New returns a new HTTP server.
func New(ts ticket.Service) *Server {
	s := &Server{
		Ticket: ts,
	}

	r := chi.NewRouter()

	r.Use(accessControl)

	r.Route("/tickets", func(r chi.Router) {
		h := ticketHandler{s.Ticket}
		r.Mount("/", h.router())
	})

	s.router = r

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
