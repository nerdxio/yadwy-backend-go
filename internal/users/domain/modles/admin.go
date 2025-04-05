package modles

type Admin struct {
	id int64
}

func NewAdmin(id int64) *Admin {
	return &Admin{
		id: id,
	}
}
