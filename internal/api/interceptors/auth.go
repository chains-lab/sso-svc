package interceptors

import (
	"context"

	"github.com/chains-lab/gatekit/tokens"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	IssuerKey    contextKey = "issuer"
	SubjectIDKey contextKey = "subject"
	AudienceKey  contextKey = "audience"
)

func NewAuth(secretKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		switch info.FullMethod {
		case "/sso.SsoService/GoogleLogin",
			"/sso.SsoService/GoogleCallback":
			// эти методы открыты — просто идём дальше без аутентификации
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		toks := md["authorization"]
		if len(toks) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token not supplied")
		}

		data, err := tokens.VerifyServiceJWT(ctx, toks[0], "your-secret-key")
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		ctx = context.WithValue(ctx, IssuerKey, data.Issuer)
		ctx = context.WithValue(ctx, SubjectIDKey, data.Subject)
		ctx = context.WithValue(ctx, AudienceKey, data.Audience)

		return handler(ctx, req)
	}
}

type TokenData struct {
	Issuer   string   `json:"iss,omitempty"`
	Subject  string   `json:"sub,omitempty"`
	Audience []string `json:"aud,omitempty"`
}

func GetTokenData(ctx context.Context) (TokenData, error) {
	issuer, ok := ctx.Value(IssuerKey).(string)
	if !ok {
		return TokenData{}, status.Errorf(codes.Unauthenticated, "issuer not found in context")
	}

	subject, ok := ctx.Value(SubjectIDKey).(string)
	if !ok {
		return TokenData{}, status.Errorf(codes.Unauthenticated, "subject not found in context")
	}

	audience, ok := ctx.Value(AudienceKey).([]string)
	if !ok {
		return TokenData{}, status.Errorf(codes.Unauthenticated, "audience not found in context")
	}

	return TokenData{
		Issuer:   issuer,
		Subject:  subject,
		Audience: audience,
	}, nil
}
