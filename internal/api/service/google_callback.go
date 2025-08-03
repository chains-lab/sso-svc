package service

import (
	"context"
	"encoding/json"

	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func (s Service) GoogleCallback(
	ctx context.Context,
	req *svc.GoogleCallbackRequest,
) (*svc.TokensPair, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Log(ctx, uuid.Nil).Error("missing metadata in Google callback request")

		return nil, responses.UnauthorizedError(ctx, "metadata not found", nil)
	}

	requestIDArr := md["x-request-id"]
	if len(requestIDArr) == 0 {
		logger.Log(ctx, uuid.Nil).Error("missing request ID in Google callback request")

		return nil, responses.UnauthorizedError(ctx, "request ID not supplied", nil)
	}

	requestID, err := uuid.Parse(requestIDArr[0])
	if err != nil {
		logger.Log(ctx, uuid.Nil).WithError(err).Errorf("invalid request ID: %s", requestIDArr[0])

		return nil, responses.UnauthorizedError(ctx, "invalid request ID", nil)
	}

	code := req.Code
	if code == "" {
		logger.Log(ctx, requestID).Error("missing code in Google callback request")

		return nil, responses.BadRequestError(ctx, requestID, responses.Violation{
			Field:       "code",
			Description: "code is required",
		})
	}

	token, err := s.cfg.GoogleOAuth().Exchange(ctx, code)
	if err != nil {
		logger.Log(ctx, requestID).WithError(err).Errorf("error exchanging code for token: %s", code)

		return nil, responses.InternalError(ctx, &requestID)
	}

	client := s.cfg.GoogleOAuth().Client(ctx, token)
	httpResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logger.Log(ctx, requestID).WithError(err).Error("error fetching userinfo from Google")

		return nil, responses.InternalError(ctx, &requestID)
	}

	defer httpResp.Body.Close()

	var ui struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(httpResp.Body).Decode(&ui); err != nil {
		logger.Log(ctx, requestID).WithError(err).Error("error decoding Google userinfo")

		return nil, responses.InternalError(ctx, &requestID)
	}

	ua := ""
	if vals := md.Get("user-agent"); len(vals) > 0 {
		ua = vals[0]
	}

	_, tokensPair, err := s.app.Login(ctx, ui.Email, roles.User, ua)
	if err != nil {
		logger.Log(ctx, requestID).WithError(err).Errorf("error logging in user with email %s", ui.Email)

		return nil, responses.AppError(ctx, requestID, err)
	}

	return responses.TokensPair(tokensPair), nil
}
