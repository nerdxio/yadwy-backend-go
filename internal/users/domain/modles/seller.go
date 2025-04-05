package modles

type Seller struct {
	id     int64
	userId int64
}

func NewSeller(id, userId int64) *Seller {
	return &Seller{
		id:     id,
		userId: userId,
	}
}

func (s *Seller) ID() int64 {
	return s.id
}

func (s *Seller) UserId() int64 {
	return s.userId
}
