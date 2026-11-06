package specifications

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
)

var ErrAlreadyDeprecated = errors.New("specifications: already deprecated")

type ID struct{ value string }

func NewID(value string) (ID, error) {
	return ID{value: value}, nil
}

func (id ID) Value() string { return id.value }

type Code struct{ value string }

func NewCode(value string) (Code, error) {
	return Code{value: value}, nil
}

type Description struct{ value string }

func NewDescription(value string) (Description, error) {
	return Description{value: value}, nil
}

type Specification struct {
	id          ID
	code        Code
	description Description
	deprecated  bool
	events      []cqrs.Event
}

func New(id ID, code Code, description Description) *Specification {
	return &Specification{
		id:          id,
		code:        code,
		description: description,
		deprecated:  false,
		events: []cqrs.Event{CreatedEvent{
			ID:          id.value,
			Code:        code.value,
			Description: description.value,
			Deprecated:  false,
		}},
	}
}

func (s *Specification) Events() []cqrs.Event {
	events := make([]cqrs.Event, len(s.events))
	copy(events, s.events)
	return events
}

func (s *Specification) ClearEvents() { s.events = nil }

func (s *Specification) Deprecate() error {
	if s.deprecated {
		return ErrAlreadyDeprecated
	}
	s.deprecated = true
	s.events = append(s.events, DeprecatedEvent{
		ID:         s.id.value,
		Deprecated: s.deprecated,
	})
	return nil
}

type Reader interface {
	Specification(context.Context, ID) (*Specification, error)
}

type Saver interface {
	Save(context.Context, *Specification) error
}
