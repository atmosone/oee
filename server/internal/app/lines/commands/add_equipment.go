package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/lines"
)

const CommandTypeAddEquipment = cqrs.CommandType("cmd.lines.add_equipment")

type AddEquipmentCommand struct {
	ID          lines.ID
	EquipmentID equipment.ID
}

func (c AddEquipmentCommand) Type() cqrs.CommandType {
	return CommandTypeAddEquipment
}

type AddEquipmentRepositoryPort interface {
	lines.Reader
	lines.Saver
}

func AddEquipment(repository AddEquipmentRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(AddEquipmentCommand)
		line, err := repository.Line(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := line.AddEquipment(cmd.EquipmentID); err != nil {
			return nil, err
		}
		defer line.ClearEvents()
		if err := repository.Save(ctx, line); err != nil {
			return nil, err
		}
		return line.Events(), nil
	}
}
