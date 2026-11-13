package analytics

import (
	"time"
)

type Production struct {
	Good   int64 `json:"good"`
	Defect int64 `json:"defect"`
}

func (p Production) Quality() float64 {
	if p.Good == 0 && p.Defect == 0 {
		return 1
	}
	if p.Good == 0 {
		return 0
	}
	return float64(p.Good+p.Defect) / float64(p.Good)
}

type Equipment struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type Workshift struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Interval
}

type Downtime struct {
	ID              string `json:"id"`
	SpecificationID string `json:"specification_id"`
	Interval
}

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

type TimelineItem struct {
	DowntimeID      string
	SpecificationID string
	Interval
}

type Report struct {
	Workshift  Workshift
	Equipment  Equipment
	Production []Production
	Downtimes  []Downtime
}

func (r Report) PlanDuration(now time.Time) time.Duration {
	return r.Workshift.Duration(now)
}

func (r Report) DowntimeDuration(now time.Time) time.Duration {
	total := time.Duration(0)
	for _, downtime := range r.Downtimes {
		total += r.Workshift.Allocation(now, downtime.Interval)
	}
	return total
}

func (r Report) Timeline(now time.Time) []TimelineItem {
	items := make([]TimelineItem, 0, len(r.Downtimes))
	for _, downtime := range r.Downtimes {
		inter, ok := r.Workshift.Intersection(now, downtime.Interval)
		if !ok {
			continue
		}
		items = append(items, TimelineItem{
			DowntimeID:      downtime.ID,
			SpecificationID: downtime.SpecificationID,
			Interval:        inter,
		})
	}
	return items
}

func (r Report) Pareto(now time.Time) map[string]time.Duration {
	pareto := make(map[string]time.Duration)
	for _, downtime := range r.Downtimes {
		dur := r.Workshift.Allocation(now, downtime.Interval)
		if dur > 0 {
			pareto[downtime.SpecificationID] += dur
		}
	}
	return pareto
}

func (r Report) Availability(now time.Time) float64 {
	plan := r.PlanDuration(now).Seconds()
	if plan <= 0 {
		return 1.0
	}
	avail := (plan - r.DowntimeDuration(now).Seconds()) / plan
	if avail < 0 {
		return 0
	}
	return avail
}
