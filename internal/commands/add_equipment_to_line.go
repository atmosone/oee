package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/lines"
)

const CommandTypeAddEquipmentToLine = cqrs.CommandType("add_equipment_to_line")

type AddEquipmentToLineCommand struct {
	ID          lines.ID     `json:"id"`
	EquipmentID equipment.ID `json:"equipment_id"`
}

func (c AddEquipmentToLineCommand) Type() cqrs.CommandType {
	return CommandTypeAddEquipmentToLine
}

type AddEquipmentToLineRepositoryPort interface {
	lines.Reader
	lines.Saver
}

func AddEquipmentToLine(repository AddEquipmentToLineRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(AddEquipmentToLineCommand)
		line, err := repository.Line(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := line.AddEquipment(cmd.EquipmentID); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, line); err != nil {
			return nil, err
		}
		return line.Events, nil
	}
}
