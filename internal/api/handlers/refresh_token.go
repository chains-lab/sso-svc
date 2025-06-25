package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/auth"
	"github.com/google/uuid"
)

func (a Service) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.TokensPairResponse, error) {
	requestID := uuid.New()
	meta := Meta(ctx)

	curToken := req.RefreshToken

	session, tokensPair, err := a.app.Refresh(ctx, meta.InitiatorID, meta.SessionID, req.Agent, curToken)
	if err != nil {
		return nil, responses.AppError(ctx, requestID, err)
	}

	Log(ctx, requestID).Infof("Session %s refreshed successfully", session.ID)

	return responses.TokensPair(tokensPair), nil
}
