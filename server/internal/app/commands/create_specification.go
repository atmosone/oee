package commands

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain"
)

const CommandTypeCreateSpecification = cqrs.CommandType("create_specification")

type CreateSpecificationCommand struct {
	SpecificationID          string `json:"specification_id"`
	SpecificationCode        string `json:"specification_code"`
	SpecificationDescription string `json:"specification_description"`
}

func (c CreateSpecificationCommand) Validate() error {
	return errors.Join(
		domain.SpecificationID(c.SpecificationID).Validate(),
		domain.SpecificationCode(c.SpecificationCode).Validate(),
		domain.SpecificationDescription(c.SpecificationDescription).Validate(),
	)
}

func (c CreateSpecificationCommand) Type() cqrs.CommandType {
	return CommandTypeCreateSpecification
}

type CreateSpecificationRepository interface {
	domain.SpecificationSaver
}

func CreateSpecification(repository CreateSpecificationRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateSpecificationCommand)
		specification := domain.NewSpecification(
			domain.SpecificationID(cmd.SpecificationID),
			domain.SpecificationCode(cmd.SpecificationCode),
			domain.SpecificationDescription(cmd.SpecificationDescription),
		)
		defer specification.Clear()
		if err := repository.SaveSpecification(ctx, specification); err != nil {
			return nil, err
		}
		return specification.Events(), nil
	}
}
