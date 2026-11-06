package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/lines"
)

const CommandTypeRemoveEquipment = cqrs.CommandType("cmd.lines.remove_equipment")

type RemoveEquipmentCommand struct {
	ID          lines.ID
	EquipmentID equipment.ID
}

func (c RemoveEquipmentCommand) Type() cqrs.CommandType {
	return CommandTypeRemoveEquipment
}

type RemoveEquipmentRepositoryPort interface {
	lines.Reader
	lines.Saver
}

func RemoveEquipment(repository RemoveEquipmentRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(RemoveEquipmentCommand)
		line, err := repository.Line(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := line.RemoveEquipment(cmd.EquipmentID); err != nil {
			return nil, err
		}
		defer line.ClearEvents()
		if err := repository.Save(ctx, line); err != nil {
			return nil, err
		}
		return line.Events(), nil
	}
}
