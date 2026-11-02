package commands

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain"
)

const CommandTypeCreateEquipment = cqrs.CommandType("create_equipment")

type CreateEquipmentCommand struct {
	EquipmentID    string `json:"equipment_id"`
	EquipmentName  string `json:"equipment_name"`
	EquipmentUnits string `json:"equipment_units"`
}

func (c CreateEquipmentCommand) Type() cqrs.CommandType {
	return CommandTypeCreateEquipment
}

func (c CreateEquipmentCommand) Validate() error {
	return errors.Join(
		domain.EquipmentID(c.EquipmentID).Validate(),
		domain.EquipmentName(c.EquipmentName).Validate(),
		domain.EquipmentUnits(c.EquipmentUnits).Validate(),
	)
}

type CreateEquipmentRepository interface {
	domain.EquipmentSaver
}

func CreateEquipment(repository CreateEquipmentRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateEquipmentCommand)
		equipment := domain.NewEquipment(
			domain.EquipmentID(cmd.EquipmentID),
			domain.EquipmentName(cmd.EquipmentName),
			domain.EquipmentUnits(cmd.EquipmentUnits),
		)
		defer equipment.Clear()
		if err := repository.SaveEquipment(ctx, equipment); err != nil {
			return nil, err
		}
		return equipment.Events(), nil
	}
}
