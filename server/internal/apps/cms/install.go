package cms

import (
	"database/sql"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
	"github.com/atmosone/oee/internal/apps/cms/adapters/postgres"
	"github.com/atmosone/oee/internal/apps/cms/commands"
	"github.com/atmosone/oee/internal/apps/cms/queries"
)

func Install(app cqrs.App, db *sql.DB) *http.Server {
	app.OnCommand(commands.CommandTypeCreateEquipment,
		commands.CreateEquipment(
			postgres.NewCreateEquipmentRepository(db)))

	app.OnCommand(commands.CommandTypeCreateEquipmentLine,
		commands.CreateEquipmentLine(
			postgres.NewCreateEquipmentLineRepository(db)))

	app.OnCommand(commands.CommandTypeExpandEquipmentLine,
		commands.ExpandEquipmentLine(
			postgres.NewExpandEquipmentLineRepository(db)))

	app.OnCommand(commands.CommandTypeCreateSpecification,
		commands.CreateSpecification(
			postgres.NewCreateSpecificationRepository(db)))

	app.OnCommand(commands.CommandTypeCreateQualificationRule,
		commands.CreateQualificationRule(
			postgres.NewCreateQualificationRuleRepository(db)))

	app.OnQuery(queries.QueryTypeEquipmentList,
		queries.EquipmentList(
			postgres.NewEquipmentListStore(db)))

	app.OnQuery(queries.QueryTypeEquipmentLinesList,
		queries.EquipmentLinesList(
			postgres.NewEquipmentLinesListStore(db)))

	app.OnQuery(queries.QueryTypeSpecificationsList,
		queries.SpecificationsList(
			postgres.NewSpecificationsListStore(db)))

	app.OnQuery(queries.QueryTypeEquipmentDetails,
		queries.EquipmentDetails(
			postgres.NewEquipmentDetailsStore(db)))

	app.OnQuery(queries.QueryTypeEquipmentLineDetails,
		queries.EquipmentLineDetails(
			postgres.NewEquipmentLineDetailsStore(db)))

	app.OnQuery(queries.QueryTypeSpecificationDetails,
		queries.SpecificationDetails(
			postgres.NewSpecificationDetailsStore(db)))

	app.OnQuery(queries.QueryTypeAttachableSpecifications,
		queries.AttachableSpecifications(
			postgres.NewAttachableSpecificationsStore(db)))

	app.OnQuery(queries.QueryTypeExpandableEquipment,
		queries.ExpandableEquipment(
			postgres.NewExpandableEquipmentStore(db)))

	mux := http.NewServeMux()
	return &http.Server{
		Addr: ":8880",
		Handler: schema.HandlerWithOptions(
			openapi.NewHandler(app),
			schema.StdHTTPServerOptions{
				BaseRouter: mux,
				BaseURL:    "/api",
			},
		),
	}
}
