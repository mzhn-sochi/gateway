package middleware

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

const LOGGER = "logger"

func AttachLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid := c.Locals(REQUEST_ID).(string)
		l := logger.With(slog.String("requestId", rid))
		c.Locals(LOGGER, l)
		return c.Next()
	}
}
