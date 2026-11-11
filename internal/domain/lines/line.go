package lines

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
)

var ErrLineAlreadyExists = errors.New("lines: line already exists")
var ErrLineNotFound = errors.New("lines: line not found")
var ErrEquipmentAlreadyExists = errors.New("lines: equipment already exists")
var ErrEquipmentNotExists = errors.New("lines: equipment not exists")
var ErrUnknownDirection = errors.New("lines: unknown direction")
var ErrOutOfLine = errors.New("lines: out of line")

type ID string

type Name string

type Line struct {
	ID        ID
	Name      Name
	Equipment []equipment.ID
	Events    []cqrs.Event
}

func New(id ID, name Name) *Line {
	return &Line{
		ID:   id,
		Name: name,
		Events: []cqrs.Event{CreatedEvent{
			ID:   string(id),
			Name: string(name),
		}},
	}
}

func (l *Line) AddEquipment(equipmentID equipment.ID) error {
	for _, id := range l.Equipment {
		if id == equipmentID {
			return ErrEquipmentAlreadyExists
		}
	}
	l.Equipment = append(l.Equipment, equipmentID)
	l.Events = append(l.Events, EquipmentAddedEvent{
		ID:          string(l.ID),
		EquipmentID: string(equipmentID),
		Position:    len(l.Equipment) - 1,
	})
	return nil
}

type Direction string

const (
	DirectionRight = Direction("right")
	DirectionLeft  = Direction("left")
)

func (l *Line) MoveEquipment(equipmentID equipment.ID, direction Direction) error {
	currentIdx := -1
	for i, id := range l.Equipment {
		if id == equipmentID {
			currentIdx = i
			break
		}
	}
	if currentIdx == -1 {
		return ErrEquipmentNotExists
	}
	var targetIdx int
	switch direction {
	case DirectionLeft:
		targetIdx = currentIdx - 1
		if targetIdx < 0 {
			return ErrOutOfLine
		}
	case DirectionRight:
		targetIdx = currentIdx + 1
		if targetIdx >= len(l.Equipment) {
			return ErrOutOfLine
		}
	default:
		return ErrUnknownDirection
	}
	displacedEquipmentID := l.Equipment[targetIdx]
	l.Equipment[currentIdx], l.Equipment[targetIdx] = l.Equipment[targetIdx], l.Equipment[currentIdx]
	l.Events = append(l.Events,
		EquipmentMovedEvent{
			ID:          string(l.ID),
			EquipmentID: string(equipmentID),
			Position:    targetIdx,
		},
		EquipmentMovedEvent{
			ID:          string(l.ID),
			EquipmentID: string(displacedEquipmentID),
			Position:    currentIdx,
		},
	)
	return nil
}

func (l *Line) RemoveEquipment(equipmentID equipment.ID) error {
	targetIdx := -1
	for i, id := range l.Equipment {
		if id == equipmentID {
			targetIdx = i
			break
		}
	}
	if targetIdx == -1 {
		return ErrEquipmentNotExists
	}
	l.Events = append(l.Events, EquipmentRemovedEvent{
		ID:          string(l.ID),
		EquipmentID: string(equipmentID),
	})
	for i := targetIdx + 1; i < len(l.Equipment); i++ {
		shiftedEquipmentID := l.Equipment[i]
		newPosition := i - 1
		l.Events = append(l.Events, EquipmentMovedEvent{
			ID:          string(l.ID),
			EquipmentID: string(shiftedEquipmentID),
			Position:    newPosition,
		})
	}
	l.Equipment = append(l.Equipment[:targetIdx], l.Equipment[targetIdx+1:]...)
	return nil
}

type Reader interface {
	Line(context.Context, ID) (*Line, error)
}

type Saver interface {
	Save(context.Context, *Line) error
}
