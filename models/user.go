package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type (
	User struct {
		ID              int          `json:"-"`
		Uuid            uuid.UUID    `json:"uuid"`
		Email           string       `json:"email"`
		Phone           string       `json:"phone"`
		ZipCode			string 		 `json:"zipCode"`
		SearchTerms []*Term `json:"terms" gorm:"many2many:user_terms"`
		DeletedAt       *time.Time   `json:"-"`
	}
)

func NewUser() *User {
	return &User{}
}

