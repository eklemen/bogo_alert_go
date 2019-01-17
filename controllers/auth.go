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
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Println("-=-=-=-=-=-=-=-", string(bytes))
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println("ERROR:::::", err)
	return err == nil
}

func Register(c echo.Context) error {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		ZipCode  string `json:"zipCode"`
		Phone    string `json:"phone"`
	}{}
	u := models.NewUser()

	if err := c.Bind(&req); err != nil {
		return err
	}
	hash, err := HashPassword(req.Password)

	if err != nil {
		return err
	}

	f := app.DB.Where(
		&models.User{
			Email: req.Email,
		}).First(&u)

	if f.RecordNotFound() {
		uid, _ := uuid.NewV4()
		u.Uuid = uid
		u.Password = hash
		u.Email = req.Email
		u.ZipCode = req.ZipCode
		u.Phone = req.Phone

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
	claims["user"] = u
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	return &UserWithToken{User: &u, Token: t}, nil
}

func Login(c echo.Context) error {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	u := &models.User{Email: req.Email}
	us := models.NewUser()
	u.Email = req.Email
	f := app.DB.
		Preload("Terms").
		Where(u).
		First(&us)

	if f.RecordNotFound() {
		return c.JSON(http.StatusNotFound, "User is not registered.")
	}
	if CheckPasswordHash(req.Password, us.Password) {
		res, err := AssignToken(*us)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, res)
	}
	return c.JSON(http.StatusUnauthorized, "Incorrect email or password.")
}
