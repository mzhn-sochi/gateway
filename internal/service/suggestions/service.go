package suggestions

import (
	"context"
	"fmt"
	"github.com/mzhn-sochi/gateway/api/shop_suggestions"
	"github.com/mzhn-sochi/gateway/internal/config"
	"github.com/mzhn-sochi/gateway/internal/controllers"
	"github.com/mzhn-sochi/gateway/internal/entity"
	"github.com/mzhn-sochi/gateway/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type Service struct {
	config *config.Config
	logger *slog.Logger
	client shop_suggestions.ShopSuggestionsServiceClient
}

func New(config *config.Config, logger *slog.Logger) controllers.SuggestionsService {
	l := logger.With("service", "suggestions")

	host := config.Services.Suggestions.Host
	port := config.Services.Suggestions.Port

	l.Info("connecting to grpc service", slog.String("host", host), slog.Int("port", port))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		l.Error("error with connection to grpc service", slog.String("err", err.Error()))
		panic(err)
	}

	client := shop_suggestions.NewShopSuggestionsServiceClient(conn)

	return &Service{
		config: config,
		client: client,
		logger: l,
	}
}

func (s *Service) GetSuggestions(ctx context.Context, lon, lat float32, count uint32) ([]*entity.Suggestion, error) {
	rid := ctx.Value(middleware.REQUEST_ID)
	l := s.logger.With("request_id", rid)

	l.Debug("start get suggestions")
	res, err := s.client.GetSuggestion(ctx, &shop_suggestions.SuggestionOptions{
		Lon:   lon,
		Count: count,
		Lat:   lat,
	})
	if err != nil {
		return nil, err
	}

	suggestions := make([]*entity.Suggestion, len(res.Suggestions))
	for i, item := range res.Suggestions {
		suggestions[i] = &entity.Suggestion{
			Distance: item.Distance,
			Subtitle: item.Subtitle,
			Title:    item.Title,
		}
	}

	return suggestions, nil
}
