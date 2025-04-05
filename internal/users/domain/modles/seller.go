package modles

type Seller struct {
	id int64
}

func NewSeller(id int64) *Seller {
	return &Seller{
		id: id,
	}
}

func (s *Seller) ID() int64 {
	return s.id
}
