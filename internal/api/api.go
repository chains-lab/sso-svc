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
	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg config.Config, log logger.Logger, app *app.App) error {
	server := service.NewService(cfg, app)
	authInterceptor := interceptors.NewAuth(cfg.JWT.Service.SecretKey, cfg.JWT.User.AccessToken.SecretKey)
	logInterceptor := logger.UnaryLogInterceptor(log)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logInterceptor,
			authInterceptor,
		),
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
