package controllers

import (
	"fmt"
	"github.com/eklemen/bogo_alert/app"
	"github.com/eklemen/bogo_alert/models"
	"github.com/labstack/echo"
	"net/http"
)

func UpdateSearchTerms(c echo.Context) error {
	uid := c.Get("userId").(int)
	u := &models.User{ID: uid}

	req := struct {
		Terms []string `json:"terms"`
	}{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	terms := []models.Term{}
	for _, t := range req.Terms {
		terms = append(terms, models.Term{Keyword: t})
		fmt.Println("terms", t)
	}
	fmt.Println("UID", uid)

	app.DB.Model(&u).Association("Terms").Append(terms)
	return c.JSON(http.StatusOK, u)
}
