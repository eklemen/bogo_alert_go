package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/eklemen/bogo_alert/app"
	"github.com/eklemen/bogo_alert/models"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
)

func GetStoresByZip(c echo.Context) error {
	uid := c.Get("userId").(int)
	u := &models.User{ID: uid}

	req := struct {
		ZipCode string `json:"zipCode"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	app.DB.Model(&u).Update("zipCode", req.ZipCode)

	purl := "https://services.publix.com/api/v1/storelocation?types=R,G,H,N,S&includeOpenAndCloseDates=true&useLegacySchedule=false&count=10&zipCode="
	pzip := fmt.Sprintf("%s%v", purl, req.ZipCode)

	r, _ := http.NewRequest("GET", pzip, nil)

	res, _ := http.DefaultClient.Do(r)

	resData, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	resSt := struct {
		Stores []models.StoreResponse `json:"STORES"`
	}{}

	json.Unmarshal(resData, &resSt)

	return c.JSON(http.StatusOK, resSt)
}
