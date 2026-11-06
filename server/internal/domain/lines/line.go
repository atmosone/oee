package lines

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
)

var ErrEquipmentAlreadyExists = errors.New("lines: equipment already exists")
var ErrEquipmentNotExists = errors.New("lines: equipment not exists")
var ErrUnknownDirection = errors.New("lines: unknown direction")
var ErrOutOfLine = errors.New("lines: out of line")

type ID struct{ value string }

type Name struct{ value string }

type Line struct {
	id        ID
	name      Name
	equipment []equipment.ID
	events    []cqrs.Event
}

func New(id ID, name Name) *Line {
	return &Line{
		id:   id,
		name: name,
		events: []cqrs.Event{CreatedEvent{
			ID:   id.value,
			Name: name.value,
		}},
	}
}

func (l *Line) Events() []cqrs.Event {
	events := make([]cqrs.Event, len(l.events))
	copy(events, l.events)
	return events
}

func (l *Line) ClearEvents() { l.events = nil }

func (l *Line) AddEquipment(equipmentID equipment.ID) error {
	for _, id := range l.equipment {
		if id.Equal(equipmentID) {
			return ErrEquipmentAlreadyExists
		}
	}
	l.equipment = append(l.equipment, equipmentID)
	l.events = append(l.events, EquipmentAddedEvent{
		ID:          l.id.value,
		EquipmentID: equipmentID.Value(),
		Position:    len(l.equipment) - 1,
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
	for i, id := range l.equipment {
		if id.Equal(equipmentID) {
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
		if targetIdx >= len(l.equipment) {
			return ErrOutOfLine
		}
	default:
		return ErrUnknownDirection
	}
	displacedEquipmentID := l.equipment[targetIdx]
	l.equipment[currentIdx], l.equipment[targetIdx] = l.equipment[targetIdx], l.equipment[currentIdx]
	l.events = append(l.events,
		EquipmentMovedEvent{
			ID:          l.id.value,
			EquipmentID: equipmentID.Value(),
			Position:    targetIdx,
		},
		EquipmentMovedEvent{
			ID:          l.id.value,
			EquipmentID: displacedEquipmentID.Value(),
			Position:    currentIdx,
		},
	)
	return nil
}

func (l *Line) RemoveEquipment(equipmentID equipment.ID) error {
	targetIdx := -1
	for i, id := range l.equipment {
		if id.Equal(equipmentID) {
			targetIdx = i
			break
		}
	}
	if targetIdx == -1 {
		return ErrEquipmentNotExists
	}
	l.events = append(l.events, EquipmentRemovedEvent{
		ID:          l.id.value,
		EquipmentID: equipmentID.Value(),
	})
	for i := targetIdx + 1; i < len(l.equipment); i++ {
		shiftedEquipmentID := l.equipment[i]
		newPosition := i - 1
		l.events = append(l.events, EquipmentMovedEvent{
			ID:          l.id.value,
			EquipmentID: shiftedEquipmentID.Value(),
			Position:    newPosition,
		})
	}
	l.equipment = append(l.equipment[:targetIdx], l.equipment[targetIdx+1:]...)
	return nil
}

type Reader interface {
	Line(context.Context, ID) (*Line, error)
}

type Saver interface {
	Save(context.Context, *Line) error
}
