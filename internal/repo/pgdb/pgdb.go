package pgdb

import (
	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/internal/events/outbox"
)

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

func (eo OutboxEvent) ToModel() outbox.OutboxEvent {
	res := outbox.OutboxEvent{
		ID:           eo.ID,
		Topic:        eo.Topic,
		EventType:    eo.EventType,
		EventVersion: eo.EventVersion,
		Key:          eo.Key,
		Payload:      eo.Payload,
		Status:       string(eo.Status),
		Attempts:     eo.Attempts,
		NextRetryAt:  eo.NextRetryAt,
		CreatedAt:    eo.CreatedAt,
	}
	if eo.SentAt.Valid {
		t := eo.SentAt.Time
		res.SentAt = &t
	}

	return res
}
