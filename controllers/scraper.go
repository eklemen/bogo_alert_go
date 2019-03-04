package controllers

import (
	"github.com/labstack/echo"
	"net/http"
)

func SetStoresDeals(c echo.Context) error {

	return c.JSON(http.StatusOK, "")
}
