package main

import (
	"github.com/gorilla/sessions"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//"github.com/markbates/goth"
	//"github.com/markbates/goth/gothic"
	"fmt"
	"github.com/eklemen/bogo_alert/config"
	"github.com/eklemen/bogo_alert/models"
	"github.com/subosito/gotenv"
	"net/http"
	"os"
)

//var db *gorm.DB

func helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "hi")
}

func main() {
	gotenv.Load()
	// Dont forget to add postgres adapter to imports
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	if err := config.Open(); err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("DB Connected...")
	}
	defer config.DB.Close()

	// TODO: create a struct for these
	config.DB.LogMode(true)
	// Migrate the schema
	config.DB.AutoMigrate(&models.User{})

	e := echo.New()
	e.Debug = true
	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowHeaders:     []string{echo.HeaderAccessControlAllowOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Authentication strategies
	key := os.Getenv("GOTH_SESSION_SECRET")
	maxAge := 86400 * 90 // 90 days

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = false

	// Routes //
	e.GET("/hi", helloHandler)

	// Start server
	e.Logger.Fatal(e.Start(os.Getenv("SERVER_PORT")))
}
