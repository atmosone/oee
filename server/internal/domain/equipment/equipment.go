package equipment

import (
	"context"

	"github.com/atmosone/cqrs"
)

type ID struct{ value string }

func NewID(value string) (ID, error) {
	return ID{value: value}, nil
}

func (id ID) Value() string { return id.value }

func (id ID) Equal(other ID) bool {
	return id.value == other.value
}

type Name struct{ value string }

func NewName(value string) (Name, error) {
	return Name{value: value}, nil
}

type Category struct{ value string }

func NewCategory(value string) (Category, error) {
	return Category{value: value}, nil
}

type Equipment struct {
	id       ID
	name     Name
	category Category
	events   []cqrs.Event
}

func (e *Equipment) Events() []cqrs.Event {
	events := make([]cqrs.Event, len(e.events))
	copy(events, e.events)
	return events
}

func (e *Equipment) ClearEvents() { e.events = nil }

func New(id ID, name Name, category Category) *Equipment {
	return &Equipment{
		id:       id,
		name:     name,
		category: category,
		events: []cqrs.Event{CreatedEvent{
			ID:       id.value,
			Name:     name.value,
			Category: category.value,
		}},
	}
}

type Saver interface {
	Save(context.Context, *Equipment) error
}
