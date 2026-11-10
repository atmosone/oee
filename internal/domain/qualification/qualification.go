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
	ID          downtimes.ID
	EquipmentID equipment.ID
	Status      downtimes.Status
}

func (d Downtime) IsFinished() bool {
	return d.Status == downtimes.StatusFinished
}

type Specification struct {
	ID         specifications.ID
	Deprecated bool
}

func (s Specification) IsDeprecated() bool {
	return s.Deprecated
}

type Qualification struct {
	EquipmentID    equipment.ID
	Specifications map[specifications.ID]Specification
	Events         []cqrs.Event
}

func (q *Qualification) RemoveRule(specificationID specifications.ID) error {
	specification, ok := q.Specifications[specificationID]
	if !ok {
		return ErrSpecificationNotExists
	}
	delete(q.Specifications, specification.ID)
	q.Events = append(q.Events, RuleRemovedEvent{
		EquipmentID:     string(q.EquipmentID),
		SpecificationID: string(specification.ID),
	})
	return nil
}

func (q *Qualification) CreateRule(specification Specification) error {
	if specification.Deprecated {
		return ErrDeprecatedSpecification
	}
	if _, ok := q.Specifications[specification.ID]; ok {
		return ErrSpecificationAlreadyExists
	}
	q.Specifications[specification.ID] = specification
	q.Events = append(q.Events, RuleCreatedEvent{
		EquipmentID:     string(q.EquipmentID),
		SpecificationID: string(specification.ID),
	})
	return nil
}

type Comment string

func (q *Qualification) QualifyDowntime(
	downtime Downtime,
	specificationID specifications.ID,
	comment Comment,
) error {
	if downtime.EquipmentID != q.EquipmentID {
		return ErrEquipmentConflict
	}
	if !downtime.IsFinished() {
		return ErrDowntimeNotFinishedYet
	}
	if specification, ok := q.Specifications[specificationID]; !ok {
		return ErrForbiddenSpecification
	} else if specification.IsDeprecated() {
		return ErrDeprecatedSpecification
	}
	q.Events = append(q.Events, DowntimeQualifiedEvent{
		DowntimeID:      string(downtime.ID),
		SpecificationID: string(specificationID),
		Comment:         string(comment),
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
