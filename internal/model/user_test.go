package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)


func TestNewUser(t *testing.T) {
	u, err := NewUser("Edson",  "ed@gmail.com", "123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Name != "Edson" || u.Email != "ed@gmail.com" {
		t.Fatalf("fields not set correctly: %+v", u)
	}
	// Verify password was hashed correctly
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("123456")); err != nil {
		t.Fatalf("password not hashed correctly: %v", err)
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

func TestNewUserTableDriven(t *testing.T) {
	tests := []struct {
		nameTest string
		name string
		email string
		password string
		want error
	} {
		{"ok", "Bob", "bob@gmail.com", "secret", nil},
		{"missing name", "", "bob@gmail.com", "secret", ErrNameIsRequired},
		{"missing email", "Bob", "", "secret", ErrEmailIsRequired},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {
			u, err := NewUser(tt.name, tt.email, tt.password)

			if tt.want == nil {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if u == nil {
					t.Errorf("got nil user on sucess")
				}
				if u != nil && (u.Name != tt.name || u.Email != tt.email) {
					t.Errorf("fields not set correctly: %+v", u)
				}
				// Verify password was hashed correctly
				if u != nil {
					if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(tt.password)); err != nil {
						t.Errorf("password not hashed correctly: %v", err)
					}
				}
				return
			}
			if u != nil {
				t.Error("expected nil user on error")
			}
			if err != tt.want {
				t.Errorf("got error %v, want %v", err, tt.want)
			}
		})
	}
}

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

func TestUserPassword(t *testing.T) {
	user, err := NewUser("Gabi", "g@gmail.com", "123456")
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("123456"))
	assert.False(t, user.ValidatePassword("1234567"))
	assert.NotEqual(t, "123456", user.Password)
}