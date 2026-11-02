package queries

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

const QueryTypeEquipmentList = cqrs.QueryType("equipment_list")

type EquipmentListQuery struct{}

func (q EquipmentListQuery) Type() cqrs.QueryType { return QueryTypeEquipmentList }

type EquipmentListStore interface {
	EquipmentList(ctx context.Context) (schema.EquipmentList, error)
}

func EquipmentList(store EquipmentListStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		return store.EquipmentList(ctx)
	}
}
