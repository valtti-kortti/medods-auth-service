package service

import (
	"time"

	"auth-service/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *service) CreateTokens(ctx *fiber.Ctx) error {
	// получаем guid из параметра запроса
	guidString := ctx.Params("user_guid")
	guid, err := uuid.Parse(guidString)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user GUID")
	}

	// генерируем токен рефреш и его хеш
	refreshToken, hashToken, err := s.refresh.GenerateToken()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// получаем IP клиента
	ip := ctx.IP()
	if ip == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "IP address is nil")
	}

	// получаем User-Agent клиента
	userAgent := ctx.Get("User-Agent")
	if userAgent == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "User-Agent is nil")
	}

	session := repository.Session{
		UserGUID:  guid,
		UserAgent: userAgent,
		TokenHash: hashToken,
		IpAddress: ip,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	// Создаем сессию в бд и получаем ее ID
	sessionID, err := s.repo.SetSession(ctx.Context(), &session)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// генерируем access токен
	accessToken, err := s.access.GenerateToken(sessionID, ip)
	if err != nil {
		_ = s.repo.DeleteSession(ctx.Context(), sessionID)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// устанавливаем куки с токенами
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
		SameSite: "Strict",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return ctx.SendStatus(fiber.StatusOK)
}
