package redis

import (
	"encoding/json"
	"hex-arch/domain/ticket"

	"github.com/go-redis/redis"
)

const table = "tickets"

type ticketRepository struct {
	connection *redis.Client
}

func NewRedisTicketRepository(connection *redis.Client) ticket.Repository {
	return &ticketRepository{
		connection,
	}
}

func (r *ticketRepository) Create(ticket *ticket.Ticket) error {
	encoded, err := json.Marshal(ticket)

	if err != nil {
		return err
	}

	r.connection.HSet(table, ticket.ID, encoded) //Don't expire
	return nil
}

func (r *ticketRepository) FindById(id string) (*ticket.Ticket, error) {
	b, err := r.connection.HGet(table, id).Bytes()

	if err != nil {
		return nil, err
	}

	t := new(ticket.Ticket)
	err = json.Unmarshal(b, t)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *ticketRepository) FindAll() (tickets []*ticket.Ticket, err error) {
	ts := r.connection.HGetAll(table).Val()
	for key, value := range ts {
		t := new(ticket.Ticket)
		err = json.Unmarshal([]byte(value), t)

		if err != nil {
			return nil, err
		}

		t.ID = key
		tickets = append(tickets, t)
	}
	return tickets, nil
}
