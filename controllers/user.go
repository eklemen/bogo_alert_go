package controllers

import (
	"fmt"
	"github.com/eklemen/bogo_alert/app"
	"github.com/eklemen/bogo_alert/models"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

func UpdateSearchTerms(c echo.Context) error {
	uid := c.Get("userId").(int)
	u := &models.User{ID: uid}
	app.DB.Where(u).First(&u)

	fmt.Println("------", u)
	req := struct {
		Terms []string `json:"terms"`
	}{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	newTerms := []models.Term{}

	for _, tstr := range req.Terms {
		for _, tt := range u.Terms {
			tstrK := strings.ToLower(strings.TrimSpace(tstr))
			ttK := strings.ToLower(strings.TrimSpace(tt.Keyword))
			if tstrK != ttK {
				newTerms = append(newTerms, models.Term{Keyword: tstrK})
			}
		}
	}

	terms := []models.Term{}
	for _, t := range req.Terms {
		terms = append(terms, models.Term{Keyword: t})
	}
	fmt.Println("UID", uid)

	app.DB.Model(&u).Association("Terms").Replace(terms)

	return c.JSON(http.StatusOK, u)
}

func GetUser(c echo.Context) error {
	uid := c.Get("userId").(int)
	u := &models.User{ID: uid}
	app.DB.Preload("Terms").Where(&u).First(u)

	return c.JSON(http.StatusOK, u)
}

func UpdateUser(c echo.Context) error {
	uid := c.Get("userId").(int)
	req := struct {
		Email   string `json:"email"`
		ZipCode string `json:"zipCode"`
		Phone   string `json:"phone"`
	}{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	u := &models.User{ID: uid}
	app.DB.Model(&u).Updates(models.User{
		ZipCode: req.ZipCode,
		Email:   req.Email,
		Phone:   req.Phone,
	})

	res := &models.User{ID: uid}
	app.DB.Preload("Terms").First(&res)

	return c.JSON(http.StatusOK, res)
}
