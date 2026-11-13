package queries

import (
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/pkg/openapi"
)

var ErrNotFound = errors.New("queries: not found")

const (
	QueryTypeLinesList            = cqrs.QueryType("lines_list")
	QueryTypeLineTopology         = cqrs.QueryType("line_topology")
	QueryTypeCurrentWorkshift     = cqrs.QueryType("current_workshift")
	QueryTypeCurrentInterval      = cqrs.QueryType("current_interval")
	QueryTypeCurrentIntervalsList = cqrs.QueryType("current_intervals_list")
)

type LinesListQuery struct{}

func (q LinesListQuery) Type() cqrs.QueryType {
	return QueryTypeLinesList
}

type LineTopologyQuery struct {
	openapi.LineTopologyParams
}

func (q LineTopologyQuery) Type() cqrs.QueryType {
	return QueryTypeLineTopology
}

type CurrentWorkshiftQuery struct{}

func (q CurrentWorkshiftQuery) Type() cqrs.QueryType {
	return QueryTypeCurrentWorkshift
}

type CurrentIntervalQuery struct {
	openapi.CurrentIntervalParams
}

func (q CurrentIntervalQuery) Type() cqrs.QueryType {
	return QueryTypeCurrentInterval
}

type CurrentIntervalsListQuery struct{}

func (q CurrentIntervalsListQuery) Type() cqrs.QueryType {
	return QueryTypeCurrentIntervalsList
}
