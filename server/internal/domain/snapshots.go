package domain

import (
	"time"
)

type WorkshiftSnapshot struct {
	ID         string
	Name       string
	Status     string
	Duration   int64
	StartedAt  time.Time
	FinishedAt time.Time
}

func HydrateWorkshift(snapshot WorkshiftSnapshot) *Workshift {
	return &Workshift{
		reactiveEntity: &reactiveEntity{},
		id:             WorkshiftID(snapshot.ID),
		name:           WorkshiftName(snapshot.Name),
		status:         WorkshiftStatus(snapshot.Status),
		duration:       time.Duration(snapshot.Duration),
		startedAt:      snapshot.StartedAt,
		finishedAt:     snapshot.FinishedAt,
	}
}

type EquipmentLineSnapshot struct {
	ID        string
	Name      string
	Equipment []string
}

func HydrateEquipmentLine(snapshot EquipmentLineSnapshot) *EquipmentLine {
	el := &EquipmentLine{
		reactiveEntity: &reactiveEntity{},
		id:             EquipmentLineID(snapshot.ID),
		name:           EquipmentLineName(snapshot.Name),
		equipment:      make(map[EquipmentID]struct{}),
	}
	for _, item := range snapshot.Equipment {
		el.equipment[EquipmentID(item)] = struct{}{}
	}
	return el
}

type DowntimeSnapshot struct {
	ID          string
	EquipmentID string
	Status      string
	Duration    int64
	StartedAt   time.Time
	FinishedAt  time.Time
}

func HydrateDowntime(snapshot DowntimeSnapshot) *Downtime {
	return &Downtime{
		reactiveEntity: &reactiveEntity{},
		id:             DowntimeID(snapshot.ID),
		equipmentID:    EquipmentID(snapshot.EquipmentID),
		status:         DowntimeStatus(snapshot.Status),
		duration:       time.Duration(snapshot.Duration),
		startedAt:      snapshot.StartedAt,
		finishedAt:     snapshot.FinishedAt,
	}
}
