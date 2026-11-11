package equipment

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
)

var ErrEquipmentAlreadyExists = errors.New("equipment: already exists")

type ID string
type Name string
type Category string

type Equipment struct {
	ID       ID
	Name     Name
	Category Category
	Events   []cqrs.Event
}

func New(id ID, name Name, category Category) *Equipment {
	return &Equipment{
		ID:       id,
		Name:     name,
		Category: category,
		Events: []cqrs.Event{CreatedEvent{
			ID:       string(id),
			Name:     string(name),
			Category: string(category),
		}},
	}
}

type Saver interface {
	Save(context.Context, *Equipment) error
}
