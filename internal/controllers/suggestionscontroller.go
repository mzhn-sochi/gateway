package controllers

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mzhn-sochi/gateway/internal/entity"
)

type SuggestionsService interface {
	GetSuggestions(ctx context.Context, lon, lat float32, count uint32) ([]*entity.Suggestion, error)
}

type SuggestionsController struct {
	service   SuggestionsService
	validator *validator.Validate
}

func NewSuggestionsController(service SuggestionsService) *SuggestionsController {
	return &SuggestionsController{
		service:   service,
		validator: validator.New(),
	}
}

func (c *SuggestionsController) Register(router fiber.Router) {
	router.Get("/suggestions", c.getSuggestions())
}

func (c *SuggestionsController) getSuggestions() fiber.Handler {
	type suggestionsReq struct {
		Lon   float32 `query:"lon" validate:"required"`
		Lat   float32 `query:"lat" validate:"required"`
		Count uint32  `query:"count" validate:"required"`
	}

	return func(ctx *fiber.Ctx) error {
		var r suggestionsReq
		if err := ctx.QueryParser(&r); err != nil {
			return bad(err.Error())
		}

		if err := c.validator.Struct(r); err != nil {
			return bad(err.Error())
		}

		suggestions, err := c.service.GetSuggestions(ctx.Context(), r.Lon, r.Lat, r.Count)
		if err != nil {
			return internal(err.Error())
		}

		return ok(ctx, suggestions)
	}
}
