package handlers

import (
	"context"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handlers) RefreshToken(ctx context.Context, req *sso.RefreshTokenRequest) (*sso.TokensPairResponse, error) {
	requestID := uuid.New()

	curToken := req.RefreshToken

	log := h.log.WithField("request_id", requestID)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid session id")
	}

	session, tokensPair, err := h.app.Refresh(ctx, userID, sessionID, req.Agent, curToken)
	if err != nil {
		return nil, h.presenter.AppError(requestID, err)
	}

	log.Infof("Session %s refreshed successfully", session.ID)

	return responses.TokensPair(tokensPair), nil
}
