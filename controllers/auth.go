package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/eklemen/bogo_alert/app"
	"github.com/eklemen/bogo_alert/models"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(c echo.Context) error {
	u := models.NewUser()

	if err := c.Bind(u); err != nil {
		return err
	}
	hash, err := HashPassword(u.Password)

	if err != nil {
		return err
	}

	f := app.DB.Where(
		&models.User{
			Email: u.Email,
		}).First(&u)

	uid, _ := uuid.NewV4()
	if f.RecordNotFound() {
		u.Uuid = uid
		u.Password = hash

		app.DB.Create(&u)
		return c.JSON(http.StatusOK, u)
	}

	return c.JSON(http.StatusBadRequest, "User is already registered")
}

type (
	UserWithToken struct {
		User  *models.User `json:"user"`
		Token string       `json:"token"`
	}
)

func AssignToken(u models.User) (*UserWithToken, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims (the DB id is encoded below)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	//u.Token = t
	//app.DB.Update(&u)
	return &UserWithToken{User: &u, Token: t}, nil
}

func Login(c echo.Context) error {
	r := &models.Credentails{}
	if err := c.Bind(r); err != nil {
		return err
	}
	fmt.Println("u.Password", r.Password)

	u := &models.User{Email: r.Email}
	f := app.DB.
		Where(u).First(&u)

	if f.RecordNotFound() {
		return c.JSON(http.StatusNotFound, "User is not registered.")
	}
	fmt.Println("r.pass", r.Password)
	fmt.Println("u.pass", u.Password)
	if CheckPasswordHash(r.Password, u.Password) {
		fmt.Println("same-----")
		res, err := AssignToken(*u)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, res)
	}
	return c.JSON(http.StatusUnauthorized, "Incorrect email or password.")
}
