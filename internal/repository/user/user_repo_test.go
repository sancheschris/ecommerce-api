package repository

import (
	"testing"

	"github.com/sancheschris/ecommerce-api/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&model.User{})
	user, err := model.NewUser("Chris", "chris@gm.com", "secret")
	assert.NoError(t, err)
	assert.NotNil(t, user.Name)

	userDB := NewUser(db)
	err = userDB.Create(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Name)
}

func TestFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&model.User{})
	user, err := model.NewUser("Chris", "chris@gm.com", "secret")
	assert.NoError(t, err)
	assert.NotNil(t, user)

	userDB := NewUser(db)
	err = userDB.Create(user)
	assert.NoError(t, err)
	assert.Nil(t, err)

	usrByEmail, err := userDB.FindByEmail("chris@gm.com")
	assert.NoError(t, err)
	assert.NotNil(t, usrByEmail)
	assert.Equal(t, "Chris", usrByEmail.Name)
	assert.Equal(t, "chris@gm.com", usrByEmail.Email)

	assert.False(t, usrByEmail.Email == "Bob")
}

func TestCreateUserTable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&model.User{})
	userDB := NewUser(db)

	tests := []struct {
		name string
		email string
		password string
		want bool
	} {
		{"Chris", "chris@gm.com", "secret", false},
		{"", "bob@gm.com", "123", true},
		{"Tedd", "", "secret", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := model.NewUser(tt.name, tt.email, tt.password)
			if tt.want {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				err = userDB.Create(user)
				assert.NoError(t, err)
			}
		})
	}
}