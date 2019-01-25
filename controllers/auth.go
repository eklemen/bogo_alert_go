package controllers

import (
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
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(c echo.Context) error {
	req := struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		ZipCode    string `json:"zipCode"`
		Phone      string `json:"phone"`
		EmailOptIn bool   `json:"emailOptIn"`
		TextOptIn  bool   `json:"textOptIn"`
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
		u.EmailOptIn = req.EmailOptIn
		u.TextOptIn = req.TextOptIn
		u.Phone = req.Phone

		res, err := AssignToken(*u)
		if err != nil {
			return err
		}
		u.Token = res.Token

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

func AssignToken(u models.User) (*models.User, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims (the DB id is encoded below)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	u.Token = t
	return &u, nil
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
		app.DB.Save(&res)
		return c.JSON(http.StatusOK, res)
	}
	return c.JSON(http.StatusUnauthorized, "Incorrect email or password.")
}

func Logout(c echo.Context) error {
	uid := c.Get("userId").(int)
	u := &models.User{ID: uid}
	app.DB.Model(&u).Update("token", "")
	return c.JSON(http.StatusNoContent, "")
}
