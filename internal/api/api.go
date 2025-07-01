package api

import (
	"auth-service/internal/api/middleware"
	"auth-service/internal/service"
	"auth-service/internal/token"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Routers struct {
	Service service.Service
}

func NewRouters(r *Routers, tokenService token.JWTTokenManager) *fiber.App {
	app := fiber.New()

	// Настройка CORS (разрешенные методы, заголовки, авторизация)
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET, POST, PATCH, DELETE",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))

	appGroupToken := app.Group("/")

	appGroupToken.Get("/tokens/:user_guid", r.Service.CreateTokens)
	appGroupToken.Get("/refresh", r.Service.RefreshToken)

	appGroupUser := app.Group("/user", middleware.Authorization(tokenService))

	appGroupUser.Get("/guid", r.Service.GetUserGUID)
	appGroupUser.Get("/logout", r.Service.LogoutUser)

	return app
}
