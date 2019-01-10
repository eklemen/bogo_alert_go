package controllers

import (
	"github.com/labstack/echo"
	"github.com/eklemen/bogo_alert/models"
	"net/http"
	"github.com/eklemen/bogo_alert/config"
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
