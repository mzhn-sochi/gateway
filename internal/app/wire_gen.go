// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/mzhn-sochi/gateway/internal/config"
	"github.com/mzhn-sochi/gateway/internal/controllers"
	"github.com/mzhn-sochi/gateway/internal/logger"
	"github.com/mzhn-sochi/gateway/internal/service/authservice"
	"github.com/mzhn-sochi/gateway/internal/service/s3service"
	"github.com/mzhn-sochi/gateway/internal/service/suggestions"
	"github.com/mzhn-sochi/gateway/internal/service/ticketservice"
)

// Injectors from wire.go:

func InitApp() *App {
	configConfig := config.New()
	slogLogger := logger.New(configConfig)
	suggestionsService := suggestions.New(configConfig, slogLogger)
	suggestionsController := controllers.NewSuggestionsController(suggestionsService)
	service := authservice.New(configConfig, slogLogger)
	authController := controllers.NewAuthController(service)
	ticketserviceService := ticketservice.New(configConfig, slogLogger)
	s3Service := s3service.New(configConfig, slogLogger)
	ticketController := controllers.NewTicketController(ticketserviceService, s3Service, service)
	app := newApp(configConfig, slogLogger, suggestionsController, authController, ticketController)
	return app
}
