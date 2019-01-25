package models

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"time"
)

type (
	User struct {
		ID         int        `json:"-"`
		Uuid       uuid.UUID  `json:"uuid"`
		Email      string     `json:"email"`
		Token      string     `json:"token"`
		Password   string     `json:"-"`
		Phone      string     `json:"phone"`
		ZipCode    string     `json:"zipCode"`
		EmailOptIn bool       `json:"emailOptIn"`
		TextOptIn  bool       `json:"textOptIn"`
		Terms      []*Term    `json:"-" gorm:"many2many:user_terms"`
		Store      *Store     `json:"store"`
		StoreID    int        `json:"-"`
		DeletedAt  *time.Time `json:"-"`
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

type Marshalable interface {
	MarshalJSON() ([]byte, error)
}

func (u *User) MarshalJSON() ([]byte, error) {

	termsBare := make([]string, len(u.Terms))
	for i, t := range u.Terms {
		termsBare[i] = t.Keyword
	}

	jsonStruct := struct {
		ID         int        `json:"-"`
		Uuid       uuid.UUID  `json:"uuid"`
		Email      string     `json:"email"`
		Token      string     `json:"token"`
		Password   string     `json:"-"`
		Phone      string     `json:"phone"`
		ZipCode    string     `json:"zipCode"`
		EmailOptIn bool       `json:"emailOptIn"`
		TextOptIn  bool       `json:"textOptIn"`
		Terms      []string   `json:"terms" gorm:"many2many:user_terms"`
		Store      *Store     `json:"store"`
		StoreID    int        `json:"-"`
		DeletedAt  *time.Time `json:"-"`
	}{
		ID:         u.ID,
		Terms:      termsBare,
		Uuid:       u.Uuid,
		Email:      u.Email,
		Token:      u.Token,
		Password:   u.Password,
		Phone:      u.Phone,
		ZipCode:    u.ZipCode,
		EmailOptIn: u.EmailOptIn,
		TextOptIn:  u.TextOptIn,
		Store:      u.Store,
		StoreID:    u.StoreID,
		DeletedAt:  u.DeletedAt,
	}

	return json.Marshal(&jsonStruct)
}
