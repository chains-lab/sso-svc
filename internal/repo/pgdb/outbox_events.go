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

const outboxEventsTable = "outbox_events"

type OutboxEvent struct {
	ID           uuid.UUID `db:"id"`
	Topic        string    `db:"topic"`
	EventType    string    `db:"event_type"`
	EventVersion int       `db:"event_version"`
	Key          string    `db:"key"`
	Payload      []byte    `db:"payload"`

	Status      string     `db:"status"`
	Attempts    int        `db:"attempts"`
	NextRetryAt time.Time  `db:"next_retry_at"`
	CreatedAt   time.Time  `db:"created_at"`
	SentAt      *time.Time `db:"sent_at"`
}

type OutboxEventsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewOutboxEvents(db *sql.DB) OutboxEventsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return OutboxEventsQ{
		db:       db,
		selector: builder.Select("outbox_events.*").From(outboxEventsTable),
		inserter: builder.Insert(outboxEventsTable),
		updater:  builder.Update(outboxEventsTable),
		deleter:  builder.Delete(outboxEventsTable),
		counter:  builder.Select("COUNT(*) AS count").From(outboxEventsTable),
	}
}

func (q OutboxEventsQ) New() OutboxEventsQ {
	return NewOutboxEvents(q.db)
}

func (q OutboxEventsQ) Insert(ctx context.Context, input OutboxEvent) error {
	values := map[string]interface{}{
		"id":            input.ID,
		"topic":         input.Topic,
		"event_type":    input.EventType,
		"event_version": input.EventVersion,
		"key":           input.Key,
		"payload":       input.Payload,

		"status":        input.Status,
		"attempts":      input.Attempts,
		"next_retry_at": input.NextRetryAt,

		"created_at": input.CreatedAt,
		"sent_at":    input.SentAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", outboxEventsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q OutboxEventsQ) Update(ctx context.Context) ([]OutboxEvent, error) {
	q.updater = q.updater.Suffix("RETURNING outbox_events.*")

	query, args, err := q.updater.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building update query for %s: %w", outboxEventsTable, err)
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

	var out []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		err = rows.Scan(
			&e.ID,
			&e.Topic,
			&e.EventType,
			&e.EventVersion,
			&e.Key,
			&e.Payload,
			&e.Status,
			&e.Attempts,
			&e.NextRetryAt,
			&e.CreatedAt,
			&e.SentAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning updated outbox event: %w", err)
		}
		out = append(out, e)
	}

	return out, nil
}

func (q OutboxEventsQ) UpdateStatus(status string) OutboxEventsQ {
	q.updater = q.updater.Set("status", status)
	return q
}

func (q OutboxEventsQ) AddAttempts() OutboxEventsQ {
	q.updater = q.updater.Set("attempts", sq.Expr("attempts + 1"))
	return q
}

func (q OutboxEventsQ) UpdateAttempts(attempts int) OutboxEventsQ {
	q.updater = q.updater.Set("attempts", attempts)
	return q
}

func (q OutboxEventsQ) UpdateNextRetryAt(t time.Time) OutboxEventsQ {
	q.updater = q.updater.Set("next_retry_at", t)
	return q
}

func (q OutboxEventsQ) UpdateNextRetryAndStatus(nextRetryAt time.Time) OutboxEventsQ {
	q.updater = q.updater.Set("next_retry_at", nextRetryAt)
	return q
}

func (q OutboxEventsQ) UpdateSentAt(sentAt time.Time) OutboxEventsQ {
	q.updater = q.updater.Set("sent_at", sentAt)
	return q
}

func (q OutboxEventsQ) Get(ctx context.Context) (OutboxEvent, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return OutboxEvent{}, fmt.Errorf("building get query for %s: %w", outboxEventsTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var e OutboxEvent
	err = row.Scan(
		&e.ID,
		&e.Topic,
		&e.EventType,
		&e.EventVersion,
		&e.Key,
		&e.Payload,
		&e.Status,
		&e.Attempts,
		&e.NextRetryAt,
		&e.CreatedAt,
		&e.SentAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return OutboxEvent{}, nil
		}
		return OutboxEvent{}, err
	}

	return e, nil
}

func (q OutboxEventsQ) Select(ctx context.Context) ([]OutboxEvent, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", outboxEventsTable, err)
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

	var events []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		err = rows.Scan(
			&e.ID,
			&e.Topic,
			&e.EventType,
			&e.EventVersion,
			&e.Key,
			&e.Payload,
			&e.Status,
			&e.Attempts,
			&e.NextRetryAt,
			&e.CreatedAt,
			&e.SentAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning outbox event: %w", err)
		}
		events = append(events, e)
	}

	return events, nil
}

func (q OutboxEventsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", outboxEventsTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q OutboxEventsQ) FilterID(id uuid.UUID) OutboxEventsQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	return q
}

func (q OutboxEventsQ) FilterTopic(topic string) OutboxEventsQ {
	q.selector = q.selector.Where(sq.Eq{"topic": topic})
	q.counter = q.counter.Where(sq.Eq{"topic": topic})
	q.deleter = q.deleter.Where(sq.Eq{"topic": topic})
	q.updater = q.updater.Where(sq.Eq{"topic": topic})
	return q
}

func (q OutboxEventsQ) FilterEventType(eventType string) OutboxEventsQ {
	q.selector = q.selector.Where(sq.Eq{"event_type": eventType})
	q.counter = q.counter.Where(sq.Eq{"event_type": eventType})
	q.deleter = q.deleter.Where(sq.Eq{"event_type": eventType})
	q.updater = q.updater.Where(sq.Eq{"event_type": eventType})
	return q
}

func (q OutboxEventsQ) FilterStatus(status string) OutboxEventsQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	return q
}

func (q OutboxEventsQ) FilterReadyToSend(now time.Time) OutboxEventsQ {
	q.selector = q.selector.Where(sq.LtOrEq{"next_retry_at": now})
	q.counter = q.counter.Where(sq.LtOrEq{"next_retry_at": now})
	return q
}

func (q OutboxEventsQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", outboxEventsTable, err)
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

func (q OutboxEventsQ) Page(limit, offset uint64) OutboxEventsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q OutboxEventsQ) OrderByCreatedAtAsc() OutboxEventsQ {
	q.selector = q.selector.OrderBy("created_at ASC")
	return q
}

func (q OutboxEventsQ) OrderByCreatedAtDesc() OutboxEventsQ {
	q.selector = q.selector.OrderBy("created_at DESC")
	return q
}

func (q OutboxEventsQ) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
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
