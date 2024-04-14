package ticketservice

import (
	"context"
	"errors"
	"github.com/mzhn-sochi/gateway/api/ts"
	"io"
	"log/slog"
)

func (s *Service) ShopSummary(ctx context.Context) (map[string]int64, error) {

	l := ctx.Value("logger").(*slog.Logger).With("service", "service").With("method", "ShopSummary")

	stream, err := s.client.GetShopSummary(ctx, &ts.Empty{})
	if err != nil {
		return nil, err
	}

	r := make(map[string]int64)

	for {
		recv, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		l.Debug("receive a record", "shopId", recv.ShopId, "count", recv.Count)
		r[recv.ShopId] = recv.Count
	}

	l.Debug("received all records")

	return r, nil
}
func (s *Service) UserSummary(ctx context.Context) (map[string]int64, error) {

	l := ctx.Value("logger").(*slog.Logger).With("service", "service").With("method", "UserSummary")

	stream, err := s.client.GetUserSummary(ctx, &ts.Empty{})
	if err != nil {
		return nil, err
	}

	r := make(map[string]int64)

	for {
		recv, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		l.Debug("receive a record", "shopId", recv.UserId, "count", recv.Count)
		r[recv.UserId] = recv.Count
	}
	l.Debug("received all records")

	return r, nil
}

func (s *Service) StatusSummary(ctx context.Context) (map[string]int64, error) {

	l := ctx.Value("logger").(*slog.Logger).With("service", "service").With("method", "StatusSummary")

	stream, err := s.client.GetStatusSummary(ctx, &ts.Empty{})
	if err != nil {
		return nil, err
	}

	r := make(map[string]int64)

	for {
		recv, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		l.Debug("receive a record", "shopId", recv.StatusId, "count", recv.Count)
		r[recv.StatusId] = recv.Count
	}

	l.Debug("received all records")

	return r, nil
}
