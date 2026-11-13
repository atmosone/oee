package production

import (
	"context"
	"errors"
	"fmt"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/intervals"
)

var ErrActiveIntervalNotFound = errors.New("production: active interval not found")
var ErrValidationFailed = errors.New("production: validation failed")
var ErrNoProductionInterval = errors.New("production: no production interval")

type ActiveInterval struct {
	ID          intervals.ID
	EquipmentID equipment.ID
	Category    intervals.Category
}

type Category string

const (
	CategoryGood   Category = "good"
	CategoryDefect Category = "defect"
)

type Amount int64

type Production struct {
	Category Category
	Amount   Amount
}

func NewProduction(category Category, amount Amount) (Production, error) {
	if category != CategoryGood && category != CategoryDefect {
		return Production{}, fmt.Errorf("%w (unknown category: %s)", ErrValidationFailed, category)
	}
	if amount < 0 {
		return Production{}, fmt.Errorf("%w (negative production value)", ErrValidationFailed)
	}
	return Production{Category: category, Amount: amount}, nil
}

func (p Production) Add(other Production) (Production, error) {
	if p.Category != other.Category {
		return Production{}, fmt.Errorf("%w (category mismatch)", ErrValidationFailed)
	}
	return Production{Amount: p.Amount + other.Amount, Category: p.Category}, nil
}

type Counter struct {
	Interval   ActiveInterval
	Production Production
	Events     []cqrs.Event
}

func (c *Counter) Increase(growth Production) error {
	if c.Interval.Category != intervals.CategoryWorktime {
		return ErrNoProductionInterval
	}
	if growth.Amount == 0 {
		return nil
	}
	next, err := c.Production.Add(growth)
	if err != nil {
		return err
	}
	c.Production = next
	c.Events = append(c.Events, IncreasedEvent{
		IntervalID:  string(c.Interval.ID),
		EquipmentID: string(c.Interval.EquipmentID),
		Category:    string(c.Production.Category),
		Amount:      int64(c.Production.Amount),
		Growth:      int64(growth.Amount),
	})
	return nil
}

const (
	EventTypeIncreased = cqrs.EventType("production_increased")
)

type IncreasedEvent struct {
	IntervalID  string `json:"interval_id"`
	EquipmentID string `json:"equipment_id"`
	Category    string `json:"category"`
	Amount      int64  `json:"amount"`
	Growth      int64  `json:"growth"`
}

func (e IncreasedEvent) Type() cqrs.EventType {
	return EventTypeIncreased
}

type Reader interface {
	Counter(context.Context, equipment.ID, Category) (*Counter, error)
}

type Saver interface {
	Save(context.Context, *Counter) error
}
