package callbacks

import (
	"context"
	"encoding/json"
	"fmt"

	//"database/sql"
	//
	//"github.com/recovery-flow/comtools/httpkit"
	//"github.com/recovery-flow/comtools/httpkit/problems"

	"github.com/google/uuid"
	"github.com/recovery-flow/sso-oauth/internal/service"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/events/rabbit/evebody"
)

func AccountUpdateRole(ctx context.Context, svc *service.Service, body []byte) error {
	if svc == nil || svc.Domain == nil {
		return fmt.Errorf("service or domain layer is nil")
	}

	var event evebody.RoleUpdated
	err := json.Unmarshal(body, &event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event body: %w", err)
	}

	ID, err := uuid.Parse(event.AccountID)
	if err != nil {
		return fmt.Errorf("failed to parse account ID: %w", err)
	}

	return svc.Domain.SessionsTerminate(ctx, ID, nil)
}
