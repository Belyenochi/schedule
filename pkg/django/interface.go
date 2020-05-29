package django

import "django-go/pkg/types"

type ScheduleResult struct {
	Sn     string
	Group  string
	CpuIds []int
}

type RescheduleResult struct {
	Stage    int
	SourceSn string
	TargetSn string
	Group    string
	CpuIds   []int
}

type ScheduleInterface interface {
	Schedule(nodes []types.Node, apps []types.App, rule types.Rule) ([]ScheduleResult, error)
	Reschedule(nodeWithPods []types.NodeWithPod, rule types.Rule) ([]RescheduleResult, error)
}
