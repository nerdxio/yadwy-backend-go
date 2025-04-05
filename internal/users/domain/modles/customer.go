package modles

type Customer struct {
	id int64
}

func NewCustomer(id int64) *Customer {
	return &Customer{
		id: id,
	}
}

func (c *Customer) ID() int64 {
	return c.id
}
