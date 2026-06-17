package subscribe

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/validation"
	"github.com/sqdelitL/subscription-aggregator/internal/errs"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/http/validate"
)

type Subscribe struct {
	ID          uint64  `json:"id"`
	UserID      string  `json:"user_id"`
	ServiceName string  `json:"service_name"`
	Price       int64   `json:"price"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

func (sub Subscribe) ValidateCreate() error {
	err := validation.ValidateStruct(&sub,
		validation.Field(&sub.UserID, validation.Required, validate.RuleUUID()),
		validation.Field(&sub.ServiceName, validation.Required),
		validation.Field(&sub.Price,
			validation.Min(0),
		),
		validation.Field(&sub.StartDate,
			validation.Required,
			validate.RuleSubscribeDateFormat(),
		),
		validation.Field(&sub.EndDate,
			validation.When(sub.EndDate != nil,
				validate.RuleSubscribeDateFormat(),
			),
		),
	)

	if err != nil {
		var i validation.InternalError
		if errors.As(err, &i) {
			return fmt.Errorf("error while executing validation Subscribe: %v: %w", i, errs.InternalError)
		}
		b, _ := json.Marshal(err)
		return fmt.Errorf("%s: %w", string(b), errs.JsonValidationError)
	}
	return nil
}

func (sub Subscribe) ValidateUpdate() error {
	err := validation.ValidateStruct(&sub,
		validation.Field(&sub.ID, validation.Required))
	if err != nil {
		var i validation.InternalError
		if errors.As(err, &i) {
			return fmt.Errorf("error while executing validation Subscribe: %v: %w", i, errs.InternalError)
		}
		b, _ := json.Marshal(err)
		return fmt.Errorf("%s: %w", string(b), errs.JsonValidationError)
	}

	return sub.ValidateCreate()
}
