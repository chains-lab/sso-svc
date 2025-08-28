package dbx

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const usersEmailTable = "user_emails"

type UserEmailModel struct {
	ID       uuid.UUID `db:"user_id"`
	Email    string    `db:"email"`
	Verified bool      `db:"verified"`
}

type UserEmailQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewUsersEmail(db *sql.DB) UserEmailQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return UserEmailQ{
		db:       db,
		selector: b.Select("*").From(usersEmailTable),
		inserter: b.Insert(usersEmailTable),
		updater:  b.Update(usersEmailTable),
		deleter:  b.Delete(usersEmailTable),
		counter:  b.Select("COUNT(*) AS count").From(usersEmailTable),
	}
}

func (q UserEmailQ) New() UserEmailQ {
	return NewUsersEmail(q.db)
}

func (q UserEmailQ) Insert(ctx context.Context, input UserEmailModel) error {
	values := map[string]any{
		"user_id":  input.ID,
		"email":    input.Email,
		"verified": input.Verified,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", usersEmailTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q UserEmailQ) Update(ctx context.Context, input map[string]any) error {
	values := map[string]any{}

	if email, ok := input["email"]; ok {
		values["email"] = email
	}
	if verified, ok := input["verified"]; ok {
		values["verified"] = verified
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", usersEmailTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q UserEmailQ) Get(ctx context.Context) (UserEmailModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return UserEmailModel{}, fmt.Errorf("building get query for %s: %w", usersEmailTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var m UserEmailModel
	err = row.Scan(&m.ID, &m.Email, &m.Verified)
	if err != nil {
		return UserEmailModel{}, err
	}
	return m, nil
}

func (q UserEmailQ) Select(ctx context.Context) ([]UserEmailModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", usersEmailTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []UserEmailModel
	for rows.Next() {
		var m UserEmailModel
		if err := rows.Scan(&m.ID, &m.Email, &m.Verified); err != nil {
			return nil, fmt.Errorf("scanning %s: %w", usersEmailTable, err)
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (q UserEmailQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", usersEmailTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q UserEmailQ) FilterID(userID uuid.UUID) UserEmailQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	return q
}

func (q UserEmailQ) FilterEmail(email string) UserEmailQ {
	q.selector = q.selector.Where(sq.Eq{"email": email})
	q.counter = q.counter.Where(sq.Eq{"email": email})
	q.updater = q.updater.Where(sq.Eq{"email": email})
	q.deleter = q.deleter.Where(sq.Eq{"email": email})
	return q
}

func (q UserEmailQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", usersEmailTable, err)
	}

	var count uint64
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q UserEmailQ) Page(limit, offset uint64) UserEmailQ {
	q.counter = q.counter.Limit(limit).Offset(offset)
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q UserEmailQ) Transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	ctxWithTx := context.WithValue(ctx, TxKey, tx)

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
