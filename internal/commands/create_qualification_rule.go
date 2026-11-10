package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/qualification"
	"github.com/atmosone/oee/internal/domain/specifications"
)

const CommandTypeCreateQualificationRule = cqrs.CommandType("cmd.qualification.create_rule")

type CreateQualificationRuleCommand struct {
	EquipmentID     equipment.ID
	SpecificationID specifications.ID
}

func (c CreateQualificationRuleCommand) Type() cqrs.CommandType {
	return CommandTypeCreateQualificationRule
}

type CreateQualificationRuleRepositoryPort interface {
	qualification.SpecificationReader
	qualification.Reader
	qualification.Saver
}

func CreateQualificationRule(repository CreateQualificationRuleRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateQualificationRuleCommand)
		specification, err := repository.Specification(ctx, cmd.SpecificationID)
		if err != nil {
			return nil, err
		}
		qualification, err := repository.Qualification(ctx, cmd.EquipmentID)
		if err != nil {
			return nil, err
		}
		if err := qualification.CreateRule(specification); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, qualification); err != nil {
			return nil, err
		}
		return qualification.Events, nil
	}
}
