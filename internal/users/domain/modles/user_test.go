package modles

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			id       int
			name     string
			email    string
			password string
			role     Role
		}
		want User
	}{
		{
			name: "should create new user successfully",
			args: struct {
				id       int
				name     string
				email    string
				password string
				role     Role
			}{
				id:       1,
				name:     "John Doe",
				email:    "john@example.com",
				password: "hashedPassword123",
				role:     RoleCustomer,
			},
			want: User{
				id:       1,
				name:     "John Doe",
				email:    "john@example.com",
				password: "hashedPassword123",
				role:     RoleCustomer,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUser(tt.args.id, tt.args.name, tt.args.email, tt.args.password, tt.args.role)

			if got.ID() != tt.want.id {
				t.Errorf("NewUser().ID() = %v, want %v", got.ID(), tt.want.id)
			}
			if got.Name() != tt.want.name {
				t.Errorf("NewUser().Name() = %v, want %v", got.Name(), tt.want.name)
			}
			if got.Email() != tt.want.email {
				t.Errorf("NewUser().Email() = %v, want %v", got.Email(), tt.want.email)
			}
			if got.Password() != tt.want.password {
				t.Errorf("NewUser().Password() = %v, want %v", got.Password(), tt.want.password)
			}
			if got.Role() != tt.want.role {
				t.Errorf("NewUser().Role() = %v, want %v", got.Role(), tt.want.role)
			}
		})
	}
}
