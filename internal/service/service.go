package service

import (
	"auth-service/internal/repository"
	"auth-service/internal/token"
	"auth-service/internal/webhook"

	"github.com/gofiber/fiber/v2"
)

type Service interface {
	CreateTokens(ctx *fiber.Ctx) error
	RefreshToken(ctx *fiber.Ctx) error
	GetUserGUID(ctx *fiber.Ctx) error
	LogoutUser(ctx *fiber.Ctx) error
}

type service struct {
	repo    repository.Repository
	refresh token.RefreshManager
	access  token.JWTTokenManager
	web     webhook.Webhook
}

func NewService(repo repository.Repository, ref token.RefreshManager, jwt token.JWTTokenManager, web webhook.Webhook) Service {
	return &service{repo: repo, refresh: ref, access: jwt, web: web}
}
