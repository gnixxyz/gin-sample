package main

import (
	"gin-sample/api"
	"gin-sample/models"
	"gin-sample/pkg/middleware"
	"gin-sample/views"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// load .env environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")

	// init models
	db, _ := models.Init()

	// create gin app
	app := gin.Default()

	// middleware
	app.Use(models.Inject(db))
	app.Use(middleware.JWTMiddleware())

	// routes
	api.Routes(app)
	views.Routes(app)

	// listen to given port
	app.Run(":" + port)
}
