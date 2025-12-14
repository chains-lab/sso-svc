package pgdb

import (
	"github.com/chains-lab/sso-svc/internal/domain/entity"
)

//type txKeyType struct{}
//
//var TxKey = txKeyType{}
//
//func TxFromCtx(ctx context.Context) (*sql.Tx, bool) {
//	tx, ok := ctx.Value(TxKey).(*sql.Tx)
//	return tx, ok
//}

func (s Session) ToModel() entity.Session {
	return entity.Session{
		ID:        s.ID,
		AccountID: s.AccountID,
		LastUsed:  s.LastUsed,
		CreatedAt: s.CreatedAt,
	}
}

func (s Session) GetHashToken() string {
	return s.HashToken
}

func (a Account) ToModel() entity.Account {
	return entity.Account{
		ID:                a.ID,
		Username:          a.Username,
		Role:              string(a.Role),
		Status:            string(a.Status),
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
		UsernameUpdatedAt: a.UsernameUpdatedAt,
	}
}

func (ae AccountEmail) ToModel() entity.AccountEmail {
	return entity.AccountEmail{
		AccountID: ae.AccountID,
		Email:     ae.Email,
		Verified:  ae.Verified,
		UpdatedAt: ae.UpdatedAt,
		CreatedAt: ae.CreatedAt,
	}
}

func (ap AccountPassword) ToModel() entity.AccountPassword {
	return entity.AccountPassword{
		AccountID: ap.AccountID,
		Hash:      ap.Hash,
		UpdatedAt: ap.UpdatedAt,
		CreatedAt: ap.CreatedAt,
	}
}
