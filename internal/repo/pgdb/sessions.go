package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const sessionsTable = "sessions"

type Session struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	HashToken string    `db:"hash_token"`
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
		selector: builder.Select("sessions.*").From(sessionsTable),
		inserter: builder.Insert(sessionsTable),
		updater:  builder.Update(sessionsTable),
		deleter:  builder.Delete(sessionsTable),
		counter:  builder.Select("COUNT(*) AS count").From(sessionsTable),
	}
}

func (q SessionsQ) New() SessionsQ {
	return NewSessions(q.db)
}

func (q SessionsQ) Insert(ctx context.Context, input Session) error {
	values := map[string]interface{}{
		"id":         input.ID,
		"account_id": input.AccountID,
		"hash_token": input.HashToken,
		"last_used":  input.LastUsed,
		"created_at": input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", sessionsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q SessionsQ) Update(ctx context.Context) ([]Session, error) {
	q.updater = q.updater.
		Set("last_used", time.Now().UTC()).
		Suffix("RETURNING sessions.*")

	query, args, err := q.updater.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building update query for %s: %w", sessionsTable, err)
	}

	var rows *sql.Rows
	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Session
	for rows.Next() {
		var s Session
		err = rows.Scan(
			&s.ID,
			&s.AccountID,
			&s.HashToken,
			&s.LastUsed,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning updated session: %w", err)
		}
		out = append(out, s)
	}

	return out, nil
}

func (q SessionsQ) UpdateToken(token string) SessionsQ {
	q.updater = q.updater.Set("hash_token", token)
	return q
}

func (q SessionsQ) UpdateLastUsed(lastUsed time.Time) SessionsQ {
	q.updater = q.updater.Set("last_used", lastUsed)
	return q
}

func (q SessionsQ) Get(ctx context.Context) (Session, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Session{}, fmt.Errorf("building get query for sessions: %w", err)
	}

	var sess Session
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&sess.ID,
		&sess.AccountID,
		&sess.HashToken,
		&sess.CreatedAt,
		&sess.LastUsed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Session{}, nil
		}

		return Session{}, err
	}

	return sess, nil
}

func (q SessionsQ) Select(ctx context.Context) ([]Session, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for sessions: %w", err)
	}

	var rows *sql.Rows

	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var sess Session
		err = rows.Scan(
			&sess.ID,
			&sess.AccountID,
			&sess.HashToken,
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

func (q SessionsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for sessions: %w", err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q SessionsQ) FilterID(ID uuid.UUID) SessionsQ {
	q.selector = q.selector.Where(sq.Eq{"id": ID})
	q.deleter = q.deleter.Where(sq.Eq{"id": ID})
	q.updater = q.updater.Where(sq.Eq{"id": ID})
	q.counter = q.counter.Where(sq.Eq{"id": ID})

	return q
}

func (q SessionsQ) FilterAccountID(accountID uuid.UUID) SessionsQ {
	q.selector = q.selector.Where(sq.Eq{"account_id": accountID})
	q.deleter = q.deleter.Where(sq.Eq{"account_id": accountID})
	q.updater = q.updater.Where(sq.Eq{"account_id": accountID})
	q.counter = q.counter.Where(sq.Eq{"account_id": accountID})

	return q
}

func (q SessionsQ) OrderCreatedAt(ascending bool) SessionsQ {
	if ascending {
		q.selector = q.selector.OrderBy("created_at ASC")
	} else {
		q.selector = q.selector.OrderBy("created_at DESC")
	}
	return q
}

func (q SessionsQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for sessions: %w", err)
	}

	var count uint64
	if tx, ok := TxFromCtx(ctx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q SessionsQ) Page(limit, offset uint64) SessionsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}

func (q SessionsQ) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_, ok := TxFromCtx(ctx)
	if ok {
		return fn(ctx)
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
				err = fmt.Errorf("tx err: %v; rollback err: %v", err, rbErr)
			}
		}
	}()

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	if err = fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
