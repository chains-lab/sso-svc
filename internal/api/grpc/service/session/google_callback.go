package session

import (
	"context"
	"encoding/json"

	svc "github.com/chains-lab/sso-proto/gen/go/session"
	"github.com/chains-lab/sso-svc/internal/api/grpc/problems"
	"github.com/chains-lab/sso-svc/internal/api/grpc/responses"
	"github.com/chains-lab/sso-svc/internal/logger"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GoogleCallback(
	ctx context.Context,
	req *svc.GoogleCallbackRequest,
) (*svc.TokensPair, error) {
	code := req.Code
	if code == "" {
		logger.Log(ctx).Error("missing code in Google callback request")

		return nil, problems.InvalidArgumentError(ctx, "missing code in Google callback request", &errdetails.BadRequest_FieldViolation{
			Field:       "code",
			Description: "code is required",
		})
	}

	token, err := s.cfg.GoogleOAuth().Exchange(ctx, code)
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("error exchanging code for token: %s", code)

		return nil, problems.InternalError(ctx)
	}

	client := s.cfg.GoogleOAuth().Client(ctx, token)
	httpResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logger.Log(ctx).WithError(err).Error("error fetching userinfo from Google")

		return nil, problems.InternalError(ctx)
	}

	defer httpResp.Body.Close()

	var ui struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(httpResp.Body).Decode(&ui); err != nil {
		logger.Log(ctx).WithError(err).Error("error decoding Google userinfo")

		return nil, problems.InternalError(ctx)
	}

	_, tokensPair, err := s.app.Login(ctx, ui.Email, "TODO")
	if err != nil {
		logger.Log(ctx).WithError(err).Errorf("error logging in user with email %s", ui.Email)

		return nil, problems.AppError(ctx, err)
	}

	return responses.TokensPair(tokensPair), nil
}
