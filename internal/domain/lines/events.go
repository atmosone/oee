package lines

import "github.com/atmosone/cqrs"

const (
	EventTypeCreated          = cqrs.EventType("evt.lines.created")
	EventTypeEquipmentAdded   = cqrs.EventType("evt.lines.equipment_added")
	EventTypeEquipmentMoved   = cqrs.EventType("evt.lines.equipment_moved")
	EventTypeEquipmentRemoved = cqrs.EventType("evt.lines.equipment_removed")
)

type CreatedEvent struct {
	ID   string
	Name string
}

func (e CreatedEvent) Type() cqrs.EventType {
	return EventTypeCreated
}

type EquipmentAddedEvent struct {
	ID          string
	EquipmentID string
	Position    int
}

func (e EquipmentAddedEvent) Type() cqrs.EventType {
	return EventTypeEquipmentAdded
}

type EquipmentMovedEvent struct {
	ID          string
	EquipmentID string
	Position    int
}

func (e EquipmentMovedEvent) Type() cqrs.EventType {
	return EventTypeEquipmentMoved
}

type EquipmentRemovedEvent struct {
	ID          string
	EquipmentID string
}

func (e EquipmentRemovedEvent) Type() cqrs.EventType {
	return EventTypeEquipmentRemoved
}
