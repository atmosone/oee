package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/lines"
)

const CommandTypeRemoveEquipmentFromLine = cqrs.CommandType("remove_equipment_from_line")

type RemoveEquipmentFromLineCommand struct {
	ID          lines.ID     `json:"id"`
	EquipmentID equipment.ID `json:"equipment_id"`
}

func (c RemoveEquipmentFromLineCommand) Type() cqrs.CommandType {
	return CommandTypeRemoveEquipmentFromLine
}

type RemoveEquipmentFromLineRepositoryPort interface {
	lines.Reader
	lines.Saver
}

func RemoveEquipmentFromLine(repository RemoveEquipmentFromLineRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(RemoveEquipmentFromLineCommand)
		line, err := repository.Line(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := line.RemoveEquipment(cmd.EquipmentID); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, line); err != nil {
			return nil, err
		}
		return line.Events, nil
	}
}
