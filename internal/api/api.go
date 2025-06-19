package api

import (
	"context"
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
	AdminDeleteUserSession(context.Context, *sso.AdminSessionRequest) (*sso.SessionsListResponse, error)
	AdminTerminateUserSessions(context.Context, *sso.AdminUserRequest) (*sso.Empty, error)
}

func Run(ctx context.Context, cfg config.Config, log *logrus.Logger, app *app.App) error {
	server := handlers.NewHandlers(cfg, log.WithField("service", "sso"), app)

	authInterceptor := interceptors.NewAuth(cfg.JWT.Service.SecretKey)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)

	sso.RegisterSsoServiceServer(grpcServer, server)

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		return err
	}
	log.Infof("gRPC server listening on %s", lis.Addr())
	return grpcServer.Serve(lis)
}
