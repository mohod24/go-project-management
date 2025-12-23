package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3" // Gunakan contrib/jwt untuk v5 support
	"github.com/joho/godotenv"
	"github.com/mohod24/go-project-management/config"
	"github.com/mohod24/go-project-management/controllers"
	"github.com/mohod24/go-project-management/utils"
)

func Setup(app *fiber.App, 
	uc *controllers.UserController,
	bc *controllers.BoardController) {
	err := godotenv.Load()
		if err != nil{
		log.Fatal("Error loading .env file:", err)
	}
	
	// Public Routes
	auth := app.Group("/v1/auth")
	auth.Post("/register", uc.Register)
	auth.Post("/login", uc.Login)

	// JWT Protected Routes
	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.AppConfig.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c, "Error unauthorized", err.Error())
		},
	}))

	// User Routes
	userGroup := api.Group("/users")
	userGroup.Get("/page", uc.GetUserPagination)
	userGroup.Get("/:id", uc.GetUser)
	userGroup.Put("/:id", uc.UpdateUser)
	userGroup.Delete("/:id", uc.DeleteUser)

	// Board Routes
	boardGroup := api.Group("/boards")
	boardGroup.Post("/", bc.CreateBoard)
	boardGroup.Put("/:id", bc.UpdateBoard)
	boardGroup.Post("/:id/members", bc.AddBoardMember)
}