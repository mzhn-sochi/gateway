package app

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mzhn-sochi/gateway/api/auth"
	"github.com/mzhn-sochi/gateway/internal/config"
	"github.com/mzhn-sochi/gateway/internal/controllers"
	"github.com/mzhn-sochi/gateway/pkg/middleware"
	"log/slog"
)

type App struct {
	app    *fiber.App
	config *config.Config
	logger *slog.Logger

	suggestionsController *controllers.SuggestionsController
	AuthController        *controllers.AuthController
	ticketController      *controllers.TicketController
}

func newApp(
	config *config.Config,
	log *slog.Logger,
	suggestionsController *controllers.SuggestionsController,
	authController *controllers.AuthController,
	ticketController *controllers.TicketController,
) *App {
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
		app:                   app,
		config:                config,
		logger:                log,
		suggestionsController: suggestionsController,
		AuthController:        authController,
		ticketController:      ticketController,
	}
}

func (a *App) Run() error {

	host := a.config.App.Host
	port := a.config.App.Port

	a.app.Use(logger.New())
	a.app.Use(middleware.AttachRequestId())
	a.app.Use(middleware.AttachLogger(a.logger))

	a.app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://localhost:5173, http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	v1 := a.app.Group("/api/v1")

	au := v1.Group("/auth")
	au.Post("/sign-in", a.AuthController.SignIn())
	au.Post("/sign-up", a.AuthController.SignUp())
	au.Post("/sign-out", a.AuthController.AuthRequired(auth.Role_user), a.AuthController.SignOut())
	au.Post("/refresh", a.AuthController.Refresh())

	v1.Get("/suggestions", a.suggestionsController.GetSuggestions())

	tt := v1.Group("/tickets")
	tt.Get("/", a.ticketController.List())
	tt.Get("/:id", a.ticketController.Find())
	tt.Post("/", a.AuthController.AuthRequired(auth.Role_user), a.ticketController.Create())
	tt.Patch("/:id", a.ticketController.CloseTicket())

	v1.Get("/user/tickets", a.AuthController.AuthRequired(auth.Role_user), a.ticketController.ListUsers())
	v1.Get("/profile", a.AuthController.AuthRequired(auth.Role_user), a.AuthController.Profile())

	a.logger.Info("server started", slog.String("host", host), slog.Int("port", port))
	return a.app.Listen(fmt.Sprintf("%s:%d", host, port))
}

func (a *App) Shutdown() {
	a.app.Shutdown()
}
