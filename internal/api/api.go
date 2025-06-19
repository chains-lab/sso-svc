package api

import (
	"context"
	"fmt"
	"net"

	"github.com/chains-lab/chains-auth/internal/api/handlers"
	"github.com/chains-lab/chains-auth/internal/api/interceptors"
	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/config"
	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type SsoServiceServer interface {
	GetUser(context.Context, *sso.UserRequest) (*sso.UserResponse, error)
	RefreshToken(context.Context, *sso.RefreshTokenRequest) (*sso.TokensPairResponse, error)
	GoogleLogin(context.Context, *sso.Empty) (*sso.GoogleLoginResponse, error)
	GoogleCallback(context.Context, *sso.GoogleCallbackRequest) (*sso.TokensPairResponse, error)
	Logout(context.Context, *sso.SessionRequest) (*sso.Empty, error)
	GetUserSession(context.Context, *sso.SessionRequest) (*sso.SessionResponse, error)
	GetUserSessions(context.Context, *sso.UserRequest) (*sso.SessionsListResponse, error)
	DeleteUserSession(context.Context, *sso.SessionRequest) (*sso.SessionsListResponse, error)
	TerminateUserSessions(context.Context, *sso.UserRequest) (*sso.Empty, error)

	AdminTerminateUserSessions(context.Context, *sso.TerminateUserSessionByAdminRequest) (*sso.Empty, error)
	AdminUpdateUserRole(context.Context, *sso.UpdateUserRoleRequest) (*sso.Empty, error)
	AdminUpdateUserVerified(context.Context, *sso.UpdateUserVerifiedRequest) (*sso.Empty, error)
	AdminUpdateUserSuspended(context.Context, *sso.UpdateUserSuspendedRequest) (*sso.Empty, error)
	AdminUpdateUserSubscription(context.Context, *sso.UpdateUserSubscriptionRequest) (*sso.Empty, error)
}

func Run(ctx context.Context, cfg config.Config, log *logrus.Logger, app *app.App) error {
	// 1) Создаём реализацию хэндлеров и interceptor
	server := handlers.NewHandlers(cfg, log.WithField("service", "sso"), app)
	authInterceptor := interceptors.NewAuth(cfg.JWT.Service.SecretKey)

	// 2) Инициализируем gRPC‐сервер
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)
	sso.RegisterSsoServiceServer(grpcServer, server)

	// 3) Открываем слушатель
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	log.Infof("gRPC server listening on %s", lis.Addr())

	// 4) Запускаем Serve в горутине
	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- grpcServer.Serve(lis)
	}()

	// 5) Слушаем контекст и окончание Serve()
	select {
	case <-ctx.Done():
		// контекст отменили — аккуратно останавливаем
		log.Info("shutting down gRPC server …")
		grpcServer.GracefulStop()
		return nil
	case err := <-serveErrCh:
		// Serve() завершился (с ошибкой или EOF)
		return fmt.Errorf("gRPC Serve() exited: %w", err)
	}
}
