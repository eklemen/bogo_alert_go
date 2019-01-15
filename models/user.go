package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type (
	User struct {
		ID        int        `json:"-"`
		Uuid      uuid.UUID  `json:"uuid"`
		Email     string     `json:"email"`
		Token     string     `json:"token"`
		Password  string     `json:"-"`
		Phone     string     `json:"phone"`
		ZipCode   string     `json:"zipCode"`
		Terms     []*Term    `json:"terms" gorm:"many2many:user_terms"`
		DeletedAt *time.Time `json:"-"`
	}
)

type (
	Credentails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

func NewUser() *User {
	return &User{
		Terms: []*Term{},
	}
}
