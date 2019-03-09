package controllers

import (
	"fmt"
	"github.com/eklemen/bogo_alert/app"
	"github.com/eklemen/bogo_alert/models"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func SetStoresDeals(c echo.Context) error {
	ss := models.Scrapedstore{}
	req := struct {
		StoreID string `json:"storeId"`
		Date    string `json:"date"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", req.Date)
	fmt.Println("TIME----", t)
	if err != nil {
		return err
	}

	ss.Date = t
	ss.StoreID = req.StoreID

	app.DB.FirstOrCreate(&ss)

	return c.JSON(http.StatusOK, ss)
}
