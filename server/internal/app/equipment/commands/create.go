package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
)

const CommandTypeCreate = cqrs.CommandType("cmd.equipment.create")

type CreateCommand struct {
	ID       equipment.ID
	Name     equipment.Name
	Category equipment.Category
}

func (c CreateCommand) Type() cqrs.CommandType {
	return CommandTypeCreate
}

type CreateRepositoryPort interface {
	equipment.Saver
}

func Create(repository CreateRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateCommand)
		equipment := equipment.New(
			cmd.ID,
			cmd.Name,
			cmd.Category,
		)
		defer equipment.ClearEvents()
		if err := repository.Save(ctx, equipment); err != nil {
			return nil, err
		}
		return equipment.Events(), nil
	}
}
