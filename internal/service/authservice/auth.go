package authservice

import (
	"context"
	"fmt"
	"github.com/mzhn-sochi/gateway/api/auth"
	"github.com/mzhn-sochi/gateway/internal/config"
	"github.com/mzhn-sochi/gateway/internal/entity"
	"github.com/mzhn-sochi/gateway/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log/slog"
)

type Service struct {
	config *config.Config
	client auth.AuthClient
}

//var _ controllers.AuthService = (*Service)(nil)

func New(config *config.Config, logger *slog.Logger) *Service {

	l := logger.With("service", "auth")

	host := config.Services.AuthService.Host
	port := config.Services.AuthService.Port

	l.Info("connecting to grpc service", slog.String("host", host), slog.Int("port", port))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		l.Error("error with connection to grpc service", slog.String("err", err.Error()))
		panic(err)
	}

	client := auth.NewAuthClient(conn)

	return &Service{
		config: config,
		client: client,
	}
}

func (s *Service) SignIn(ctx context.Context, phone, password string) (*entity.Tokens, error) {

	logger := ctx.Value(middleware.LOGGER).(*slog.Logger).With("service", "auth").With("method", "SignIn")

	req := &auth.SignInRequest{
		Phone:    phone,
		Password: password,
	}

	res, err := s.client.SignIn(ctx, req)
	if err != nil {
		logger.Debug("error with sign in", slog.String("err", err.Error()))
		return nil, err
	}

	return &entity.Tokens{
		Access:  res.Access,
		Refresh: res.Refresh,
	}, nil
}

func (s *Service) SignUp(ctx context.Context, phone, password string) (*entity.Tokens, error) {
	logger := ctx.Value(middleware.LOGGER).(*slog.Logger).With("service", "auth").With("method", "SignUp")

	req := &auth.SignUpRequest{
		Phone:    phone,
		Password: password,
	}

	res, err := s.client.SignUp(ctx, req)
	if err != nil {
		logger.Debug("error with sign up", slog.String("err", err.Error()))
		return nil, err
	}

	return &entity.Tokens{
		Access:  res.Access,
		Refresh: res.Refresh,
	}, nil
}

func (s *Service) SignOut(ctx context.Context, accessToken string) error {
	logger := ctx.Value(middleware.LOGGER).(*slog.Logger).With("service", "auth").With("method", "SignOut")

	req := &auth.SignOutRequest{
		AccessToken: accessToken,
	}

	if _, err := s.client.SignOut(ctx, req); err != nil {
		logger.Debug("error with sign out", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*entity.Tokens, error) {
	logger := ctx.Value(middleware.LOGGER).(*slog.Logger).With("service", "auth").With("method", "Refresh")

	req := &auth.RefreshRequest{
		RefreshToken: refreshToken,
	}

	res, err := s.client.Refresh(ctx, req)
	if err != nil {
		logger.Debug("error with refresh", slog.String("err", err.Error()))
		return nil, err
	}

	return &entity.Tokens{
		Access:  res.Access,
		Refresh: res.Refresh,
	}, nil
}

func (s *Service) Authenticate(ctx context.Context, accessToken string, role auth.Role) error {
	logger := ctx.Value(middleware.LOGGER).(*slog.Logger).With("service", "auth").With("method", "Auth")

	req := &auth.AuthRequest{
		AccessToken: accessToken,
		Role:        role,
	}

	logger.Debug("trying to auth", slog.String("role", role.String()), slog.String("accessToken", accessToken))

	if _, err := s.client.Auth(ctx, req); err != nil {

		if e, ok := status.FromError(err); ok {
			logger.Debug("error with auth",
				slog.String("err", e.Message()),
				slog.Int("code", int(e.Code())),
				slog.String("name", e.String()),
			)

			switch e.Code() {
			case codes.Unauthenticated:
				return ErrUnauthorized
			case codes.PermissionDenied:
				return ErrForbidden
			case codes.NotFound:
				return ErrNotFound
			case codes.InvalidArgument:
				return ErrInvalidRequest
			default:
				return err
			}
		}

		logger.Debug("error with auth", slog.String("err", err.Error()))
		return err
	}

	return nil
}
