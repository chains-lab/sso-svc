package interceptor

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/sso-svc/internal/constant"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Auth(skService string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("no metadata found in incoming context")).Err()
		}
		tokenSvc := md["authorization"]
		if len(tokenSvc) == 0 {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("authorization token not supplied")).Err()
		}

		data, err := auth.VerifyServiceJWT(ctx, tokenSvc[0], skService)
		if err != nil {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("failed to verify token: %s", err)).Err()
		}

		requestIDArr := md["x-request-id"]
		if len(requestIDArr) == 0 {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("request ID not supplied")).Err()
		}

		requestID, err := uuid.Parse(requestIDArr[0])
		if err != nil {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("invalid request ID: %v", err)).Err()
		}

		ThisSvcInAudience := false
		for _, aud := range data.Audience {
			if aud == constant.ServiceName {
				ThisSvcInAudience = true
				break
			}
		}

		if !ThisSvcInAudience {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("service issuer %s not in audience %v", data.Issuer, data.Audience)).Err()
		}

		ctx = context.WithValue(ctx, RequestIDCtxKey, requestID)

		return handler(ctx, req)
	}
}
