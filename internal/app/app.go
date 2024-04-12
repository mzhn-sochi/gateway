package app

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mzhn-sochi/gateway/internal/config"
	"github.com/mzhn-sochi/gateway/internal/controllers"
	"github.com/mzhn-sochi/gateway/pkg/middleware"
)

type App struct {
	app    *fiber.App
	config *config.Config
	logger *slog.Logger

	controllers []controllers.Controller
}

func newApp(config *config.Config, log *slog.Logger,
	analyzerController *controllers.AnalyzerController,
	suggestionsController *controllers.SuggestionsController) *App {

	app := fiber.New(fiber.Config{
		AppName:       "sochya-gateway",
		CaseSensitive: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			err = ctx.Status(code).JSON(fiber.Map{
				"message": e.Message,
			})

			return nil
		},
		BodyLimit: 10 << 20,
	})

	return &App{
		app:    app,
		config: config,
		logger: log,
		controllers: []controllers.Controller{
			analyzerController,
			suggestionsController,
		},
	}
}

func (a *App) Run() error {

	host := a.config.App.Host
	port := a.config.App.Port

	a.app.Use(logger.New())
	a.app.Use(middleware.AttachRequestId())

	a.app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://localhost:5173, http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	v1 := a.app.Group("/api/v1")

	for _, c := range a.controllers {
		a.logger.Debug(fmt.Sprintf("register %s", reflect.TypeOf(c).Elem().Name()))
		c.Register(v1)
	}

	a.logger.Info("server started", slog.String("host", host), slog.Int("port", port))
	return a.app.Listen(fmt.Sprintf("%s:%d", host, port))
}

func (a *App) Shutdown() {
	a.app.Shutdown()
}
