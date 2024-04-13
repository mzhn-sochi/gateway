package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mzhn-sochi/gateway/internal/entity"
	"github.com/mzhn-sochi/gateway/pkg/file"
)

type AnalyzerService interface {
	Analyze(ctx context.Context, r file.Reader) (*entity.ImageInfo, error)
}

type AnalyzerController struct {
	service AnalyzerService
}

func NewAnalyzerController(service AnalyzerService) *AnalyzerController {
	return &AnalyzerController{
		service: service,
	}
}

func (c *AnalyzerController) Analyze() fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		f, err := ctx.FormFile("pricetag")
		if err != nil {
			return bad(err.Error())
		}

		reader, err := f.Open()
		if err != nil {
			return internal(err.Error())
		}

		r := file.NewReader(reader, f.Size)

		info, err := c.service.Analyze(ctx.Context(), r)
		if err != nil {
			return internal(err.Error())
		}

		return ok(ctx, info)

	}
}
