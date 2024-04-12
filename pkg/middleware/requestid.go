package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const REQUEST_ID = "requestId"

func AttachRequestId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestId := uuid.NewString()
		c.Locals(REQUEST_ID, requestId)
		c.Response().Header.Set("X-Request-Id", requestId)
		return c.Next()
	}
}
