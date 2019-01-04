package server

import (
	"encoding/json"
	"hexagonal-arch-go/domain/ticket"
	"net/http"

	"github.com/go-chi/chi"
)

type ticketHandler struct {
	s ticket.Service
}

func (h *ticketHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.get)
	r.Post("/", h.create)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.getById)
	})

	return r
}

func (h *ticketHandler) get(w http.ResponseWriter, r *http.Request) {
	tickets, _ := h.s.FindAllTickets()

	response, _ := json.Marshal(tickets)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *ticketHandler) getById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ticket, _ := h.s.FindTicketById(id)

	response, _ := json.Marshal(ticket)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *ticketHandler) create(w http.ResponseWriter, r *http.Request) {

	var ticket ticket.Ticket
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&ticket)
	_ = h.s.CreateTicket(&ticket)

	response, _ := json.Marshal(ticket)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
