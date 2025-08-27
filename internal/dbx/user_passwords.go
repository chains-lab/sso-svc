package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const usersPassTable = "user_passwords"

type UserPasswordModel struct {
	ID        uuid.UUID `db:"user_id"`
	PassHash  string    `db:"password_hash"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserPassQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewUsersPass(db *sql.DB) UserPassQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return UserPassQ{
		db:       db,
		selector: builder.Select("*").From(usersPassTable),
		inserter: builder.Insert(usersPassTable),
		updater:  builder.Update(usersPassTable),
		deleter:  builder.Delete(usersPassTable),
		counter:  builder.Select("COUNT(*) AS count").From(usersPassTable),
	}
}

func (q UserPassQ) applyConditions(conditions ...sq.Sqlizer) UserPassQ {
	q.selector = q.selector.Where(conditions)
	q.counter = q.counter.Where(conditions)
	q.updater = q.updater.Where(conditions)
	q.deleter = q.deleter.Where(conditions)

	return q
}

func (q UserPassQ) New() UserPassQ {
	return NewUsersPass(q.db)
}

func (q UserPassQ) Insert(ctx context.Context, input UserPasswordModel) error {
	values := map[string]interface{}{
		"user_id":       input.ID,
		"password_hash": input.PassHash,
		"updated_at":    input.UpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for table %s: %w", usersPassTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q UserPassQ) Update(ctx context.Context, input map[string]any) error {
	values := map[string]interface{}{}

	if passHash, ok := input["password_hash"]; ok {
		values["password_hash"] = passHash
	}
	if updatedAt, ok := input["updated_at"]; ok {
		values["updated_at"] = updatedAt
	} else {
		values["updated_at"] = time.Now().UTC()
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for table %s: %w", usersPassTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q UserPassQ) Get(ctx context.Context) (UserPasswordModel, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return UserPasswordModel{}, fmt.Errorf("building get query for table %s: %w", usersPassTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	var acc UserPasswordModel
	err = row.Scan(
		&acc.ID,
		&acc.PassHash,
		&acc.UpdatedAt,
	)
	if err != nil {
		return UserPasswordModel{}, err
	}

	return acc, nil
}

func (q UserPassQ) Select(ctx context.Context) ([]UserPasswordModel, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for table %s: %w", usersPassTable, err)
	}

	var rows *sql.Rows

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserPasswordModel
	for rows.Next() {
		var acc UserPasswordModel
		err := rows.Scan(
			&acc.ID,
			&acc.PassHash,
			&acc.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning table %s: %w", usersPassTable, err)
		}
		users = append(users, acc)
	}

	return users, nil
}

func (q UserPassQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for table %s: %w", usersPassTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}

	return nil
}

func (q UserPassQ) FilterID(userID uuid.UUID) UserPassQ {
	q.applyConditions(sq.Eq{"user_id": userID})

	return q
}

func (q UserPassQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %stable %s: %w", usersPassTable, err)
	}

	var count uint64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q UserPassQ) Page(limit, offset uint64) UserPassQ {
	q.counter = q.counter.Limit(limit).Offset(offset)
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}

func (q UserPassQ) Transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := q.db.BeginTx(ctx, nil)
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
