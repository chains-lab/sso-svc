package api

import (
	"context"
	"fmt"
	"net"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/handlers"
	"github.com/chains-lab/sso-svc/internal/api/interceptors"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/utils/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type AuthService interface {
	GetUser(context.Context, *svc.Empty) (*svc.UserResponse, error)
	RefreshToken(context.Context, *svc.RefreshTokenRequest) (*svc.TokensPairResponse, error)
	GoogleLogin(context.Context, *svc.Empty) (*svc.GoogleLoginResponse, error)
	GoogleCallback(context.Context, *svc.GoogleCallbackRequest) (*svc.TokensPairResponse, error)
	Logout(context.Context, *svc.Empty) (*svc.Empty, error)
	GetUserSession(context.Context, *svc.Empty) (*svc.SessionResponse, error)
	GetUserSessions(context.Context, *svc.Empty) (*svc.SessionsListResponse, error)
	DeleteUserSession(context.Context, *svc.Empty) (*svc.SessionsListResponse, error)
	TerminateUserSessions(context.Context, *svc.Empty) (*svc.Empty, error)

	AdminUpdateUserRole(context.Context, *svc.AdminUpdateUserRoleRequest) (*svc.UserResponse, error)
	AdminUpdateUserSubscription(context.Context, *svc.AdminUpdateUserSubscriptionRequest) (*svc.UserResponse, error)
	AdminUpdateUserSuspended(context.Context, *svc.AdminUpdateUserSuspendedRequest) (*svc.UserResponse, error)
	AdminUpdateUserVerified(context.Context, *svc.AdminUpdateUserVerifiedRequest) (*svc.UserResponse, error)
	AdminGetUserSessions(context.Context, *svc.AdminGetUserSessionsRequest) (*svc.SessionsListResponse, error)
	AdminGetUserSession(context.Context, *svc.AdminGetUserSessionRequest) (*svc.SessionResponse, error)
	AdminDeleteUserSession(context.Context, *svc.AdminDeleteUserSessionRequest) (*svc.Empty, error)
	AdminTerminateUserSessions(context.Context, *svc.AdminTerminateUserSessionsRequest) (*svc.Empty, error)
}

func Run(ctx context.Context, cfg config.Config, log *logrus.Logger, app *app.App) error {
	// 1) Создаём реализацию хэндлеров и interceptor
	server := handlers.NewService(cfg, app)
	authInterceptor := interceptors.NewAuth(cfg.JWT.Service.SecretKey)

	// 2) Инициализируем gRPC‐сервер
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)
	svc.RegisterServiceServer(grpcServer, server)

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
