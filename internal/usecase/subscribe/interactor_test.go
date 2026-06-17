package subscribe

import (
	"testing"
	"time"

	"github.com/sqdelitL/subscription-aggregator/internal/domain/subscribe"
)

func date(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func ptrDate(year int, month time.Month) *time.Time {
	t := date(year, month)
	return &t
}

func TestCalculateTotal(t *testing.T) {
	tests := []struct {
		name        string
		subs        []subscribe.Subscribe
		periodStart time.Time
		periodEnd   *time.Time
		expected    int64
	}{
		{
			name: "подписка полностью внутри периода",
			subs: []subscribe.Subscribe{
				{Price: 1000, StartDate: date(2025, 3), EndDate: ptrDate(2025, 5)},
			},
			periodStart: date(2025, 1),
			periodEnd:   ptrDate(2025, 7),
			expected:    3000,
		},
		{
			name: "частичное перекрытие слева (период начинается позже подписки)",
			subs: []subscribe.Subscribe{
				{Price: 500, StartDate: date(2025, 1), EndDate: ptrDate(2025, 6)},
			},
			periodStart: date(2025, 3),
			periodEnd:   ptrDate(2025, 9),
			expected:    2000,
		},
		{
			name: "частичное перекрытие справа (период заканчивается раньше подписки)",
			subs: []subscribe.Subscribe{
				{Price: 300, StartDate: date(2025, 5), EndDate: ptrDate(2025, 9)},
			},
			periodStart: date(2025, 1),
			periodEnd:   ptrDate(2025, 8),
			expected:    900,
		},
		{
			name: "подписка начинается после периода – пересечения нет",
			subs: []subscribe.Subscribe{
				{Price: 100, StartDate: date(2026, 1), EndDate: ptrDate(2026, 6)},
			},
			periodStart: date(2025, 1),
			periodEnd:   ptrDate(2025, 13),
			expected:    0,
		},
		{
			name: "подписка заканчивается до начала периода – пересечения нет",
			subs: []subscribe.Subscribe{
				{Price: 100, StartDate: date(2024, 1), EndDate: ptrDate(2024, 12)},
			},
			periodStart: date(2025, 1),
			periodEnd:   ptrDate(2025, 13),
			expected:    0,
		},
		{
			name: "открытая подписка (EndDate=nil) внутри закрытого периода",
			subs: []subscribe.Subscribe{
				{Price: 400, StartDate: date(2025, 7), EndDate: nil},
			},
			periodStart: date(2025, 1),
			periodEnd:   ptrDate(2025, 13),
			expected:    2400,
		},
		{
			name: "открытая подписка и открытый период (бесконечный)",
			subs: []subscribe.Subscribe{
				{Price: 200, StartDate: date(2025, 3), EndDate: nil},
			},
			periodStart: date(2025, 1),
			periodEnd:   nil,
			expected:    200 * int64(999999-monthNumber(date(2025, 3))+1),
		},
		{
			name: "несколько подписок с разным перекрытием",
			subs: []subscribe.Subscribe{
				{Price: 1000, StartDate: date(2025, 1), EndDate: ptrDate(2025, 6)},
				{Price: 500, StartDate: date(2025, 4), EndDate: ptrDate(2025, 9)},
			},
			periodStart: date(2025, 3),
			periodEnd:   ptrDate(2025, 8),
			expected:    1000*4 + 500*4,
		},
		{
			name:        "пустой список подписок",
			subs:        []subscribe.Subscribe{},
			periodStart: date(2025, 1),
			periodEnd:   ptrDate(2025, 13),
			expected:    0,
		},
		{
			name: "подписка ровно совпадает с периодом",
			subs: []subscribe.Subscribe{
				{Price: 700, StartDate: date(2025, 4), EndDate: ptrDate(2025, 6)},
			},
			periodStart: date(2025, 4),
			periodEnd:   ptrDate(2025, 7),
			expected:    2100,
		},
		{
			name: "одна подписка без конца, период раньше подписки",
			subs: []subscribe.Subscribe{
				{Price: 300, StartDate: date(2025, 6), EndDate: nil},
			},
			periodStart: date(2024, 1),
			periodEnd:   ptrDate(2024, 12),
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateTotal(tt.subs, tt.periodStart, tt.periodEnd)
			if got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}
