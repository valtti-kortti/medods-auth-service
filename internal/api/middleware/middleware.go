package middleware

import (
	"auth-service/internal/token"
	"github.com/gofiber/fiber/v2"
	"log"
)

func Authorization(tokenService token.JWTTokenManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			log.Println("Access token not found in cookies")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization token required",
			})
		}

		claims, err := tokenService.ValidToken(accessToken)
		if err != nil {
			log.Println("Middleware: token invalid:", err)
			return c.SendStatus(fiber.StatusUnauthorized)
		} else {
			log.Printf("Middleware: token OK: session_id=%d user_ip=%s", claims.SessionID, claims.UserIP)
		}

		return c.Next()
	}
}
