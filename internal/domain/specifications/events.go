package specifications

import (
	"github.com/atmosone/cqrs"
)

const (
	EventTypeCreated    = cqrs.EventType("evt.specifications.created")
	EventTypeDeprecated = cqrs.EventType("evt.specifications.deprecated")
)

type CreatedEvent struct {
	ID          string
	Code        string
	Description string
	Deprecated  bool
}

func (e CreatedEvent) Type() cqrs.EventType {
	return EventTypeCreated
}

type DeprecatedEvent struct {
	ID         string
	Deprecated bool
}

func (e DeprecatedEvent) Type() cqrs.EventType {
	return EventTypeDeprecated
}
