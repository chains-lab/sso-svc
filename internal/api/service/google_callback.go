package service

import (
	"context"
	"encoding/json"

	"github.com/chains-lab/gatekit/roles"
	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s Service) GoogleCallback(
	ctx context.Context,
	req *svc.GoogleCallbackRequest,
) (*svc.TokensPair, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	requestIDArr := md["x-request-id"]
	if len(requestIDArr) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "request ID not supplied")
	}

	requestID, err := uuid.Parse(requestIDArr[0])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid request ID: %v", err)
	}

	code := req.Code
	if code == "" {
		Log(ctx, requestID).Error("missing code in Google callback request")

		return nil, status.Errorf(codes.InvalidArgument, "missing code")
	}

	token, err := s.cfg.GoogleOAuth().Exchange(ctx, code)
	if err != nil {
		Log(ctx, requestID).WithError(err).Errorf("error exchanging code for token: %s", code)

		return nil, status.Errorf(codes.Internal, "failed to exchange code")
	}

	client := s.cfg.GoogleOAuth().Client(ctx, token)
	httpResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		Log(ctx, requestID).WithError(err).Error("error fetching userinfo from Google")

		return nil, status.Errorf(codes.Internal, "failed to fetch user info")
	}
	defer httpResp.Body.Close()

	var ui struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(httpResp.Body).Decode(&ui); err != nil {
		Log(ctx, requestID).WithError(err).Error("error decoding Google userinfo")

		return nil, status.Errorf(codes.Internal, "invalid user info")
	}

	ua := ""
	if vals := md.Get("user-agent"); len(vals) > 0 {
		ua = vals[0]
	}

	_, tokensPair, err := s.app.Login(ctx, ui.Email, roles.User, ua)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.TokensPair(tokensPair), nil
}
