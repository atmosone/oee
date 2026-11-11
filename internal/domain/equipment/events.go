package equipment

import (
	"github.com/atmosone/cqrs"
)

const EventTypeCreated = cqrs.EventType("equipment_created")

type CreatedEvent struct {
	ID       string
	Name     string
	Category string
}

func (e CreatedEvent) Type() cqrs.EventType {
	return EventTypeCreated
}
