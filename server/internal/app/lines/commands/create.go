package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/lines"
)

const CommandTypeCreate = cqrs.CommandType("cmd.lines.create")

type CreateCommand struct {
	ID   lines.ID
	Name lines.Name
}

func (c CreateCommand) Type() cqrs.CommandType {
	return CommandTypeCreate
}

type CreateRepositoryPort interface {
	lines.Saver
}

func Create(repository CreateRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateCommand)
		line := lines.New(cmd.ID, cmd.Name)
		defer line.ClearEvents()
		if err := repository.Save(ctx, line); err != nil {
			return nil, err
		}
		return line.Events(), nil
	}
}
