package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/qualification"
	"github.com/atmosone/oee/internal/domain/specifications"
)

const CommandTypeRemoveQualificationRule = cqrs.CommandType("cmd.qualification.remove_rule")

type RemoveQualificationRuleCommand struct {
	EquipmentID     equipment.ID
	SpecificationID specifications.ID
}

func (c RemoveQualificationRuleCommand) Type() cqrs.CommandType {
	return CommandTypeRemoveQualificationRule
}

type RemoveQualificationRuleRepositoryPort interface {
	qualification.Reader
	qualification.Saver
}

func RemoveQualificationRule(repository RemoveQualificationRuleRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(RemoveQualificationRuleCommand)
		qualification, err := repository.Qualification(ctx, cmd.EquipmentID)
		if err != nil {
			return nil, err
		}
		if err := qualification.RemoveRule(cmd.SpecificationID); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, qualification); err != nil {
			return nil, err
		}
		return qualification.Events, nil
	}
}
