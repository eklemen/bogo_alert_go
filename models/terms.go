package models

type (
	Term struct {
		ID      int    `json:"-"`
		Keyword string `json:"keyword"`
	}
)
