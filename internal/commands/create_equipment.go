package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
)

const CommandTypeCreateEquipment = cqrs.CommandType("cmd.equipment.create")

type CreateEquipmentCommand struct {
	ID       equipment.ID
	Name     equipment.Name
	Category equipment.Category
}

func (c CreateEquipmentCommand) Type() cqrs.CommandType {
	return CommandTypeCreateEquipment
}

type CreateEquipmentRepositoryPort interface {
	equipment.Saver
}

func CreateEquipment(repository CreateEquipmentRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateEquipmentCommand)
		equipment := equipment.New(
			cmd.ID,
			cmd.Name,
			cmd.Category,
		)
		if err := repository.Save(ctx, equipment); err != nil {
			return nil, err
		}
		return equipment.Events, nil
	}
}
