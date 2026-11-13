package intervals

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
)

var ErrNotFound = errors.New("intervals: not found")
var ErrAlreadyExists = errors.New("intervals: already exists")
var ErrActiveNotExists = errors.New("intervals: active not exists")
var ErrAlreadyFinished = errors.New("intervals: already finished")
var ErrValidationFailed = errors.New("intervals: validation failed")

type ID string

type Status string

const (
	StatusActive   = Status("active")
	StatusFinished = Status("finished")
)

type Category string

const (
	CategoryDowntime = Category("downtime")
	CategoryWorktime = Category("worktime")
)

type Interval struct {
	ID          ID
	EquipmentID equipment.ID
	Status      Status
	Category    Category
	StartedAt   time.Time
	FinishedAt  time.Time
	Duration    time.Duration
}

type Equipment struct {
	ID     equipment.ID
	Active *Interval
	Events []cqrs.Event
}

func (e *Equipment) Finish(finishedAt time.Time) error {
	if e.Active == nil {
		return ErrActiveNotExists
	}
	if finishedAt.IsZero() {
		return fmt.Errorf("%w (invalid finish time value)", ErrValidationFailed)

	}
	if e.Active.StartedAt.After(finishedAt) {
		return fmt.Errorf("%w (expired finish time value)", ErrValidationFailed)
	}
	e.Active.Status = StatusFinished
	e.Active.FinishedAt = finishedAt.UTC()
	e.Active.Duration = e.Active.FinishedAt.Sub(e.Active.StartedAt)
	e.Events = append(e.Events, FinishedEvent{
		ID:          string(e.Active.ID),
		EquipmentID: string(e.ID),
		Status:      string(e.Active.Status),
		Category:    string(e.Active.Category),
		FinishedAt:  e.Active.FinishedAt,
		Duration:    e.Active.Duration,
	})
	e.Active = nil
	return nil
}

func (e *Equipment) StartInterval(id ID, category Category, startedAt time.Time) error {
	if id == "" {
		return fmt.Errorf("%w (empty id value)", ErrValidationFailed)
	}
	if category != CategoryDowntime && category != CategoryWorktime {
		return fmt.Errorf("%w (unknown category value: %s)", ErrValidationFailed, category)
	}
	if startedAt.IsZero() {
		return fmt.Errorf("%w (zero started at timestamp value)", ErrValidationFailed)
	}
	if e.Active != nil {
		if e.Active.Category == category {
			return fmt.Errorf("%w (no interval category switching)", ErrValidationFailed)
		}
		if err := e.Finish(startedAt); err != nil {
			return err
		}
	}
	e.Active = &Interval{
		ID:          id,
		EquipmentID: e.ID,
		Status:      StatusActive,
		Category:    category,
		StartedAt:   startedAt.UTC(),
	}
	e.Events = append(e.Events, StartedEvent{
		ID:          string(e.Active.ID),
		EquipmentID: string(e.ID),
		Status:      string(e.Active.Status),
		Category:    string(e.Active.Category),
		StartedAt:   e.Active.StartedAt,
	})
	return nil
}

type Saver interface {
	Save(context.Context, *Equipment) error
}

type Reader interface {
	Equipment(context.Context, equipment.ID) (*Equipment, error)
}
