package lines

import "github.com/atmosone/cqrs"

const (
	EventTypeCreated          = cqrs.EventType("line_created")
	EventTypeEquipmentAdded   = cqrs.EventType("equipment_added_to_line")
	EventTypeEquipmentMoved   = cqrs.EventType("equipment_moved_in_line")
	EventTypeEquipmentRemoved = cqrs.EventType("equipment_removed_from_line")
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
