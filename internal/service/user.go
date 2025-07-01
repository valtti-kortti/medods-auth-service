package service

import (
	"auth-service/internal/repository"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
)

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
