package views

import (
	"time"

	"github.com/atmosone/oee/pkg/openapi"
)

type TimelineWorkshift struct {
	ID string `json:"id"`
	Interval
}

type TimelineRecordCategory string

const TimelineRecordCategoryWorktime = TimelineRecordCategory("worktime")
const TimelineRecordCategoryDowntime = TimelineRecordCategory("downtime")

type TimelineRecord struct {
	ID       string                 `json:"id"`
	Category TimelineRecordCategory `json:"category"`
	Interval
}

type TimelineEquipment struct {
	ID   string `json:"equipment_id"`
	Name string `json:"name"`
}

type Timeline struct {
	Workshift TimelineWorkshift `json:"workshift"`
	Equipment TimelineEquipment `json:"equipment"`
	Records   []TimelineRecord  `json:"records"`
}

type ApexchartsTimeline struct {
	Name string                  `json:"name"`
	Data []ApexchartsTimelineBar `json:"data"`
}

type ApexchartsTimelineBar struct {
	X string   `json:"x"`
	Y [2]int64 `json:"y"`
}

func (t Timeline) ApexchartsTimeline(now time.Time) []openapi.ApexchartsTimeline {
	downtimes := openapi.ApexchartsTimeline{
		Name: openapi.IntervalCategory(TimelineRecordCategoryDowntime),
		Data: make([]openapi.ApexchartsTimelineBar, 0),
	}
	worktimes := openapi.ApexchartsTimeline{
		Name: openapi.IntervalCategory(TimelineRecordCategoryWorktime),
		Data: make([]openapi.ApexchartsTimelineBar, 0),
	}
	for _, record := range t.Records {
		intersection, ok := t.Workshift.Intersection(now, record.Interval)
		if !ok {
			continue
		}
		switch record.Category {
		case TimelineRecordCategoryDowntime:
			downtimes.Data = append(downtimes.Data, openapi.ApexchartsTimelineBar{
				Id: openapi.IntervalID(record.ID),
				X:  openapi.EquipmentName(t.Equipment.Name),
				Y: []int{
					int(intersection.Start().UnixMilli()),
					int(intersection.Finish(now).UnixMilli()),
				},
			})
		case TimelineRecordCategoryWorktime:
			worktimes.Data = append(worktimes.Data, openapi.ApexchartsTimelineBar{
				Id: openapi.IntervalID(record.ID),
				X:  openapi.EquipmentName(t.Equipment.Name),
				Y: []int{
					int(intersection.Start().UnixMilli()),
					int(intersection.Finish(now).UnixMilli()),
				},
			})
		}
	}
	return []openapi.ApexchartsTimeline{downtimes, worktimes}
}
