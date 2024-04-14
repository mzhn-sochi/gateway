package ticketservice

import (
	"context"
	"fmt"
	"github.com/mzhn-sochi/gateway/api/ts"
	"github.com/mzhn-sochi/gateway/internal/config"
	"github.com/mzhn-sochi/gateway/internal/entity"
	"github.com/mzhn-sochi/gateway/internal/entity/dto"
	"github.com/mzhn-sochi/gateway/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"log/slog"
)

//var _ controllers.TicketsService = (*Service)(nil)

type Service struct {
	config *config.Config
	client ts.TicketServiceClient
}

func New(config *config.Config, logger *slog.Logger) *Service {

	l := logger.With("service", "ts")

	host := config.Services.TicketService.Host
	port := config.Services.TicketService.Port

	l.Info("connecting to grpc service", slog.String("host", host), slog.Int("port", port))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		l.Error("error with connection to grpc service", slog.String("err", err.Error()))
		panic(err)
	}

	client := ts.NewTicketServiceClient(conn)

	return &Service{
		config: config,
		client: client,
	}
}

func (s *Service) Find(ctx context.Context, id string) (*entity.Ticket, error) {
	req := &ts.FindByIdRequest{TicketId: id}

	ticket, err := s.client.FindById(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Printf("ticket from req %+v", ticket.Item)

	var item *entity.Item
	if ticket.Item != nil {
		item = &entity.Item{
			Product:     ticket.Item.Product,
			Description: ticket.Item.Description,
			Price:       ticket.Item.Price,
			Amount:      ticket.Item.Amount,
			Unit:        ticket.Item.Unit,
			Overprice:   ticket.Item.Overprice,
		}
	}

	log.Printf("item %+v", item)

	return &entity.Ticket{
		Id:          ticket.Id,
		UserId:      ticket.UserId,
		ImageUrl:    ticket.ImageUrl,
		Status:      entity.Role(ticket.Status),
		ShopName:    ticket.ShopName,
		ShopAddress: ticket.ShopAddress,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		Reason:      ticket.Reason,
		Item:        item,
	}, nil
}
func (s *Service) List(ctx context.Context, filters *entity.TicketFilters) ([]*entity.Ticket, uint64, error) {

	logger := ctx.Value(middleware.LOGGER).(*slog.Logger)

	req := &ts.ListRequest{
		Filter: &ts.Filter{
			UserId: filters.UserId,
		},
		Bounds: &ts.Bounds{
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	logger.Debug("ticket service list request", slog.Any("request", req))
	response, err := s.client.List(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	tt := make([]*entity.Ticket, 0, len(response.Tickets))
	for _, t := range response.Tickets {

		var item *entity.Item
		if t.Item != nil {
			item = &entity.Item{
				Product:     t.Item.Product,
				Description: t.Item.Description,
				Price:       t.Item.Price,
				Amount:      t.Item.Amount,
				Unit:        t.Item.Unit,
				Overprice:   t.Item.Overprice,
			}
		}

		tt = append(tt, &entity.Ticket{
			Id:          t.Id,
			UserId:      t.UserId,
			Status:      entity.Role(t.Status),
			ImageUrl:    t.ImageUrl,
			ShopName:    t.ShopName,
			ShopAddress: t.ShopAddress,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
			Reason:      t.Reason,
			Item:        item,
		})
	}

	return tt, uint64(response.Count), nil
}
func (s *Service) Create(ctx context.Context, t *dto.CreateTicket) (string, error) {
	l := ctx.Value(middleware.LOGGER).(*slog.Logger).With("service", "ts").With("method", "create")

	req := &ts.CreateRequest{
		UserId:   t.UserId,
		ImageUrl: t.ImageUrl,
		ShopAddr: t.ShopAddr,
		ShopName: t.ShopName,
	}

	l.Debug("ticket service create request", slog.Any("request", req))

	response, err := s.client.Create(ctx, req)
	if err != nil {
		return "", err
	}

	return response.TicketId, nil
}

func (s *Service) Close(ctx context.Context, id string) error {

	req := &ts.CloseTicketRequest{TicketId: id}

	if _, err := s.client.CloseTicket(ctx, req); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return ErrTicketNotFound
			}
		}

		return fmt.Errorf("error with closing ticket: %w", err)
	}

	return nil
}
