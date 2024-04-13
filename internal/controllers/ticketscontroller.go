package controllers

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mzhn-sochi/gateway/internal/entity"
	"github.com/mzhn-sochi/gateway/pkg/file"
	"github.com/mzhn-sochi/gateway/pkg/middleware"
)

type TicketsService interface {
	Find(ctx context.Context, id string) (*entity.Ticket, error)
	List(ctx context.Context, filters *entity.TicketFilters) ([]*entity.Ticket, uint64, error)
	Create(ctx context.Context, userId string, url string, addr string) (string, error)
	//Update(t *ts.Ticket) error
	//Delete(id string) error
}

type FileUploader interface {
	Upload(ctx context.Context, reader file.Reader) (string, error)
}

type TicketController struct {
	service      TicketsService
	fileUploader FileUploader
	validator    *validator.Validate
}

func NewTicketController(service TicketsService, fileUploader FileUploader) *TicketController {
	return &TicketController{
		service:      service,
		validator:    validator.New(),
		fileUploader: fileUploader,
	}
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
		Tickets []*entity.Ticket `json:"tickets"`
		Total   uint64           `json:"total"`
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

func (c *TicketController) ListUsers() fiber.Handler {

	type query struct {
		Limit  uint64 `query:"limit"`
		Offset uint64 `query:"offset"`
	}

	return func(ctx *fiber.Ctx) error {
		logger := ctx.Locals(middleware.LOGGER).(*slog.Logger).With("service", "tickets").With("method", "ListUsers")
		ctx.Locals(middleware.LOGGER, logger)

		var q query
		if err := ctx.QueryParser(&q); err != nil {
			return err
		}

		logger.Debug("list tickets of users", slog.Uint64("limit", q.Limit), slog.Uint64("offset", q.Offset))

		u := ctx.Locals("user").(*entity.UserClaims)

		filters := &entity.TicketFilters{
			Filters: entity.Filters{
				Limit:  q.Limit,
				Offset: q.Offset,
			},
			UserId: &u.Id,
		}

		logger.Debug("list users")

		tickets, _, err := c.service.List(ctx.Context(), filters)
		if err != nil {
			return internal(err.Error())
		}

		return ok(ctx, tickets)
	}
}

func (c *TicketController) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		addr := ctx.FormValue("address", "")
		if addr == "" {
			return bad("address is required")
		}

		f, err := ctx.FormFile("pricetag")
		if err != nil {
			return bad(err.Error())
		}

		logger := ctx.Locals(middleware.LOGGER).(*slog.Logger).With("service", "tickets").With("method", "Create")
		ctx.Locals(middleware.LOGGER, logger)

		reader, err := f.Open()
		if err != nil {
			return internal(err.Error())
		}

		ctype := f.Header.Get("Content-Type")
		logger.Debug("upload file", slog.String("file", f.Filename), slog.String("content-type", ctype))

		r := file.NewReader(reader, f.Size, ctype)
		logger.Debug("upload file", slog.String("file", f.Filename))

		url, err := c.fileUploader.Upload(ctx.Context(), r)
		if err != nil {
			return internal(err.Error())
		}

		uid := "test-user"
		if u, ok := ctx.Locals("user").(*entity.UserClaims); ok {
			uid = u.Id
		}

		ticketId, err := c.service.Create(ctx.Context(), uid, url, addr)
		if err != nil {
			return internal(err.Error())
		}

		return ok(ctx, ticketId)
	}
}
