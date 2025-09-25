package pgdb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/sso-svc/internal/data/schemas"
	"github.com/google/uuid"
)

const usersEmailTable = "users_email"

type UsersEmailQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewUsersEmail(db *sql.DB) UsersEmailQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return UsersEmailQ{
		db:       db,
		selector: b.Select("*").From(usersEmailTable),
		inserter: b.Insert(usersEmailTable),
		updater:  b.Update(usersEmailTable),
		deleter:  b.Delete(usersEmailTable),
		counter:  b.Select("COUNT(*) AS count").From(usersEmailTable),
	}
}

func (q UsersEmailQ) Insert(ctx context.Context, input schemas.UserEmailModel) error {
	values := map[string]any{
		"user_id":  input.ID,
		"email":    input.Email,
		"verified": input.Verified,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", usersEmailTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q UsersEmailQ) Update(ctx context.Context, input schemas.UserEmailUpdateInput) error {
	values := map[string]any{}

	if input.Email != nil {
		values["email"] = *input.Email
	}
	if input.Verified != nil {
		values["verified"] = *input.Verified
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", usersEmailTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q UsersEmailQ) Get(ctx context.Context) (schemas.UserEmailModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return schemas.UserEmailModel{}, fmt.Errorf("building get query for %s: %w", usersEmailTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var m schemas.UserEmailModel
	err = row.Scan(&m.ID, &m.Email, &m.Verified)
	if err != nil {
		return schemas.UserEmailModel{}, err
	}
	return m, nil
}

func (q UsersEmailQ) Select(ctx context.Context) ([]schemas.UserEmailModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", usersEmailTable, err)
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

	var out []schemas.UserEmailModel
	for rows.Next() {
		var m schemas.UserEmailModel
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

func (q UsersEmailQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", usersEmailTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q UsersEmailQ) FilterID(userID uuid.UUID) schemas.UsersEmailQ {
	q.selector = q.selector.Where(sq.Eq{"user_id": userID})
	q.counter = q.counter.Where(sq.Eq{"user_id": userID})
	q.updater = q.updater.Where(sq.Eq{"user_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"user_id": userID})
	return q
}

func (q UsersEmailQ) FilterEmail(email string) schemas.UsersEmailQ {
	q.selector = q.selector.Where(sq.Eq{"email": email})
	q.counter = q.counter.Where(sq.Eq{"email": email})
	q.updater = q.updater.Where(sq.Eq{"email": email})
	q.deleter = q.deleter.Where(sq.Eq{"email": email})
	return q
}

func (q UsersEmailQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", usersEmailTable, err)
	}

	var count uint
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

func (q UsersEmailQ) Page(limit, offset uint) schemas.UsersEmailQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}
