package mongo

import (
	"hexagonal-arch-go/domain/ticket"

	"github.com/globalsign/mgo"
)

const DB_NAME = "ticketing"
const COLLECTION = "tickets"

type ticketRepository struct {
	session *mgo.Session
}

func NewMongoTicketRepository(session *mgo.Session) ticket.Repository {
	return &ticketRepository{
		session,
	}
}

func (r *ticketRepository) Create(ticket *ticket.Ticket) error {
	err := r.session.DB(DB_NAME).C(COLLECTION).Insert(ticket)

	if err != nil {
		return err
	}

	return nil
}

func (r *ticketRepository) FindById(id string) (*ticket.Ticket, error) {
	t := new(ticket.Ticket)
	err := r.session.DB(DB_NAME).C(COLLECTION).FindId(t.ID).One(t)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *ticketRepository) FindAll() (tickets []*ticket.Ticket, err error) {
	err = r.session.DB(DB_NAME).C(COLLECTION).Find(nil).All(tickets)

	if err != nil {
		return nil, err
	}

	return tickets, nil
}
