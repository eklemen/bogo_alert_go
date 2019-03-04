package controllers

import (
	"github.com/eklemen/bogo_alert/app"
	"github.com/eklemen/bogo_alert/models"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

func UpdateSearchTerms(c echo.Context) error {
	uid := c.Get("userId").(int)
	u := &models.User{ID: uid}
	app.DB.Preload("Terms").
		Preload("Store").
		Where(u).
		First(&u)

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

	// get all terms in db that match the given array
	// build out 2 arrays (will be 2 nested loops)
	// one array hold what will be added
	// one array will be what to remove

	app.DB.Model(&u).Association("Terms").Replace(terms)

	return c.JSON(http.StatusOK, u)
}

func GetUser(c echo.Context) error {
	uid := c.Get("userId").(int)
	u := &models.User{ID: uid}
	app.DB.Preload("Terms").
		Preload("Store").
		Where(&u).First(u)

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
	app.DB.Preload("Terms").
		Preload("Store").
		First(&res)

	return c.JSON(http.StatusOK, res)
}

func UpdateUserStore(c echo.Context) error {
	uid := c.Get("userId").(int)

	s := &models.Store{}
	u := &models.User{ID: uid}

	if err := c.Bind(&s); err != nil {
		return err
	}

	i := app.DB.Where(models.Store{StoreNum: s.StoreNum}).
		FirstOrCreate(&s)
	if i.Error != nil {
		return i.Error
	}
	app.DB.Model(&u).Update("StoreID", s.ID)

	app.DB.Preload("Terms").
		Preload("Store").
		First(&u)

	return c.JSON(http.StatusOK, u)
}

func GetStoreIds(c echo.Context) error {
	s := models.Term{}
	app.DB.Find(&s)

	return c.JSON(http.StatusOK, s)
}
