package specifications

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
)

var ErrSpecificationNotFound = errors.New("specifications: specification not found")
var ErrAlreadyDeprecated = errors.New("specifications: already deprecated")

type ID string

type Code string

type Description string

type Specification struct {
	ID          ID
	Code        Code
	Description Description
	Deprecated  bool
	Events      []cqrs.Event
}

func New(id ID, code Code, description Description) *Specification {
	return &Specification{
		ID:          id,
		Code:        code,
		Description: description,
		Deprecated:  false,
		Events: []cqrs.Event{CreatedEvent{
			ID:          string(id),
			Code:        string(code),
			Description: string(description),
			Deprecated:  false,
		}},
	}
}

func (s *Specification) Deprecate() error {
	if s.Deprecated {
		return ErrAlreadyDeprecated
	}
	s.Deprecated = true
	s.Events = append(s.Events, DeprecatedEvent{
		ID:         string(s.ID),
		Deprecated: s.Deprecated,
	})
	return nil
}

type Reader interface {
	Specification(context.Context, ID) (*Specification, error)
}

type Saver interface {
	Save(context.Context, *Specification) error
}
