//go:build wireinject
// +build wireinject

package app

import (
	"github.com/mzhn-sochi/gateway/internal/config"
	"github.com/mzhn-sochi/gateway/internal/logger"

	"github.com/google/wire"
	"github.com/mzhn-sochi/gateway/internal/controllers"
	"github.com/mzhn-sochi/gateway/internal/service/analyzerservice"
)

func InitApp() *App {
	panic(wire.Build(
		newApp,
		wire.NewSet(logger.New),
		wire.NewSet(config.New),
		wire.NewSet(controllers.NewAnalyzerController),
		wire.Bind(new(controllers.AnalyzerService), new(*analyzerservice.Service)),

		wire.NewSet(analyzerservice.New),
	))
}
