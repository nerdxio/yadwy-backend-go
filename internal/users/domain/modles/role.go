package modles

type Role string

const (
	RoleCustomer Role = "CUSTOMER"
	RoleAdmin    Role = "ADMIN"
	RoleSeller   Role = "SELLER"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleCustomer, RoleAdmin, RoleSeller:
		return true
	}
	return false
}

func (r Role) String() string {
	return string(r)
}
