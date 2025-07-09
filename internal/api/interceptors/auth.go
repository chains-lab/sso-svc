package interceptors

import (
	"context"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		toksServ := md["authorization"]
		if len(toksServ) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token not supplied")
		}

		data, err := tokens.VerifyServiceJWT(ctx, toksServ[0], skService)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		toksUser := md["x-user-token"]
		if len(toksUser) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "user token not supplied")
		}

		requestIDArr := md["x-request-id"]
		if len(requestIDArr) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "request ID not supplied")
		}

		userData, err := tokens.VerifyUserJWT(ctx, toksUser[0], skUser)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid user token: %v", err)
		}

		userID, err := uuid.Parse(userData.Subject)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid user ID: %v", err)
		}

		requestID, err := uuid.Parse(requestIDArr[0])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid request ID: %v", err)
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
