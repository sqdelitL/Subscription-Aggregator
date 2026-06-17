package respository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sqdelitL/subscription-aggregator/internal/domain/subscribe"
)

type CostFilter struct {
	Start       time.Time
	End         *time.Time
	UserID      *uuid.UUID
	ServiceName *string
}

type SubscriptionWriter interface {
	Create(ctx context.Context, sub *subscribe.Subscribe) (uint64, error)
	Update(ctx context.Context, sub *subscribe.Subscribe) error
	Delete(ctx context.Context, id uint64) error
}

type SubscriptionReader interface {
	GetByID(ctx context.Context, id uint64) (*subscribe.Subscribe, error)
	GetAll(ctx context.Context) ([]subscribe.Subscribe, error)
	FindSubscriptions(ctx context.Context, filter CostFilter) ([]subscribe.Subscribe, error)
}
