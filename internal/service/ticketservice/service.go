package ticketservice

import (
	"context"
	"fmt"
	"github.com/mzhn-sochi/gateway/api/share"
	"github.com/mzhn-sochi/gateway/api/ts"
	"github.com/mzhn-sochi/gateway/internal/config"
	"github.com/mzhn-sochi/gateway/internal/entity"
	"github.com/mzhn-sochi/gateway/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

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

func (s *Service) Find(ctx context.Context, id string) (*ts.Ticket, error) {
	req := &ts.FindByIdRequest{TicketId: id}

	ticket, err := s.client.FindById(ctx, req)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}
func (s *Service) List(ctx context.Context, filters *entity.TicketFilters) ([]*ts.Ticket, uint64, error) {

	logger := ctx.Value(middleware.LOGGER).(*slog.Logger)

	req := &ts.ListRequest{
		Filter: &ts.Filter{
			UserId: filters.UserId,
		},
		Bounds: &share.Bounds{
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	logger.Debug("ticket service list request", slog.Any("request", req))
	response, err := s.client.List(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return response.Tickets, uint64(response.Count), nil
}
