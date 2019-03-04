package main

import (
	"fmt"
	"github.com/eklemen/bogo_alert/app"
	"github.com/eklemen/bogo_alert/controllers"
	"github.com/eklemen/bogo_alert/middlewares"
	"github.com/eklemen/bogo_alert/models"
	"github.com/gorilla/sessions"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/subosito/gotenv"
	"os"
)

func main() {
	gotenv.Load()
	// Dont forget to add postgres adapter to imports
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	if err := app.Open(); err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("DB Connected...")
	}
	defer app.DB.Close()

	// TODO: create a struct for these
	app.DB.LogMode(true)
	// Migrate the schema
	app.DB.AutoMigrate(&models.User{})
	app.DB.AutoMigrate(&models.Term{})
	app.DB.AutoMigrate(&models.Store{})

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
	e.POST("/register", controllers.Register)
	e.POST("/login", controllers.Login)

	u := e.Group("/api")
	u.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))
	u.Use(middlewares.LoadUserIntoContext)
	u.POST("/logout", controllers.Logout)
	u.GET("/user", controllers.GetUser)
	u.PUT("/user", controllers.UpdateUser)
	u.POST("/user/store", controllers.UpdateUserStore)
	u.POST("/terms", controllers.UpdateSearchTerms)
	u.GET("/terms", controllers.GetStoreIds)
	u.GET("/terms/search/:term", controllers.TypeAhead)
	u.POST("/stores", controllers.GetStoresByZip)

	//p := e.Group("/s")
	//p.POST("/storesDeals", controllers.SetStoresDeals)

	// Start server
	e.Logger.Fatal(e.Start(os.Getenv("SERVER_PORT")))
}
