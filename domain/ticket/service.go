package ticket

import (
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateTicket(ticket *Ticket) error
	FindTicketById(id string) (*Ticket, error)
	FindAllTickets() ([]*Ticket, error)
}

type service struct {
	repo Repository
}

func NewTicketService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) CreateTicket(ticket *Ticket) error {
	ticket.ID = uuid.New().String()
	ticket.Created = time.Now()
	ticket.Updated = time.Now()
	ticket.Status = "open"
	return s.repo.Create(ticket)
}

func (s *service) FindTicketById(id string) (*Ticket, error) {
	return s.repo.FindById(id)
}

func (s *service) FindAllTickets() ([]*Ticket, error) {
	return s.repo.FindAll()
}
