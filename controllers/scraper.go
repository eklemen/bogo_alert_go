package controllers

import (
	"github.com/eklemen/bogo_alert/app"
	"github.com/eklemen/bogo_alert/models"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func SetStoresDeals(c echo.Context) error {
	ss := models.Scrapedstore{}
	req := struct {
		StoreID string   `json:"storeId"`
		Date    string   `json:"date"`
		Deals   []string `json:"deals"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	deals := []models.Deal{}

	for _, d := range req.Deals {
		deal := models.Deal{}
		app.DB.FirstOrCreate(&deal, models.Deal{Item: d})
		deals = append(deals, deal)
	}

	ss.Date = t
	ss.StoreID = req.StoreID
	app.DB.FirstOrCreate(&ss)
	app.DB.Model(&ss).Association("Deals").Append(deals)

	return c.JSON(http.StatusOK, ss)
}
