package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const sessionsTable = "sessions"

type SessionModel struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	Token     string    `db:"token"`
	Client    string    `db:"client"`
	LastUsed  time.Time `db:"last_used"`
	CreatedAt time.Time `db:"created_at"`
}

type SessionsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewSessions(db *sql.DB) SessionsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return SessionsQ{
		db:       db,
		selector: builder.Select("*").From(sessionsTable),
		inserter: builder.Insert(sessionsTable),
		updater:  builder.Update(sessionsTable),
		deleter:  builder.Delete(sessionsTable),
		counter:  builder.Select("COUNT(*) AS count").From(sessionsTable),
	}
}

func (s SessionsQ) New() SessionsQ {
	return NewSessions(s.db)
}

type SessionInsertInput struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Token     string
	Client    string
	LastUsed  time.Time
	CreatedAt time.Time
}

func (s SessionsQ) Insert(ctx context.Context, input SessionInsertInput) error {
	values := map[string]interface{}{
		"id":         input.ID,
		"account_id": input.AccountID,
		"token":      input.Token,
		"client":     input.Client,
		"last_used":  input.LastUsed,
		"created_at": input.CreatedAt,
	}

	query, args, err := s.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for accounts: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = s.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

type SessionUpdateInput struct {
	Token    string
	Client   string
	LastUsed time.Time
}

func (s SessionsQ) Update(ctx context.Context, input SessionUpdateInput) error {
	values := map[string]interface{}{
		"token":     input.Token,
		"client":    input.Client,
		"last_used": input.LastUsed,
	}

	query, args, err := s.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for accounts: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = s.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (s SessionsQ) Delete(ctx context.Context) error {
	query, args, err := s.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for sessions: %w", err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = s.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (s SessionsQ) Select(ctx context.Context) ([]SessionModel, error) {
	query, args, err := s.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for sessions: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []SessionModel
	for rows.Next() {
		var sess SessionModel
		err = rows.Scan(
			&sess.ID,
			&sess.AccountID,
			&sess.Token,
			&sess.Client,
			&sess.CreatedAt,
			&sess.LastUsed,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning session row: %w", err)
		}
		sessions = append(sessions, sess)
	}
	return sessions, nil
}

func (s SessionsQ) Count(ctx context.Context) (int, error) {
	query, args, err := s.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for sessions: %w", err)
	}

	var count int
	err = s.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s SessionsQ) Get(ctx context.Context) (SessionModel, error) {
	query, args, err := s.selector.Limit(1).ToSql()
	if err != nil {
		return SessionModel{}, fmt.Errorf("building get query for sessions: %w", err)
	}

	row := s.db.QueryRowContext(ctx, query, args...)
	var sess SessionModel
	err = row.Scan(
		&sess.ID,
		&sess.AccountID,
		&sess.Token,
		&sess.Client,
		&sess.CreatedAt,
		&sess.LastUsed,
	)
	if err != nil {
		return SessionModel{}, err
	}
	return sess, nil
}

func (s SessionsQ) FilterID(id uuid.UUID) SessionsQ {
	s.selector = s.selector.Where(sq.Eq{"id": id})
	s.counter = s.counter.Where(sq.Eq{"id": id})
	s.deleter = s.deleter.Where(sq.Eq{"id": id})
	s.updater = s.updater.Where(sq.Eq{"id": id})
	return s
}

func (s SessionsQ) FilterAccountID(accountID uuid.UUID) SessionsQ {
	s.selector = s.selector.Where(sq.Eq{"account_id": accountID})
	s.counter = s.counter.Where(sq.Eq{"account_id": accountID})
	s.deleter = s.deleter.Where(sq.Eq{"account_id": accountID})
	s.updater = s.updater.Where(sq.Eq{"account_id": accountID})
	return s
}

func (s SessionsQ) Transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s SessionsQ) Page(limit, offset uint64) SessionsQ {
	s.counter = s.counter.Limit(limit).Offset(offset)
	s.selector = s.selector.Limit(limit).Offset(offset)
	return s
}
