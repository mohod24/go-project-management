package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mohod24/go-project-management/config"
	"github.com/mohod24/go-project-management/controllers"
	"github.com/mohod24/go-project-management/database/seed"
	"github.com/mohod24/go-project-management/repositories"
	"github.com/mohod24/go-project-management/routes"
	"github.com/mohod24/go-project-management/services"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()
	app := fiber.New()

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	routes.Setup(app, userController)
	port := config.AppConfig.AppPort
	log.Println("Server running on port " + port)
	app.Listen(":" + port)
	log.Fatal(app.Listen(":" + port))

}