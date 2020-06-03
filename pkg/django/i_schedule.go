package django

import "django-go/pkg/types"

type ScheduleInterface interface {
	Schedule(nodes []types.Node, apps []types.App, rule types.Rule) ([]types.ScheduleResult, error)
}