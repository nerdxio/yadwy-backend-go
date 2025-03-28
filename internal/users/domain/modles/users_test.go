package modles

import "testing"

func TestNewUser(t *testing.T) {
	// Arrange
	expected := User{
		name:     "valid user",
		email:    "john@example.com",
		password: "secret123",
		role:     RoleCustomer,
	}

	// Act
	actual, err := NewUser(expected.name, expected.email, expected.password, expected.role)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if *actual != expected {
		t.Errorf("got %+v, want %+v", *actual, expected)
	}
}
