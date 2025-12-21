package routes

import (
	jwtware "github.com/gofiber/contrib/jwt" // Gunakan contrib/jwt untuk v5 support
	"github.com/gofiber/fiber/v2"
	"github.com/mohod24/go-project-management/config"
	"github.com/mohod24/go-project-management/controllers"
	"github.com/mohod24/go-project-management/utils"
)

func Setup(app *fiber.App, uc *controllers.UserController) {
	// Public Routes
	auth := app.Group("/v1/auth")
	auth.Post("/register", uc.Register)
	auth.Post("/login", uc.Login)

	// JWT Protected Routes
	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.AppConfig.JWTSecret)}, // PERBAIKAN: Bukan string literal
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c, "Error unauthorized", err.Error())
		},
	}))

	// User Routes
	userGroup := api.Group("/users")
	userGroup.Get("/page", uc.GetUserPagination)
	userGroup.Get("/:id", uc.GetUser)
}