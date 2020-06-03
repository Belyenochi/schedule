package django

import "django-go/pkg/types"

type RescheduleInterface interface {
	Reschedule(nodeWithPods []types.NodeWithPod, rule types.Rule) ([]types.RescheduleResult, error)
}