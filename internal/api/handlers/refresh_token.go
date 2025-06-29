package handlers

import (
	"context"

	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/api/responses"
	"github.com/google/uuid"
)

func (s Service) RefreshToken(ctx context.Context, req *svc.RefreshTokenRequest) (*svc.TokensPairResponse, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	curToken := req.RefreshToken

	session, tokensPair, err := s.app.Refresh(ctx, meta.InitiatorID, meta.SessionID, req.Agent, curToken)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Infof("Session %s refreshed successfully", session.ID)

	return responses.TokensPair(tokensPair), nil
}
