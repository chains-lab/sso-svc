package api

import (
	"context"
	"fmt"
	"net"

	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/interceptors"
	"github.com/chains-lab/sso-svc/internal/api/service"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService interface {
	GetUser(context.Context, *emptypb.Empty) (*svc.User, error)
	RefreshToken(context.Context, *svc.RefreshTokenRequest) (*svc.TokensPair, error)
	GoogleLogin(context.Context, *emptypb.Empty) (*svc.GoogleLoginResponse, error)
	GoogleCallback(context.Context, *svc.GoogleCallbackRequest) (*svc.TokensPair, error)
	Logout(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	GetUserSession(context.Context, *emptypb.Empty) (*svc.Session, error)
	GetUserSessions(context.Context, *emptypb.Empty) (*svc.SessionsList, error)
	DeleteUserSession(context.Context, *emptypb.Empty) (*svc.SessionsList, error)
	TerminateUserSessions(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
}

type AdminService interface {
	GetUserByAdmin(context.Context, *svc.GetUserByAdminRequest) (*svc.User, error)
	CreateAdminByAdmin(context.Context, *svc.CreateAdminByAdminRequest) (*svc.User, error)
	GetUserSessionsByAdmin(context.Context, *svc.GetUserSessionsByAdminRequest) (*svc.SessionsList, error)
	TerminateUserSessionsByAdmin(context.Context, *svc.TerminateUserSessionsByAdminRequest) (*emptypb.Empty, error)
}

func Run(ctx context.Context, cfg config.Config, log *logrus.Logger, app *app.App) error {
	// 1) Создаём реализацию хэндлеров и interceptor
	server := service.NewService(cfg, app)
	authInterceptor := interceptors.NewAuth(cfg.JWT.Service.SecretKey, cfg.JWT.User.AccessToken.SecretKey)

	// 2) Инициализируем gRPC‐сервер
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)
	svc.RegisterUserServiceServer(grpcServer, server)
	svc.RegisterAdminServiceServer(grpcServer, server)

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
