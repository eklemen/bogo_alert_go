package controllers

import (
	"github.com/eklemen/bogo_alert/config"
	"github.com/eklemen/bogo_alert/models"
	"github.com/labstack/echo"
	"net/http"
)

func Register(c echo.Context) error {
	email := c.FormValue("email")
	u := models.NewUser()

	f := config.DB.Where(
		&models.User{
			Email: u.Email,
		}).First(&u)
	return c.JSON(http.StatusOK, f)
}
