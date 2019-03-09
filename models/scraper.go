package models

import "time"

type (
	Scrapedstore struct {
		ID      int       `json:"id"`
		StoreID string    `json:"storeId"`
		Date    time.Time `json:"date"`
		Deals   []Deal    `json:"deals" gorm:"many2many:scrapedstore_deals"`
	}
)

type (
	Deal struct {
		ID   int    `json:"id"`
		Item string `json:"item"`
	}
)
