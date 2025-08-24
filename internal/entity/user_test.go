package entity

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNewUser(t *testing.T) {
	u, err := NewUser("Edson",  "ed@gmail.com", "123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Name != "Edson" || u.Email != "ed@gmail.com" || u.Password != "123456" {
		t.Fatalf("fields not set correctly: %+v", u)
	}
}

func TestNewUserWithAssert(t * testing.T) {
	user, err := NewUser("Chris", "chris@gmail.com", "secret")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Name)
	assert.Equal(t, "Chris", user.Name)
	assert.Equal(t, "chris@gmail.com", user.Email)
}

func TestUserWhenNameIsRequired(t *testing.T) {
	u, err := NewUser("", "chris@gmail.com", "secret")
	assert.Nil(t, u)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestUserWhenEmailIsRequired(t *testing.T) {
	u, err := NewUser("Chris", "", "secret")
	assert.Nil(t, u)
	assert.Equal(t, ErrEmailIsRequired, err)
}

// Test using table-driven wat to test both Validate() and NewUser
func TestUserValidate(t *testing.T) {
	tests := []struct {
		name string
		user User
		wantErr error
	} {
		{"ok", User{Name: "Alice", Email: "a@b.com", Password: "x"}, nil},
		{"missing name", User{Name: "", Email: "a@b.com", Password: "x"}, ErrNameIsRequired},
		{"missing email", User{Name: "Alice", Email: "", Password: "x"}, ErrEmailIsRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := tt.user.Validate()
			if ans != tt.wantErr {
				t.Errorf("got %d, want %d", ans, tt.wantErr)
			}
			if !errors.Is(ans, tt.wantErr) {
				t.Fatalf("got error %v, want %v", ans, tt.wantErr)
			}
		})
	}
}