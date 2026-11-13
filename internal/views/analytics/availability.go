package analytics

import "time"

type DowntimeInterval struct {
	Interval
}

func ComputeWorkshiftEquipmentAvailabilityRate(
	now time.Time,
	workshift Workshift,
	downtimes []DowntimeInterval,
) float64 {
	workshiftDuration := workshift.Duration(now).Seconds()
	if workshiftDuration == 0 {
		return 0
	}
	var downtimeDuration float64
	for _, downtime := range downtimes {
		downtimeDuration += workshift.Allocation(now, downtime.Interval).Seconds()
	}
	if downtimeDuration == 0 {
		return 1
	}
	return (workshiftDuration - downtimeDuration) / workshiftDuration
}

func ComputeWorkshiftEquipmentQualityRate(production Production) float64 {
	if production.Good == 0 && production.Defect == 0 {
		return 1
	}
	return float64(production.Good) / float64(production.Good+production.Defect)
}

func ComputeWorkshiftEquipmentPerformanceRate() float64 { return 1 }

func ComputeWorkshiftEquipmentOEE(availability, performance, quality float64) float64 {
	return availability * performance * quality
}
