package queries

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

const QueryTypeEquipmentLinesList = cqrs.QueryType("equipment_lines_list")

type EquipmentLinesListQuery struct{}

func (q EquipmentLinesListQuery) Type() cqrs.QueryType { return QueryTypeEquipmentLinesList }

type EquipmentLinesListStore interface {
	EquipmentLinesList(ctx context.Context) (schema.EquipmentLinesList, error)
}

func EquipmentLinesList(store EquipmentLinesListStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		return store.EquipmentLinesList(ctx)
	}
}
