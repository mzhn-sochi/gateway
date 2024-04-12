package analyzerservice

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/mzhn-sochi/gateway/api/pricetaganalyzer"
	"github.com/mzhn-sochi/gateway/internal/config"
	"github.com/mzhn-sochi/gateway/internal/controllers"
	"github.com/mzhn-sochi/gateway/internal/entity"
	"github.com/mzhn-sochi/gateway/pkg/file"
	"github.com/mzhn-sochi/gateway/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var _ controllers.AnalyzerService = (*Service)(nil)

type Service struct {
	config *config.Config
	logger *slog.Logger
	client pricetaganalyzer.PriceTagAnalyzerServiceClient
}

func New(config *config.Config, logger *slog.Logger) *Service {

	l := logger.With("service", "pricetaganalyzer")

	// conn := grpc.Dial(fmt.Sprintf("%s:%d", config.Services.Processing.Host, config.Services.PriceTagAnalyzer.Port), grpc.WithInsecure())

	host := config.Services.PriceTagAnalyzer.Host
	port := config.Services.PriceTagAnalyzer.Port

	l.Info("connecting to grpc service", slog.String("host", host), slog.Int("port", port))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		l.Error("error with connection to grpc service", slog.String("err", err.Error()))
		panic(err)
	}

	client := pricetaganalyzer.NewPriceTagAnalyzerServiceClient(conn)

	return &Service{
		config: config,
		client: client,
		logger: l,
	}
}

// Analyze implements controllers.AnalyzerService.
func (s *Service) Analyze(ctx context.Context, reader file.Reader) (*entity.ImageInfo, error) {

	rid := ctx.Value(middleware.REQUEST_ID)
	l := s.logger.With("request_id", rid)

	l.Debug("create stream for analyze image")
	stream, err := s.client.AnalyzeImage(ctx)
	if err != nil {
		return nil, err
	}

	chunkCount := int(math.Ceil(float64(reader.Size()) / (1 << 20)))
	l.Debug("chunk count", slog.Int("count", chunkCount))

	for i := 0; i < chunkCount; i++ {

		chunk := make([]byte, 1<<20)

		read, err := reader.Read(chunk)
		if err != nil {
			return nil, err
		}
		l.Debug("read bytes", slog.Int("read", read))

		l.Debug("sending image chunk", slog.Int("chunk", i))
		err = stream.Send(&pricetaganalyzer.ImageChunk{
			Content: chunk,
		})
		if err != nil {
			return nil, err
		}
		l.Debug("sent image chunk", slog.Int("chunk", i))
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {

		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.InvalidArgument {
				l.Error("grpc error", slog.String("err", e.String()))
				return nil, fmt.Errorf("invalid image")
			}
		}
		return nil, err
	}

	l.Debug("received reply", slog.String("reply", reply.String()))

	var measure *entity.Measure
	if reply.Measure != nil {
		measure = &entity.Measure{
			Amount: reply.Measure.Amount,
			Unit:   reply.Measure.Unit,
		}
	}

	info := &entity.ImageInfo{
		Product:     reply.Product,
		Price:       reply.Price,
		Description: reply.Description,
		Measure:     measure,
		Attributes:  reply.Attributes,
	}

	return info, nil
}
