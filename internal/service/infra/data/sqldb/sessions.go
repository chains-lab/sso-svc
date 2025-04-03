package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/hs-zavet/sso-oauth/internal/service/domain/models"
)

const sessionsTable = "sessions"

type Sessions interface {
	New() Sessions

	Insert(ctx context.Context, sess models.Session) error
	Update(ctx context.Context, updates map[string]interface{}) error
	Delete(ctx context.Context) error

	Select(ctx context.Context) ([]models.Session, error)
	Count(ctx context.Context) (int, error)
	Get(ctx context.Context) (*models.Session, error)

	Filter(filters map[string]interface{}) Sessions

	Transaction(fn func(ctx context.Context) error) error

	Page(limit, offset uint64) Sessions
}

type sessions struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewSessions(db *sql.DB) Sessions {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &sessions{
		db:       db,
		selector: builder.Select("*").From(sessionsTable),
		inserter: builder.Insert(sessionsTable),
		updater:  builder.Update(sessionsTable),
		deleter:  builder.Delete(sessionsTable),
		counter:  builder.Select("COUNT(*) AS count").From(sessionsTable),
	}
}

func (s *sessions) New() Sessions {
	return NewSessions(s.db)
}

func (s *sessions) Insert(ctx context.Context, sess models.Session) error {
	values := map[string]interface{}{
		"id":         sess.ID,
		"account_id": sess.AccountID,
		"token":      sess.Token,
		"client":     sess.Client,
		"ip":         sess.IP,
	}
	values["created_at"] = time.Now().UTC()
	values["last_used"] = time.Now().UTC()

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

func (s *sessions) Update(ctx context.Context, updates map[string]any) error {
	updates["last_used"] = time.Now().UTC()
	query, args, err := s.updater.
		SetMap(updates).
		Where(sq.Eq{"id": updates["id"]}).
		ToSql()
	if err != nil {
		return fmt.Errorf("building update query for sessions: %w", err)
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

func (s *sessions) Delete(ctx context.Context) error {
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

func (s *sessions) Select(ctx context.Context) ([]models.Session, error) {
	query, args, err := s.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for sessions: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var sess models.Session
		err = rows.Scan(
			&sess.ID,
			&sess.AccountID,
			&sess.Token,
			&sess.Client,
			&sess.IP,
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

func (s *sessions) Count(ctx context.Context) (int, error) {
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

func (s *sessions) Get(ctx context.Context) (*models.Session, error) {
	query, args, err := s.selector.Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("building get query for sessions: %w", err)
	}

	row := s.db.QueryRowContext(ctx, query, args...)
	var sess models.Session
	err = row.Scan(
		&sess.ID,
		&sess.AccountID,
		&sess.Token,
		&sess.Client,
		&sess.IP,
		&sess.CreatedAt,
		&sess.LastUsed,
	)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (s *sessions) Filter(filters map[string]any) Sessions {
	var validFilters = map[string]bool{
		"id":         true,
		"account_id": true,
		"token":      true,
		"client":     true,
		"ip":         true,
	}

	for key, value := range filters {
		if _, exits := validFilters[key]; !exits {
			continue
		}
		s.selector = s.selector.Where(sq.Eq{key: value})
		s.counter = s.counter.Where(sq.Eq{key: value})
		s.deleter = s.deleter.Where(sq.Eq{key: value})
		s.updater = s.updater.Where(sq.Eq{key: value})
	}
	return s
}

func (s *sessions) Transaction(fn func(ctx context.Context) error) error {
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

func (s *sessions) Page(limit, offset uint64) Sessions {
	s.counter = s.counter.Limit(limit).Offset(offset)
	s.selector = s.selector.Limit(limit).Offset(offset)
	return s
}
