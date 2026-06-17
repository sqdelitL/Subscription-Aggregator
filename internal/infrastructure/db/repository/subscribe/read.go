package subscribe

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sqdelitL/subscription-aggregator/internal/domain/port/repository"
	"github.com/sqdelitL/subscription-aggregator/internal/domain/subscribe"
	"github.com/sqdelitL/subscription-aggregator/internal/errs"
)

type SubscriptionReaderImpl struct {
	db *sql.DB
}

func NewSubscriptionReader(db *sql.DB) *SubscriptionReaderImpl {
	return &SubscriptionReaderImpl{db: db}
}

func (s *SubscriptionReaderImpl) GetByID(ctx context.Context, id uint64) (*subscribe.Subscribe, error) {
	const query = `
		SELECT id, user_id, service_name, price, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`
	row := s.db.QueryRowContext(ctx, query, id)
	sub := &subscribe.Subscribe{}
	err := row.Scan(&sub.ID, &sub.UserID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFoundSubscribeError
		}
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionReaderImpl) GetAll(ctx context.Context) ([]subscribe.Subscribe, error) {
	const query = `
		SELECT id, user_id, service_name, price, start_date, end_date
		FROM subscriptions
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []subscribe.Subscribe
	for rows.Next() {
		var sub subscribe.Subscribe
		err := rows.Scan(&sub.ID, &sub.UserID, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return subs, nil
}

func (s *SubscriptionReaderImpl) FindSubscriptions(ctx context.Context, filter repository.CostFilter) ([]subscribe.Subscribe, error) {
	start := filter.Start

	var endParam time.Time
	if filter.End != nil {
		endParam = *filter.End
	} else {
		endParam = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
	}

	query := `
        SELECT id, user_id, service_name, price, start_date, end_date
        FROM subscriptions
        WHERE (end_date IS NULL OR end_date >= $1)
          AND start_date < $2
          AND ($3::uuid IS NULL OR user_id = $3)
          AND ($4::text IS NULL OR service_name = $4)
        ORDER BY start_date
    `

	rows, err := s.db.QueryContext(ctx, query,
		start,
		endParam,
		filter.UserID,
		filter.ServiceName,
	)
	if err != nil {
		return nil, fmt.Errorf("query subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []subscribe.Subscribe
	for rows.Next() {
		var sub subscribe.Subscribe
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.ServiceName, &sub.Price,
			&sub.StartDate, &sub.EndDate); err != nil {
			return nil, fmt.Errorf("scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}
	return subs, rows.Err()
}
