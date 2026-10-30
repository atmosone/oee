package queries

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/adapters/openapi/qualification"
)

const QueryTypeEquipmentSpecificationsList = cqrs.QueryType("equipment_specifications_list")

type EquipmentSpecificationsListQuery struct {
	qualification.EquipmentSpecificationsListParams
}

func (q EquipmentSpecificationsListQuery) Type() cqrs.QueryType {
	return QueryTypeEquipmentSpecificationsList
}

type EquipmentSpecificationsListStore interface {
	EquipmentSpecificationsList(
		ctx context.Context,
		params qualification.EquipmentSpecificationsListParams,
	) (qualification.EquipmentSpecificationsList, error)
}

func EquipmentSpecificationsList(
	store EquipmentSpecificationsListStore,
) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		query := q.(EquipmentSpecificationsListQuery)
		return store.EquipmentSpecificationsList(ctx, query.EquipmentSpecificationsListParams)
	}
}

const QueryTypeUnqualifiedDowntimesList = cqrs.QueryType("unqualified_downtimes_list")

type UnqualifiedDowntimesListQuery struct {
	qualification.UnqualifiedDowntimesListParams
}

func (q UnqualifiedDowntimesListQuery) Type() cqrs.QueryType {
	return QueryTypeUnqualifiedDowntimesList
}

type UnqualifiedDowntimesListStore interface {
	UnqualifiedDowntimesList(
		ctx context.Context,
		params qualification.UnqualifiedDowntimesListParams,
	) (qualification.UnqualifiedDowntimesList, error)
}

func UnqualifiedDowntimesList(
	store UnqualifiedDowntimesListStore,
) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		query := q.(UnqualifiedDowntimesListQuery)
		return store.UnqualifiedDowntimesList(ctx, query.UnqualifiedDowntimesListParams)
	}
}

const QueryTypeEquipment = cqrs.QueryType("equipment")

type EquipmentQuery struct {
	qualification.EquipmentParams
}

func (q EquipmentQuery) Type() cqrs.QueryType {
	return QueryTypeEquipment
}

type EquipmentStore interface {
	Equipment(
		ctx context.Context, params qualification.EquipmentParams,
	) (*qualification.Equipment, error)
}

func Equipment(store EquipmentStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		query := q.(EquipmentQuery)
		return store.Equipment(ctx, query.EquipmentParams)
	}
}

const QueryTypeCurrentEquipmentWorkshift = cqrs.QueryType("current_equipment_workshift")

type CurrentEquipmentWorkshiftQuery struct{}

func (q CurrentEquipmentWorkshiftQuery) Type() cqrs.QueryType {
	return QueryTypeCurrentEquipmentWorkshift
}

type CurrentEquipmentWorkshiftStore interface {
	CurrentEquipmentWorkshift(ctx context.Context) (*qualification.CurrentEquipmentWorkshift, error)
}

func CurrentEquipmentWorkshift(store CurrentEquipmentWorkshiftStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		return store.CurrentEquipmentWorkshift(ctx)
	}
}

const QueryTypeCurrentEquipmentDowntime = cqrs.QueryType("current_equipment_downtime")

type CurrentEquipmentDowntimeQuery struct {
	qualification.CurrentEquipmentDowntimeParams
}

func (q CurrentEquipmentDowntimeQuery) Type() cqrs.QueryType {
	return QueryTypeCurrentEquipmentDowntime
}

type CurrentEquipmentDowntimeStore interface {
	CurrentEquipmentDowntime(
		ctx context.Context, params qualification.CurrentEquipmentDowntimeParams,
	) (*qualification.CurrentEquipmentDowntime, error)
}

func CurrentEquipmentDowntime(store CurrentEquipmentDowntimeStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		query := q.(CurrentEquipmentDowntimeQuery)
		return store.CurrentEquipmentDowntime(ctx, query.CurrentEquipmentDowntimeParams)
	}
}
