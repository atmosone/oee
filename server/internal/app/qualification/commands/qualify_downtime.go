package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/downtimes"
	"github.com/atmosone/oee/internal/domain/qualification"
	"github.com/atmosone/oee/internal/domain/specifications"
)

const CommandTypeQualifyDowntime = cqrs.CommandType("cmd.qualification.qualify_downtime")

type QualifyDowntimeCommand struct {
	DowntimeID      downtimes.ID
	SpecificationID specifications.ID
	Comment         qualification.Comment
}

func (c QualifyDowntimeCommand) Type() cqrs.CommandType {
	return CommandTypeQualifyDowntime
}

type QualifyDowntimeRepositoryPort interface {
	qualification.DowntimeReader
	qualification.Reader
	qualification.Saver
}

func QualifyDowntime(repository QualifyDowntimeRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(QualifyDowntimeCommand)
		downtime, err := repository.Downtime(ctx, cmd.DowntimeID)
		if err != nil {
			return nil, err
		}
		qualification, err := repository.Qualification(ctx, downtime.EquipmentID())
		if err != nil {
			return nil, err
		}
		if err := qualification.QualifyDowntime(downtime,
			cmd.SpecificationID,
			cmd.Comment,
		); err != nil {
			return nil, err
		}
		defer qualification.ClearEvents()
		if err := repository.Save(ctx, qualification); err != nil {
			return nil, err
		}
		return qualification.Events(), nil
	}
}
