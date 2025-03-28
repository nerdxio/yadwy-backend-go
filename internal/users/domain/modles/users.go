package modles

type User struct {
	id       int
	name     string
	email    string
	password string
	role     Role
}

func NewUser(name, email, password string, role Role) (*User, error) {
	return &User{
		name:     name,
		email:    email,
		password: password,
		role:     role,
	}, nil
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Role() Role {
	return u.role
}
