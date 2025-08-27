package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNameIsRequired = errors.New("name is required")
	ErrEmailIsRequired = errors.New("email is required")
)

type User struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"-"`
	Orders []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(name, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}   
	newUser := &User{
	       Name:      name,
	       Email:     email,
	       Password:  string(hash),
	       CreatedAt: time.Now(),
	       Orders:    []Order{},
       }
	err = newUser.Validate()
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameIsRequired
	}
	if u.Email == "" {
		return ErrEmailIsRequired
	}
	return nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}