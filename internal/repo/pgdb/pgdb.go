package pgdb

import (
	"context"
	"database/sql"

	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/contracts"
)

type txKeyType struct{}

var TxKey = txKeyType{}

func TxFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(TxKey).(*sql.Tx)
	return tx, ok
}

func (a Account) ToEntity() entity.Account {
	return entity.Account{
		ID:                a.ID,
		Username:          a.Username,
		Role:              a.Role,
		Status:            a.Status,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
		UsernameUpdatedAt: a.UsernameUpdatedAt,
	}
}

func (a AccountPassword) ToEntity() entity.AccountPassword {
	return entity.AccountPassword{
		AccountID: a.AccountID,
		Hash:      a.Hash,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (ae AccountEmail) ToEntity() entity.AccountEmail {
	return entity.AccountEmail{
		AccountID: ae.AccountID,
		Email:     ae.Email,
		Verified:  ae.Verified,
		CreatedAt: ae.CreatedAt,
		UpdatedAt: ae.UpdatedAt,
	}
}

func (s Session) ToEntity() entity.Session {
	return entity.Session{
		ID:        s.ID,
		AccountID: s.AccountID,
		LastUsed:  s.LastUsed,
		CreatedAt: s.CreatedAt,
	}
}

func (eo OutboxEvent) ToModel() contracts.OutboxEvent {
	res := contracts.OutboxEvent{
		ID:           eo.ID,
		Topic:        eo.Topic,
		EventType:    eo.EventType,
		EventVersion: eo.EventVersion,
		Key:          eo.Key,
		Payload:      eo.Payload,
		Status:       eo.Status,
		Attempts:     eo.Attempts,
		NextRetryAt:  eo.NextRetryAt,
		CreatedAt:    eo.CreatedAt,
		SentAt:       eo.SentAt,
	}

	return res
}
