package service

import (
	"auth-service/internal/repository"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
)

// GetUserGUID godoc
// @Summary Получить GUID пользователя
// @Description Возвращает GUID пользователя из валидной сессии (требуется access token в cookies)
// @Tags User
// @Produce json
// @Success 200 {object} map[string]string "Пример: {"user_guid": "550e8400-e29b-41d4-a716-446655440000"}"
// @Failure 401 "Невалидный/отсутствующий access token"
// @Failure 404 "Сессия не найдена"
// @Router /user/guid [get]
// @Security ApiKeyAuth
func (s *service) GetUserGUID(ctx *fiber.Ctx) error {
	accessToken := ctx.Cookies("access_token")
	claims, err := s.access.ValidToken(accessToken)
	if err != nil {
		log.Println("Access token invalid")
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	session, err := s.repo.GetSession(ctx.Context(), claims.SessionID)
	if err != nil {
		if errors.Is(err, repository.ErrorNotFound) {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"user_guid": session.UserGUID})
}

// LogoutUser godoc
// @Summary Выход пользователя из системы
// @Description Удаляет сессию пользователя по access token (требуется валидный access token в cookies)
// @Tags User
// @Produce json
// @Success 200 "Сессия успешно удалена"
// @Failure 401 {object} map[string]string "Невалидный/отсутствующий access token"
// @Failure 404 {object} map[string]string "Сессия не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /user/logout [get]
// @Security ApiKeyAuth
func (s *service) LogoutUser(ctx *fiber.Ctx) error {
	accessToken := ctx.Cookies("access_token")
	claims, err := s.access.ValidToken(accessToken)
	if err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	err = s.repo.DeleteSession(ctx.Context(), claims.SessionID)
	if err != nil {
		if errors.Is(err, repository.ErrorNotFound) {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusOK)
}
