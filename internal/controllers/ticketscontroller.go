package controllers

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mzhn-sochi/gateway/api/ts"
	"github.com/mzhn-sochi/gateway/internal/entity"
	"github.com/mzhn-sochi/gateway/pkg/middleware"
	"log/slog"
)

type TicketsService interface {
	Find(ctx context.Context, id string) (*ts.Ticket, error)
	List(ctx context.Context, filters *entity.TicketFilters) ([]*ts.Ticket, uint64, error)
	//Create(t *ts.Ticket) error
	//Update(t *ts.Ticket) error
	//Delete(id string) error
}

type TicketController struct {
	service   TicketsService
	validator *validator.Validate
}

func NewTicketController(service TicketsService) *TicketController {
	return &TicketController{service: service, validator: validator.New()}
}

func (c *TicketController) Find() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")

		logger := ctx.Locals(middleware.LOGGER).(*slog.Logger).With("service", "tickets").With("method", "Find")

		ctx.Locals(middleware.LOGGER, logger)

		t, err := c.service.Find(ctx.Context(), id)
		if err != nil {
			return internal(err.Error())
		}

		return ok(ctx, t)
	}
}

func (c *TicketController) List() fiber.Handler {

	type query struct {
		Limit  uint64 `query:"limit"`
		Offset uint64 `query:"offset"`
		UserId string `query:"userId"`
	}

	type response struct {
		Tickets []*ts.Ticket `json:"tickets"`
		Total   uint64       `json:"count"`
	}

	return func(ctx *fiber.Ctx) error {
		var q query
		if err := ctx.QueryParser(&q); err != nil {
			return err
		}

		logger := ctx.Locals(middleware.LOGGER).(*slog.Logger).With("service", "tickets").With("method", "List")
		ctx.Locals(middleware.LOGGER, logger)

		logger.Debug("list tickets", slog.Uint64("limit", q.Limit), slog.Uint64("offset", q.Offset), slog.String("userId", q.UserId))

		filters := &entity.TicketFilters{
			Filters: entity.Filters{
				Limit:  q.Limit,
				Offset: q.Offset,
			},
		}

		if q.UserId != "" {
			filters.UserId = &q.UserId
		}

		tickets, total, err := c.service.List(ctx.Context(), filters)
		if err != nil {
			return internal(err.Error())
		}

		return ok(ctx, &response{
			Tickets: tickets,
			Total:   total,
		})
	}
}
