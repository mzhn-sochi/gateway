package controllers

import "github.com/gofiber/fiber/v2"

func internal(message string) error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}

func bad(message string) error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}

func ok(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(&fiber.Map{
		"data": data,
	})
}
