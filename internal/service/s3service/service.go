package s3service

import (
	"context"
	"fmt"
	"github.com/mzhn-sochi/gateway/api/s3"
	"github.com/mzhn-sochi/gateway/internal/config"
	"github.com/mzhn-sochi/gateway/pkg/file"
	"github.com/mzhn-sochi/gateway/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log/slog"
	"math"
)

type S3Service struct {
	config *config.Config
	client s3.S3Client
}

func New(config *config.Config, logger *slog.Logger) *S3Service {

	l := logger.With("service", "s3")

	host := config.Services.S3Service.Host
	port := config.Services.S3Service.Port

	l.Info("connecting to grpc service", slog.String("host", host), slog.Int("port", port))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		l.Error("error with connection to grpc service", slog.String("err", err.Error()))
		panic(err)
	}

	client := s3.NewS3Client(conn)

	return &S3Service{
		config: config,
		client: client,
	}
}

func (s *S3Service) Upload(ctx context.Context, reader file.Reader) (string, error) {

	l := ctx.Value(middleware.LOGGER).(*slog.Logger).With("service", "s3").With("method", "upload")

	chunkCount := int(math.Ceil(float64(reader.Size()) / (1 << 20)))
	l.Debug("chunk count", slog.Int("count", chunkCount))

	stream, err := s.client.Upload(ctx)
	if err != nil {
		return "", err
	}

	l.Debug(
		"send image",
		slog.Int("count", chunkCount),
		slog.String("contentType", reader.ContentType()),
		slog.Int64("size", reader.Size()),
	)

	for i := 0; i < chunkCount; i++ {

		chunk := make([]byte, 1<<20)

		read, err := reader.Read(chunk)
		if err != nil {
			return "", err
		}
		l.Debug("read bytes", slog.Int("read", read))

		l.Debug("sending image chunk", slog.Int("chunk", i))
		err = stream.Send(&s3.Object{
			Image: &s3.Image{
				Chunk: chunk,
			},
			Meta: &s3.Metadata{
				ContentType: reader.ContentType(),
			},
		})
		if err != nil {
			return "", err
		}
		l.Debug("sent image chunk", slog.Int("chunk", i))
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.InvalidArgument {
				l.Error("grpc error", slog.String("err", e.String()))
				return "", fmt.Errorf("invalid image")
			}
		}

		return "", err
	}

	l.Debug("received reply", slog.String("reply", reply.String()))

	return reply.Name, nil
}
