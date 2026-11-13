package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/production"
)

const CommandTypeIncreaseProduction = cqrs.CommandType("increase_production")

type IncreaseProductionCommand struct {
	EquipmentID equipment.ID        `json:"equipment_id"`
	Category    production.Category `json:"category"`
	Amount      production.Amount   `json:"amount"`
}

func (c IncreaseProductionCommand) Type() cqrs.CommandType {
	return CommandTypeIncreaseProduction
}

type IncreaseProductionRepository interface {
	production.Reader
	production.Saver
}

func IncreaseProduction(repository IncreaseProductionRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(IncreaseProductionCommand)
		growth, err := production.NewProduction(cmd.Category, cmd.Amount)
		if err != nil {
			return nil, err
		}
		counter, err := repository.Counter(ctx, cmd.EquipmentID, cmd.Category)
		if err != nil {
			return nil, err
		}
		if err := counter.Increase(growth); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, counter); err != nil {
			return nil, err
		}
		return counter.Events, nil
	}
}
