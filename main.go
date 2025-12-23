package main

// @title Go Project Management API
// @version 1.0
// @description This is a sample server for a project management application.
// @termsOfService http://swagger.io/terms/
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

// @contact.name API Support
func main() {
	// Load environment variables and connect to the database
	config.LoadEnv()
	config.ConnectDB()

	// Seed initial admin user
	seed.SeedAdmin()
	app := fiber.New()

	// Initialize repositories, services, and controllers
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Initialize Board components
	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	// Setup routes
	routes.Setup(app, userController, boardController)
	port := config.AppConfig.AppPort
	log.Println("Server running on port " + port)
	app.Listen(":" + port)
	log.Fatal(app.Listen(":" + port))

}