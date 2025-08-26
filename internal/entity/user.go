package entity

import (
	"errors"
	"time"
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
       newUser := &User{
	       Name:      name,
	       Email:     email,
	       Password:  password,
	       CreatedAt: time.Now(),
	       Orders:    []Order{},
       }
       err := newUser.Validate()
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