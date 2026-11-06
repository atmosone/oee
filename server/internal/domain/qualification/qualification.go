package qualification

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/downtimes"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/specifications"
)

var (
	ErrEquipmentConflict          = errors.New("qualification: equipment conflict")
	ErrForbiddenSpecification     = errors.New("qualification: forbidden specification")
	ErrSpecificationNotExists     = errors.New("qualification: specification not exits")
	ErrDeprecatedSpecification    = errors.New("qualification: deprecated specification")
	ErrDowntimeNotFinishedYet     = errors.New("qualification: downtime not finished yet")
	ErrSpecificationAlreadyExists = errors.New("qualification: specification already exists")
)

type Downtime struct {
	id          downtimes.ID
	equipmentID equipment.ID
	status      downtimes.Status
}

func (d Downtime) EquipmentID() equipment.ID {
	return d.equipmentID
}

func (d Downtime) IsFinished() bool {
	return d.status == downtimes.StatusFinished
}

type Specification struct {
	id         specifications.ID
	deprecated bool
}

func (s Specification) IsDeprecated() bool {
	return s.deprecated
}

type EquipmentID string

type Qualification struct {
	equipmentID    equipment.ID
	specifications map[specifications.ID]Specification
	events         []cqrs.Event
}

func (q *Qualification) Events() []cqrs.Event {
	events := make([]cqrs.Event, len(q.events))
	copy(events, q.events)
	return events
}

func (q *Qualification) ClearEvents() { q.events = nil }

func (q *Qualification) RemoveRule(specificationID specifications.ID) error {
	specification, ok := q.specifications[specificationID]
	if !ok {
		return ErrSpecificationNotExists
	}
	delete(q.specifications, specification.id)
	q.events = append(q.events, RuleRemovedEvent{
		EquipmentID:     q.equipmentID.Value(),
		SpecificationID: specification.id.Value(),
	})
	return nil
}

func (q *Qualification) CreateRule(specification Specification) error {
	if specification.deprecated {
		return ErrDeprecatedSpecification
	}
	if _, ok := q.specifications[specification.id]; ok {
		return ErrSpecificationAlreadyExists
	}
	q.specifications[specification.id] = specification
	q.events = append(q.events, RuleCreatedEvent{
		EquipmentID:     q.equipmentID.Value(),
		SpecificationID: specification.id.Value(),
	})
	return nil
}

type Comment struct{ value string }

func NewComment(value string) (Comment, error) {
	return Comment{value: value}, nil
}

func (q *Qualification) QualifyDowntime(
	downtime Downtime,
	specificationID specifications.ID,
	comment Comment,
) error {
	if !downtime.equipmentID.Equal(q.equipmentID) {
		return ErrEquipmentConflict
	}
	if !downtime.IsFinished() {
		return ErrDowntimeNotFinishedYet
	}
	if specification, ok := q.specifications[specificationID]; !ok {
		return ErrForbiddenSpecification
	} else if specification.IsDeprecated() {
		return ErrDeprecatedSpecification
	}
	q.events = append(q.events, DowntimeQualifiedEvent{
		DowntimeID:      downtime.id.Value(),
		SpecificationID: specificationID.Value(),
		Comment:         comment.value,
	})
	return nil
}

type DowntimeReader interface {
	Downtime(context.Context, downtimes.ID) (Downtime, error)
}

type SpecificationReader interface {
	Specification(context.Context, specifications.ID) (Specification, error)
}

type Reader interface {
	Qualification(context.Context, equipment.ID) (*Qualification, error)
}

type Saver interface {
	Save(context.Context, *Qualification) error
}
