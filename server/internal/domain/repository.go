package domain

import "context"

type WorkshiftProvider interface {
	Workshift(context.Context, WorkshiftID) (*Workshift, error)
}

type WorkshiftSaver interface {
	SaveWorkshift(context.Context, *Workshift) error
}

type WorkshiftRepository interface {
	WorkshiftProvider
	WorkshiftSaver
}

type EquipmentSaver interface {
	SaveEquipment(context.Context, *Equipment) error
}

type EquipmentRepository interface {
	EquipmentSaver
}

type EquipmentLineProvider interface {
	EquipmentLine(context.Context, EquipmentLineID) (*EquipmentLine, error)
}

type EquipmentLineSaver interface {
	SaveEquipmentLine(context.Context, *EquipmentLine) error
}

type EquipmentLineRepository interface {
	EquipmentLineProvider
	EquipmentLineSaver
}

type DowntimeSaver interface {
	SaveDowntime(context.Context, *Downtime) error
}

type DowntimeProvider interface {
	Downtime(context.Context, DowntimeID) (*Downtime, error)
}

type DowntimeRepository interface {
	DowntimeProvider
	DowntimeSaver
}

type SpecificationSaver interface {
	SaveSpecification(context.Context, *Specification) error
}

type SpecificationRepository interface {
	SpecificationSaver
}

type QualificationRuleSaver interface {
	SaveQualificationRule(context.Context, *QualificationRule) error
}

type QualificationRuleRepository interface {
	QualificationRuleSaver
}

type QualificationSaver interface {
	SaveQualification(context.Context, *Qualification) error
}

type QualificationRepository interface {
	QualificationSaver
}
