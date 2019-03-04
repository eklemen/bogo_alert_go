package models

import "time"

type (
	Scrapedstore struct {
		ID      int       `json:"-"`
		StoreID string    `json:"storeId"`
		Date    time.Time `json:"date"`
		Deals   []*Deal   `json:"-" gorm:"many2many:scrapedstore_deals"`
	}
)

type (
	Deal struct {
		ID   int    `json:"-"`
		Item string `json:"item"`
	}
)
