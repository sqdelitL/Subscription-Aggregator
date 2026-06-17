package subscribe

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sqdelitL/subscription-aggregator/internal/domain/subscribe"
	"github.com/sqdelitL/subscription-aggregator/internal/errs"
)

type SubscriptionWriterImpl struct {
	db *sql.DB
}

func NewSubscriptionWriter(db *sql.DB) *SubscriptionWriterImpl {
	return &SubscriptionWriterImpl{db: db}
}

func (s *SubscriptionWriterImpl) Create(ctx context.Context, sub *subscribe.Subscribe) (uint64, error) {
	const query = `
		INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id uint64
	err := s.db.QueryRowContext(ctx, query,
		sub.UserID,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SubscriptionWriterImpl) Update(ctx context.Context, sub *subscribe.Subscribe) error {
	const query = `
		UPDATE subscriptions
		SET user_id = $1,
		    service_name = $2,
		    price = $3,
		    start_date = $4,
		    end_date = $5
		WHERE id = $6
	`
	
	result, err := s.db.ExecContext(ctx, query,
		sub.UserID,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.NotFoundSubscribeError
	}
	return nil
}

func (s *SubscriptionWriterImpl) Delete(ctx context.Context, id uint64) error {
	const query = `DELETE FROM subscriptions WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.NotFoundSubscribeError
	}
	return nil
}
