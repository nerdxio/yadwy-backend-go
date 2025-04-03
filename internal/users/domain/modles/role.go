package modles

import "fmt"

type Role string

const (
	RoleCustomer Role = "CUSTOMER"
	RoleAdmin    Role = "ADMIN"
	RoleSeller   Role = "SELLER"
)

func NewRole(roleStr string) (Role, error) {
	role := Role(roleStr)
	if !role.IsValid() {
		return "", fmt.Errorf("invalid role: %s. Must be one of: CUSTOMER, ADMIN, SELLER", roleStr)
	}
	return role, nil
}

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
