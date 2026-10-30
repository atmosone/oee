package domain

import (
	"fmt"
	"time"

	"github.com/atmosone/cqrs"
)

type Downtime struct {
	*reactiveEntity
	id          DowntimeID
	equipmentID EquipmentID
	status      DowntimeStatus
	duration    time.Duration
	startedAt   time.Time
	finishedAt  time.Time
}

func NewDowntime(id DowntimeID, equipmentID EquipmentID, startedAt time.Time) *Downtime {
	d := &Downtime{
		reactiveEntity: &reactiveEntity{},
		id:             id,
		equipmentID:    equipmentID,
		status:         DowntimeStatusActive,
		startedAt:      startedAt.UTC(),
	}
	d.events = append(d.events, DowntimeStartedEvent{
		DowntimeID:        string(d.id),
		EquipmentID:       string(d.equipmentID),
		DowntimeStatus:    string(d.status),
		DowntimeStartedAt: d.startedAt,
	})
	return d
}

func (d *Downtime) Finish(finishedAt time.Time) error {
	const shortDuration = 1 * time.Minute
	if d.status == DowntimeStatusFinished {
		return fmt.Errorf("%w (downtime already finished)", ErrConflict)
	}
	if finishedAt.Before(d.startedAt) {
		return fmt.Errorf("%w (downtime finish time should be after start time)", ErrUnprocessable)
	}
	d.status = DowntimeStatusFinished
	d.finishedAt = finishedAt.UTC()
	d.duration = d.finishedAt.Sub(d.startedAt)
	d.events = append(d.events, DowntimeFinishedEvent{
		DowntimeID:         string(d.id),
		EquipmentID:        string(d.equipmentID),
		DowntimeStatus:     string(d.status),
		DowntimeDuration:   d.duration,
		DowntimeFinishedAt: d.finishedAt,
	})
	return nil
}

type Specification struct {
	*reactiveEntity
	id          SpecificationID
	code        SpecificationCode
	description SpecificationDescription
}

func NewSpecification(
	id SpecificationID,
	code SpecificationCode,
	description SpecificationDescription,
) *Specification {
	s := &Specification{
		reactiveEntity: &reactiveEntity{},
		id:             id,
		code:           code,
		description:    description,
	}
	s.events = append(s.events, SpecificationCreatedEvent{
		SpecificationID:          string(s.id),
		SpecificationCode:        string(s.code),
		SpecificationDescription: string(s.description),
		SpecificationCreatedAt:   time.Now().UTC(),
	})
	return s
}

type QualificationRule struct {
	*reactiveEntity
	equipmentID     EquipmentID
	specificationID SpecificationID
	deprecated      bool
}

func NewQualificationRule(equipmentID EquipmentID, specificationID SpecificationID) *QualificationRule {
	qr := &QualificationRule{
		reactiveEntity:  &reactiveEntity{},
		equipmentID:     equipmentID,
		specificationID: specificationID,
		deprecated:      false,
	}
	qr.events = append(qr.events, QualificationRuleCreatedEvent{
		EquipmentID:                string(qr.equipmentID),
		SpecificationID:            string(qr.specificationID),
		QualificationRuleCreatedAt: time.Now().UTC(),
	})
	return qr
}

type Qualification struct {
	*reactiveEntity
	downtimeID      DowntimeID
	specificationID SpecificationID
	comment         QualificationComment
}

func NewQualification(
	downtimeID DowntimeID,
	specificationID SpecificationID,
	comment QualificationComment,
) *Qualification {
	q := &Qualification{
		reactiveEntity:  &reactiveEntity{},
		downtimeID:      downtimeID,
		specificationID: specificationID,
		comment:         comment,
	}
	q.events = append(q.events, QualificationCreatedEvent{
		DowntimeID:             string(q.downtimeID),
		SpecificationID:        string(q.specificationID),
		QualificationComment:   string(q.comment),
		QualificationCreatedAt: time.Now().UTC(),
	})
	return q
}

type reactiveEntity struct{ events []cqrs.Event }

func (e reactiveEntity) Events() []cqrs.Event {
	if l := len(e.events); l > 0 {
		events := make([]cqrs.Event, l)
		copy(events, e.events)
		return events
	}
	return nil
}

func (e *reactiveEntity) Clear() { e.events = nil }
