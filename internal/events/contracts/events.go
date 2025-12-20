package contracts

import "github.com/umisto/sso-svc/internal/domain/entity"

type AccountCreatedPayload struct {
	Account entity.Account `json:"account"`
	Email   string         `json:"email,omitempty"`
}

const AccountCreatedEvent = "account.created"

const AccountLoginEvent = "account.login"

type AccountLoginPayload struct {
	Account entity.Account `json:"account"`
	Email   string         `json:"email"`
}

const AccountPasswordChangeEvent = "account.password.change"

type AccountPasswordChangePayload struct {
	Account entity.Account `json:"account"`
	Email   string         `json:"email"`
}

const AccountUsernameChangeEvent = "account.username.change"

type AccountUsernameChangePayload struct {
	Account entity.Account `json:"account"`
	Email   string         `json:"email"`
}
