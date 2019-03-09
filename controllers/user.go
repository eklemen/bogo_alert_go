package controllers

import (
	"github.com/eklemen/bogo_alert/app"
	"github.com/eklemen/bogo_alert/models"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

func inList(checkItem string, list []string) bool {
	clean := func(s string) string {
		return strings.ToLower(strings.TrimSpace(s))
	}
	found := false
	for _, i := range list {
		if clean(checkItem) == clean(i) {
			found = true
		}
	}
	return found
}

func termsToStrings(terms []*models.Term) []string {
	strs := []string{}
	for _, t := range terms {
		strs = append(strs, t.Keyword)
	}
	return strs
}

func UpdateSearchTerms(c echo.Context) error {
	uid := c.Get("userId").(int)
	u := &models.User{ID: uid}
	app.DB.Preload("Terms").
		First(&u)

	req := struct {
		Terms []string `json:"terms"`
	}{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	// Add all new terms from the request
	newTerms := []models.Term{}
	for _, reqTerm := range req.Terms {
		if !inList(reqTerm, termsToStrings(u.Terms)) {
			newTerms = append(newTerms, models.Term{Keyword: reqTerm})
		}
	}

	// Remove any terms not in the request
	removeTerms := []models.Term{}
	for _, currentTerm := range u.Terms {
		if !inList(currentTerm.Keyword, req.Terms) {
			removeTerms = append(removeTerms, models.Term{ID: currentTerm.ID})
		}
	}

	// Create all new associations
	for _, newTerm := range newTerms {
		var term models.Term
		// Already checked, just need to create
		t := app.DB.FirstOrCreate(&term, newTerm)
		if t.Error != nil {
			app.DB.Rollback()
			panic(t.Error)
			return t.Error
		}
		app.DB.Model(&u).Association("Terms").Append(term)
	}

	// Remove all old associations
	if len(removeTerms) != 0 {
		app.DB.Model(&u).Association("Terms").Delete(removeTerms)
	}

	return c.JSON(http.StatusOK, u.Terms)
}

func GetUser(c echo.Context) error {
	uid := c.Get("userId").(int)
	u := &models.User{ID: uid}
	app.DB.Preload("Terms").
		Preload("Store").
		Where(&u).First(u)

	return c.JSON(http.StatusOK, u)
}

func GetUserTerms(c echo.Context) error {
	uid := c.Get("userId").(int)
	u := &models.User{ID: uid, Terms: []*models.Term{}}
	app.DB.Preload("Terms").First(&u)

	return c.JSON(http.StatusOK, u.Terms)
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
