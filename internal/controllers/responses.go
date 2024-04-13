package controllers

import "github.com/gofiber/fiber/v2"

func internal(message string) error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}

func bad(message string) error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}

func unauthorized(message string) error {
	return fiber.NewError(fiber.StatusUnauthorized, message)
}

func forbidden(message string) error {
	return fiber.NewError(fiber.StatusForbidden, message)
}

func notFound(message string) error {
	return fiber.NewError(fiber.StatusNotFound, message)
}

func ok(ctx *fiber.Ctx, data ...interface{}) error {

	if len(data) == 0 {
		return ctx.SendStatus(fiber.StatusOK)
	}

	return ctx.JSON(&fiber.Map{
		"data": data[0],
	})
}
