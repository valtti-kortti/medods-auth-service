package service

import (
	"errors"
	"time"

	"auth-service/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func (s *service) RefreshToken(ctx *fiber.Ctx) error {
	// получае токены из кук
	accessToken := ctx.Cookies("access_token")
	refreshToken := ctx.Cookies("refresh_token")

	// проверка на пустоту токенов
	if accessToken == "" || refreshToken == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	// проверяем на валидность токен access и получаем его пейлоад
	claims, err := s.access.ValidToken(accessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// ищим сессию по id сессии
	session, err := s.repo.GetSession(ctx.Context(), claims.SessionID)
	if err != nil {
		if errors.Is(err, repository.ErrorNotFound) {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// проверяем на валидность токен refresh
	if err = s.refresh.CompareHash(refreshToken, session.TokenHash); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// проверяем изменения юзер агента, если изменился удаляем сессию
	if session.UserAgent != ctx.Get("User-Agent") {
		_ = s.repo.DeleteSession(ctx.Context(), claims.SessionID)
		return fiber.NewError(fiber.StatusUnauthorized, "User-Agent is changed")
	}

	// проверяем изменения IP, если изменился отправляем вебхук
	if session.IpAddress != ctx.IP() {
		go s.web.SendWebhook(ctx.IP(), session.IpAddress, session.UserGUID)
	}
	//генерируем новый токен refresh
	newRefreshToken, hashToken, err := s.refresh.GenerateToken()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ip := ctx.IP()
	if ip == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "IP address is nil")
	}

	// удаляем прошлую сессию
	err = s.repo.DeleteSession(ctx.Context(), session.SessionID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// создаем новую сессию
	sessionNew := repository.Session{
		UserGUID:  session.UserGUID,
		UserAgent: session.UserAgent,
		TokenHash: hashToken,
		IpAddress: ip,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	sessionID, err := s.repo.SetSession(ctx.Context(), &sessionNew)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// генерируем новый токен access
	accessTokenNew, err := s.access.GenerateToken(sessionID, ip)
	if err != nil {
		_ = s.repo.DeleteSession(ctx.Context(), sessionID)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessTokenNew,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
		SameSite: "Strict",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return ctx.SendStatus(fiber.StatusOK)
}
