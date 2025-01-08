package sqlcore

import (
	"context"
	"database/sql"
	"errors"
)

func (q *Queries) BeginTx(ctx context.Context) (*Queries, *sql.Tx, error) {
	db, ok := q.db.(*sql.DB)
	if !ok {
		return nil, nil, errors.New("underlying DBTX is not *sql.DB")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	return q.WithTx(tx), tx, nil
}
