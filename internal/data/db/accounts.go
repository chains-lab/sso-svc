package db

import (
	"net/http"

	"github.com/cifra-city/sso-oauth/internal/data/db/sqlcore"
	"github.com/google/uuid"
)

type Accounts interface {
	Create(r *http.Request, email string, passHash string) (sqlcore.Account, error)

	GetById(r *http.Request, id uuid.UUID) (sqlcore.Account, error)
	GetByEmail(r *http.Request, email string) (sqlcore.Account, error)
}

type accounts struct {
	queries *sqlcore.Queries
}

func NewAccount(queries *sqlcore.Queries) Accounts {
	return &accounts{queries: queries}
}

func (a *accounts) Create(r *http.Request, email string, passHash string) (sqlcore.Account, error) {
	return a.queries.CreateAccount(r.Context(), email)
}

func (a *accounts) GetById(r *http.Request, id uuid.UUID) (sqlcore.Account, error) {
	return a.queries.GetAccountByID(r.Context(), id)
}

func (a *accounts) GetByEmail(r *http.Request, email string) (sqlcore.Account, error) {
	return a.queries.GetAccountByEmail(r.Context(), email)
}
