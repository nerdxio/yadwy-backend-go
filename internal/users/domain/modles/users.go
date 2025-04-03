package modles

type User struct {
	id       int
	name     string
	email    string
	password string
	role     Role
}

func NewUser(id int, name, email, password string, role Role) *User {
	return &User{
		id:       id,
		name:     name,
		email:    email,
		password: password,
		role:     role,
	}
}

func (u *User) Name() string {
	return u.name
}
func (u *User) ID() int {
	return u.id
}
func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}
func (u *User) Role() Role {
	return u.role
}
