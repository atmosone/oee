package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/specifications"
)

const CommandTypeCreateSpecification = cqrs.CommandType("cmd.specifications.create")

type CreateSpecificationCommand struct {
	ID          specifications.ID
	Code        specifications.Code
	Description specifications.Description
}

func (c CreateSpecificationCommand) Type() cqrs.CommandType {
	return CommandTypeCreateSpecification
}

type CreateSpecificationRepositoryPort interface {
	specifications.Saver
}

func CreateSpecification(repository CreateSpecificationRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateSpecificationCommand)
		specification := specifications.New(
			cmd.ID,
			cmd.Code,
			cmd.Description,
		)
		if err := repository.Save(ctx, specification); err != nil {
			return nil, err
		}
		return specification.Events, nil
	}
}
