package api

import (
	"context"
	"fmt"
	"net"

	"github.com/chains-lab/chains-auth/internal/api/handlers"
	"github.com/chains-lab/chains-auth/internal/api/interceptors"
	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/chains-auth/internal/tools/config"
	"github.com/chains-lab/proto-storage/gen/go/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type AuthService interface {
	GetUser(context.Context, *auth.Empty) (*auth.UserResponse, error)
	RefreshToken(context.Context, *auth.RefreshTokenRequest) (*auth.TokensPairResponse, error)
	GoogleLogin(context.Context, *auth.Empty) (*auth.GoogleLoginResponse, error)
	GoogleCallback(context.Context, *auth.GoogleCallbackRequest) (*auth.TokensPairResponse, error)
	Logout(context.Context, *auth.Empty) (*auth.Empty, error)
	GetUserSession(context.Context, *auth.Empty) (*auth.SessionResponse, error)
	GetUserSessions(context.Context, *auth.Empty) (*auth.SessionsListResponse, error)
	DeleteUserSession(context.Context, *auth.Empty) (*auth.SessionsListResponse, error)
	TerminateUserSessions(context.Context, *auth.Empty) (*auth.Empty, error)

	AdminUpdateUserRole(context.Context, *auth.AdminUpdateUserRoleRequest) (*auth.UserResponse, error)
	AdminUpdateUserSubscription(context.Context, *auth.AdminUpdateUserSubscriptionRequest) (*auth.UserResponse, error)
	AdminUpdateUserSuspended(context.Context, *auth.AdminUpdateUserSuspendedRequest) (*auth.UserResponse, error)
	AdminUpdateUserVerified(context.Context, *auth.AdminUpdateUserVerifiedRequest) (*auth.UserResponse, error)
	AdminGetUserSessions(context.Context, *auth.AdminGetUserSessionsRequest) (*auth.SessionsListResponse, error)
	AdminGetUserSession(context.Context, *auth.AdminGetUserSessionRequest) (*auth.SessionResponse, error)
	AdminDeleteUserSession(context.Context, *auth.AdminDeleteUserSessionRequest) (*auth.Empty, error)
	AdminTerminateUserSessions(context.Context, *auth.AdminTerminateUserSessionsRequest) (*auth.Empty, error)
}

func Run(ctx context.Context, cfg config.Config, log *logrus.Logger, app *app.App) error {
	// 1) Создаём реализацию хэндлеров и interceptor
	server := handlers.NewService(cfg, app)
	authInterceptor := interceptors.NewAuth(cfg.JWT.Service.SecretKey)

	// 2) Инициализируем gRPC‐сервер
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)
	auth.RegisterAuthServiceServer(grpcServer, server)

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
		log.Info("shutting down gRPC server …")
		grpcServer.GracefulStop()
		return nil
	case err := <-serveErrCh:
		return fmt.Errorf("gRPC Serve() exited: %w", err)
	}
}
