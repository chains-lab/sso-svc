package pgdb

import (
	"context"
	"database/sql"

	"github.com/umisto/sso-svc/internal/domain/entity"
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
