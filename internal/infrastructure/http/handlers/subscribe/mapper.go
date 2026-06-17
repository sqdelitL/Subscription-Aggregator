package subscribe

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	domain "github.com/sqdelitL/subscription-aggregator/internal/domain/subscribe"
	"github.com/sqdelitL/subscription-aggregator/internal/errs"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/util"
)

func toDomain(sub *Subscribe) (*domain.Subscribe, error) {
	if sub == nil {
		return nil, fmt.Errorf("subscribe is nil. %w", errs.JsonValidationError)
	}
	var id uint64
	if sub.ID != 0 {
		id = sub.ID
	}

	userID, err := uuid.Parse(sub.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id. %w", errs.JsonValidationError)
	}

	start, err := util.ParseMMYYYY(sub.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date. %v. %w", err, errs.JsonValidationError)
	}

	var end *time.Time
	if sub.EndDate != nil {
		parsed, err := util.ParseMMYYYY(*sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date. %v. %w", err, errs.JsonValidationError)
		}
		end = &parsed
	}

	return &domain.Subscribe{
		ID:          id,
		UserID:      userID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

func fromDomain(d *domain.Subscribe) (*Subscribe, error) {
	if d == nil {
		return nil, fmt.Errorf("subscribe is nil. %w", errs.JsonValidationError)
	}

	var endStr *string
	if d.EndDate != nil {
		formatted := util.FormatMMYYYY(*d.EndDate)
		endStr = &formatted
	}

	return &Subscribe{
		ID:          d.ID,
		UserID:      d.UserID.String(),
		ServiceName: d.ServiceName,
		Price:       d.Price,
		StartDate:   util.FormatMMYYYY(d.StartDate),
		EndDate:     endStr,
	}, nil
}

func fromDomains(d []domain.Subscribe) ([]Subscribe, error) {
	subs := make([]Subscribe, 0, len(d))
	for _, sub := range d {
		s, err := fromDomain(&sub)
		if err != nil {
			return nil, err
		}
		subs = append(subs, *s)
	}
	return subs, nil
}
