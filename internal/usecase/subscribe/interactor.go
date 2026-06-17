package subscribe

import (
	"context"
	"time"

	"github.com/sqdelitL/subscription-aggregator/internal/domain/port/respository"
	"github.com/sqdelitL/subscription-aggregator/internal/domain/subscribe"
)

type Interactor struct {
	writer respository.SubscriptionWriter
	reader respository.SubscriptionReader
}

func New(writer respository.SubscriptionWriter, reader respository.SubscriptionReader) *Interactor {
	return &Interactor{writer: writer, reader: reader}
}

func (i *Interactor) Create(ctx context.Context, sub *subscribe.Subscribe) (uint64, error) {
	if err := sub.Validate(); err != nil {
		return 0, err
	}

	return i.writer.Create(ctx, sub)
}

func (i *Interactor) Get(ctx context.Context, subID uint64) (*subscribe.Subscribe, error) {
	return i.reader.GetByID(ctx, subID)
}

func (i *Interactor) Update(ctx context.Context, sub *subscribe.Subscribe) error {
	if err := sub.Validate(); err != nil {
		return err
	}
	return i.writer.Update(ctx, sub)
}

func (i *Interactor) Delete(ctx context.Context, subID uint64) error {
	return i.writer.Delete(ctx, subID)
}

func (i *Interactor) List(ctx context.Context) ([]subscribe.Subscribe, error) {
	return i.reader.GetAll(ctx)
}

func (i *Interactor) GetTotalCost(ctx context.Context, filter respository.CostFilter) (int64, error) {
	subs, err := i.reader.FindSubscriptions(ctx, filter)
	if err != nil {
		return 0, err
	}
	return calculateTotal(subs, filter.Start, filter.End), nil
}

func calculateTotal(subs []subscribe.Subscribe, periodStart time.Time, periodEnd *time.Time) int64 {
	startMonth := monthNumber(periodStart)

	endMonth := 999999
	if periodEnd != nil {
		endMonth = monthNumber(*periodEnd) - 1
	}

	var total int64
	for _, sub := range subs {
		subStart := monthNumber(sub.StartDate)
		subEnd := 999999
		if sub.EndDate != nil {
			subEnd = monthNumber(*sub.EndDate)
		}

		activeStart := max(subStart, startMonth)
		activeEnd := min(subEnd, endMonth)

		if activeStart <= activeEnd {
			months := activeEnd - activeStart + 1
			total += sub.Price * int64(months)
		}
	}
	return total
}

func monthNumber(t time.Time) int {
	return t.Year()*12 + int(t.Month())
}
