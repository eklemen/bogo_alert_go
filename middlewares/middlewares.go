package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var DB *gorm.DB

func LoadUserIntoContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		usr := c.Get("user").(*jwt.Token)
		claims := usr.Claims.(jwt.MapClaims)

		// set the user id into Context
		id := int(claims["id"].(float64))
		c.Set("userId", id)

		//// set the user uuid into Context
		//uid, _ := uuid.FromString(claims["uuid"].(string))
		//c.Set("uuid", uid)

		//user := models.User{}
		//err := DB.First(&user, id).Error
		//if err != nil {
		//	c.JSON(http.StatusNotFound, "Not found")
		//}
		//c.Set("user", user)
		return next(c)
	}
}
