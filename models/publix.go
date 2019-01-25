package models

type (
	StoreResponse struct {
		WASTORENUM string `json:"WASTORENUM"`
		NAME       string `json:"NAME"`
		ADDR       string `json:"ADDR"`
		CITY       string `json:"CITY"`
		STATE      string `json:"STATE"`
		CLAT       string `json:"CLAT"`
		CLON       string `json:"CLON"`
	}
)

type (
	Store struct {
		ID       int    `json:"-"`
		StoreNum string `json:"storeNum"`
		Name     string `json:"name"`
		Address  string `json:"address"`
		City     string `json:"city"`
		State    string `json:"state"`
		Lat      string `json:"lat"`
		Long     string `json:"long"`
	}
)
