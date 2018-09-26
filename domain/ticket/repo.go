package ticket

type Repository interface {
	Create(ticket *Ticket) error
	FindById(id string) (*Ticket, error)
	FindAll() ([]*Ticket, error)
}
