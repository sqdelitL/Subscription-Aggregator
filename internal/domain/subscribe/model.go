package subscribe

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqdelitL/subscription-aggregator/internal/errs"
)

type Subscribe struct {
	ID          uint64
	UserID      uuid.UUID
	ServiceName string
	Price       int64
	StartDate   time.Time
	EndDate     *time.Time
}

func (sub *Subscribe) Validate() error {
	if sub.Price < 0 {
		return errs.NegativeSubscribePriceError
	}

	if sub.EndDate != nil {
		if sub.StartDate.After(*sub.EndDate) {
			return errs.WrongSubscribeDateIntervalError
		}
	}
	return nil
}
