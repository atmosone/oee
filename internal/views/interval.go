package views

import "time"

type Interval struct {
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
}

func (i Interval) Start() time.Time { return i.StartedAt }

func (i Interval) Finish(now time.Time) time.Time {
	if i.FinishedAt.IsZero() {
		return now
	}
	return i.FinishedAt
}

func (i Interval) Duration(now time.Time) time.Duration {
	d := i.Finish(now).Sub(i.StartedAt)
	if d < 0 {
		return 0
	}
	return d
}

func (i Interval) Intersection(now time.Time, other Interval) (Interval, bool) {
	var start, finish time.Time

	if i.Start().After(other.Start()) {
		start = i.Start()
	} else {
		start = other.Start()
	}

	if i.Finish(now).Before(other.Finish(now)) {
		finish = i.Finish(now)
	} else {
		finish = other.Finish(now)
	}

	if !finish.After(start) {
		return Interval{}, false
	}

	return Interval{
		StartedAt:  start,
		FinishedAt: finish,
	}, true
}

func (i Interval) Allocation(now time.Time, other Interval) time.Duration {
	inter, ok := i.Intersection(now, other)
	if !ok {
		return 0
	}
	return inter.Duration(now)
}
