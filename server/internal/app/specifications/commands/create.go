package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/specifications"
)

const CommandTypeCreate = cqrs.CommandType("cmd.specifications.create")

type CreateCommand struct {
	ID          specifications.ID
	Code        specifications.Code
	Description specifications.Description
}

func (c CreateCommand) Type() cqrs.CommandType {
	return CommandTypeCreate
}

type CreateRepositoryPort interface {
	specifications.Saver
}

func Create(repository CreateRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateCommand)
		specification := specifications.New(
			cmd.ID,
			cmd.Code,
			cmd.Description,
		)
		defer specification.ClearEvents()
		if err := repository.Save(ctx, specification); err != nil {
			return nil, err
		}
		return specification.Events(), nil
	}
}
