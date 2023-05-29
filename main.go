package main

import (
	"github.com/Arueljust/controllers"
	"github.com/Arueljust/initializers"
	"github.com/Arueljust/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.Connection()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/v1/signup", controllers.SignUp)
	r.POST("/v1/login", controllers.Login)
	r.GET("/v1/validate", middleware.Auth, controllers.Validate)
	r.Run()
}
