package queries

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

const QueryTypeExpandableEquipment = cqrs.QueryType("expandable_equipment")

type ExpandableEquipmentQuery struct {
	schema.ExpandableEquipmentParams
}

func (q ExpandableEquipmentQuery) Type() cqrs.QueryType {
	return QueryTypeExpandableEquipment
}

type ExpandableEquipmentStore interface {
	ExpandableEquipment(ctx context.Context, params schema.ExpandableEquipmentParams) (schema.ExpandableEquipmentList, error)
}

func ExpandableEquipment(store ExpandableEquipmentStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		return store.ExpandableEquipment(ctx, q.(ExpandableEquipmentQuery).ExpandableEquipmentParams)
	}
}
