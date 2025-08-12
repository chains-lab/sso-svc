package api

import (
	"context"
	"fmt"
	"net"

	sessProto "github.com/chains-lab/sso-proto/gen/go/svc/session"
	sessionAdminProto "github.com/chains-lab/sso-proto/gen/go/svc/sessionadmin"
	userProto "github.com/chains-lab/sso-proto/gen/go/svc/user"
	userAdminProto "github.com/chains-lab/sso-proto/gen/go/svc/useradmin"
	"github.com/chains-lab/sso-svc/internal/api/grpc/service/sessionadmin"
	"github.com/chains-lab/sso-svc/internal/api/grpc/service/useradmin"

	"github.com/chains-lab/sso-svc/internal/api/grpc/interceptor"
	"github.com/chains-lab/sso-svc/internal/api/grpc/service/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/service/user"
	"github.com/chains-lab/sso-svc/internal/app"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg config.Config, log logger.Logger, app *app.App) error {
	userSVC := user.NewService(cfg, app)
	sessionsSVC := session.NewService(cfg, app)
	sessionsAminSVC := sessionadmin.NewService(cfg, app)
	userAdminSVC := useradmin.NewService(cfg, app)
	authInterceptor := interceptor.Auth(cfg.JWT.Service.SecretKey)
	logInterceptor := logger.UnaryLogInterceptor(log)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logInterceptor,
			authInterceptor,
		),
	)

	sessProto.RegisterSessionServiceServer(grpcServer, sessionsSVC)
	userProto.RegisterUserServiceServer(grpcServer, userSVC)
	sessionAdminProto.RegisterSessionAdminServiceServer(grpcServer, sessionsAminSVC)
	userAdminProto.RegisterUserAdminServiceServer(grpcServer, userAdminSVC)

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
