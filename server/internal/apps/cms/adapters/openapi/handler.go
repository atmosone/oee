package openapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
	"github.com/atmosone/oee/internal/apps/cms/commands"
	"github.com/atmosone/oee/internal/apps/cms/queries"
)

type App interface {
	cqrs.CommandHandler
	cqrs.QueryHandler
}

type Handler struct{ app App }

// AssetsCounter implements [schema.ServerInterface].
func (h *Handler) AssetsCounter(w http.ResponseWriter, r *http.Request) {
	handleQuery(r, w, h.app, queries.AssetsCounterQuery{})
}

// ExpandEquipmentLine implements [schema.ServerInterface].
func (h *Handler) ExpandEquipmentLine(w http.ResponseWriter, r *http.Request) {
	handleCommand[commands.ExpandEquipmentLineCommand](r, w, h.app)
}

// AttachableSpecifications implements [schema.ServerInterface].
func (h *Handler) AttachableSpecifications(w http.ResponseWriter, r *http.Request, params schema.AttachableSpecificationsParams) {
	handleQuery(r, w, h.app, queries.AttachableSpecificationsQuery{AttachableSpecificationsParams: params})
}

// EquipmentDetails implements [schema.ServerInterface].
func (h *Handler) EquipmentDetails(w http.ResponseWriter, r *http.Request, params schema.EquipmentDetailsParams) {
	handleQuery(r, w, h.app, queries.EquipmentDetailsQuery{EquipmentDetailsParams: params})
}

// EquipmentLineDetails implements [schema.ServerInterface].
func (h *Handler) EquipmentLineDetails(w http.ResponseWriter, r *http.Request, params schema.EquipmentLineDetailsParams) {
	handleQuery(r, w, h.app, queries.EquipmentLineDetailsQuery{EquipmentLineDetailsParams: params})
}

// EquipmentList implements [schema.ServerInterface].
func (h *Handler) EquipmentList(w http.ResponseWriter, r *http.Request) {
	handleQuery(r, w, h.app, queries.EquipmentListQuery{})
}

// EquipmentLinesList implements [schema.ServerInterface].
func (h *Handler) EquipmentLinesList(w http.ResponseWriter, r *http.Request) {
	handleQuery(r, w, h.app, queries.EquipmentLinesListQuery{})
}

// ExpandableEquipment implements [schema.ServerInterface].
func (h *Handler) ExpandableEquipment(w http.ResponseWriter, r *http.Request, params schema.ExpandableEquipmentParams) {
	handleQuery(r, w, h.app, queries.ExpandableEquipmentQuery{ExpandableEquipmentParams: params})
}

// SpecificationDetails implements [schema.ServerInterface].
func (h *Handler) SpecificationDetails(w http.ResponseWriter, r *http.Request, params schema.SpecificationDetailsParams) {
	handleQuery(r, w, h.app, queries.SpecificationDetailsQuery{SpecificationDetailsParams: params})
}

// SpecificationsList implements [schema.ServerInterface].
func (h *Handler) SpecificationsList(w http.ResponseWriter, r *http.Request) {
	handleQuery(r, w, h.app, queries.SpecificationsListQuery{})
}

// CreateEquipment implements [schema.ServerInterface].
func (h *Handler) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	handleCommand[commands.CreateEquipmentCommand](r, w, h.app)
}

// CreateEquipmentLine implements [schema.ServerInterface].
func (h *Handler) CreateEquipmentLine(w http.ResponseWriter, r *http.Request) {
	handleCommand[commands.CreateEquipmentLineCommand](r, w, h.app)
}

// CreateQualificationRule implements [schema.ServerInterface].
func (h *Handler) CreateQualificationRule(w http.ResponseWriter, r *http.Request) {
	handleCommand[commands.CreateQualificationRuleCommand](r, w, h.app)
}

// CreateSpecification implements [schema.ServerInterface].
func (h *Handler) CreateSpecification(w http.ResponseWriter, r *http.Request) {
	handleCommand[commands.CreateSpecificationCommand](r, w, h.app)
}

// Events implements [schema.ServerInterface].
func (h *Handler) Events(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func handleQuery[Q cqrs.Query](r *http.Request, w http.ResponseWriter, app cqrs.QueryHandler, q Q) {
	res, err := app.HandleQuery(r.Context(), q)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func handleCommand[C cqrs.Command](r *http.Request, w http.ResponseWriter, app cqrs.CommandHandler) {
	var cmd C
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := app.HandleCommand(r.Context(), cmd); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NewHandler(app App) *Handler {
	return &Handler{app: app}
}
