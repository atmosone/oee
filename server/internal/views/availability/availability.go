package availability

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/downtimes"
)

type Workshift struct {
	ID         string
	Name       string
	StartedAt  time.Time
	FinishedAt time.Time
}

type Interval struct {
	StartedAt time.Time
	FinishedAt time.Time
}

func (i Interval) Allocation(other Interval) time.Duration {
	if i.FinishedAt.IsZero() || other.FinishedAt.IsZero() {
		return 0
	}
	var start, finish time.Time
	if i.StartedAt.After(other.StartedAt) {
		start = i.StartedAt
	} else {
		start = other.StartedAt
	}
	if i.FinishedAt.Before(other.FinishedAt) {
		finish = i.FinishedAt
	} else {
		finish = other.FinishedAt
	}
	if start.After(finish) {
		return 0
	}
	return finish.Sub(start)
}

type Allocation struct {
	WorkshiftID string
	EquipmentID string
	DowntimeID  string
	Allocation  time.Duration
}

type OnDowntimeStartedRepository interface {
	ActiveWorkshft(context.Context) (Workshift, error)
	SaveAllocation(context.Context, Allocation) error
}

func OnDowntimeStarted(repository OnDowntimeStartedRepository) cqrs.EventHandlerFunc {
	return func(ctx context.Context, e cqrs.Event) {
		evt := e.(downtimes.StartedEvent)
		workshift, err := repository.ActiveWorkshft(ctx)
		if err != nil {
			return
		}
		if err := repository.SaveAllocation(ctx, Allocation{
			WorkshiftID: workshift.ID,
			EquipmentID: evt.EquipmentID,
			DowntimeID:  evt.ID,
			Allocation:  0,
		}); err != nil {
			return
		}
	}
}

type OnDowntimeFinishedRepository interface {
	Allocations(context.Context, string) ([]Allocation, error)
	Workshift(context.Context, string) (Workshift, error)
}

func OnDowntimeFinished(repository OnDowntimeFinishedRepository) cqrs.EventHandlerFunc {
	return func(ctx context.Context, e cqrs.Event) {
		evt := e.(downtimes.FinishedEvent)
		allocations, err := repository.Allocations(ctx, evt.ID)
		if err != nil {
			return
		}
		for _, allocation := range allocations {
			workshift, err := repository.Workshift(ctx, allocation.WorkshiftID)
			if err != nil {
				return
			}
			if workshift.FinishedAt.Before(evt.FinishedAt) {
				allocation.Allocation = 
			}
		}
	}
}
