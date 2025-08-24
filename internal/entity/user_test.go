package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNewUser(t *testing.T) {
	// arrange
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
	assert.NoError(t, err)
	assert.NotNil(t, user)
}