package interceptors

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type MetaData struct {
	Issuer         string     `json:"iss,omitempty"`
	Subject        string     `json:"sub,omitempty"`
	Audience       []string   `json:"aud,omitempty"`
	InitiatorID    uuid.UUID  `json:"initiator_id,omitempty"`
	SessionID      uuid.UUID  `json:"session_id,omitempty"`
	SubscriptionID uuid.UUID  `json:"subscription_id,omitempty"`
	Verified       bool       `json:"verified,omitempty"`
	Role           roles.Role `json:"role,omitempty"`
	RequestID      uuid.UUID  `json:"request_id,omitempty"`
}

func NewAuth(skService, skUser string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		switch info.FullMethod {
		case "/sso.UserService/GoogleLogin",
			"/sso.UserService/GoogleCallback":
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, responses.UnauthorizedError(ctx, "metadata not found", nil)
		}
		toksServ := md["authorization"]
		if len(toksServ) == 0 {
			return nil, responses.UnauthorizedError(ctx, "authorization token not supplied", nil)
		}

		data, err := tokens.VerifyServiceJWT(ctx, toksServ[0], skService)
		if err != nil {
			return nil, responses.UnauthorizedError(ctx, fmt.Sprintf("failed to verify token: %s", err), nil)
		}

		toksUser := md["x-user-token"]
		if len(toksUser) == 0 {
			return nil, responses.UnauthorizedError(ctx, "user token not supplied", nil)
		}

		requestIDArr := md["x-request-id"]
		if len(requestIDArr) == 0 {
			return nil, responses.UnauthorizedError(ctx, "request ID not supplied", nil)
		}

		userData, err := tokens.VerifyUserJWT(ctx, toksUser[0], skUser)
		if err != nil {
			return nil, responses.UnauthorizedError(ctx, fmt.Sprintf("invalid user token: %v", err), nil)
		}

		userID, err := uuid.Parse(userData.Subject)
		if err != nil {
			return nil, responses.UnauthorizedError(ctx, fmt.Sprintf("invalid user ID: %v", err), nil)
		}

		requestID, err := uuid.Parse(requestIDArr[0])
		if err != nil {
			return nil, responses.UnauthorizedError(ctx, fmt.Sprintf("invalid request ID: %v", err), nil)
		}

		ctx = context.WithValue(ctx, MetaCtxKey, MetaData{
			Issuer:         data.Issuer,
			Subject:        data.Subject,
			Audience:       data.Audience,
			InitiatorID:    userID,
			SessionID:      userData.Session,
			SubscriptionID: userData.Subscription,
			Verified:       userData.Verified,
			Role:           userData.Role,
			RequestID:      requestID,
		})

		return handler(ctx, req)
	}
}
