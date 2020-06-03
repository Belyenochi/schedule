package calculate

import (
	"fmt"
	"sort"
	"strconv"
	"django-go/pkg/django"
	"django-go/pkg/types"
	"django-go/pkg/util"
)

type CalculateSchedule struct {
	start int64
}

func NewSchedule(start int64) django.ScheduleInterface {
	return &CalculateSchedule{start}
}

func (schedule *CalculateSchedule) Schedule(nodes []types.Node, apps []types.App, rule types.Rule) ([]types.ScheduleResult, error) {

	nodeWithPods := sortAndInitNodeWithPods(nodes, rule)

	allPods := sortAndInitPods(apps)

	allMaxInstancePerNodeLimit := util.ToAllMaxInstancePerNodeLimit(rule, types.FromApps(apps))

	schedule.calculate(nodeWithPods, allPods, rule, allMaxInstancePerNodeLimit)

	results := make([]types.ScheduleResult, 0)
	for _, nwp := range nodeWithPods {
		for _, pod := range nwp.Pods {
			result := types.ScheduleResult{
				Sn:     nwp.Node.Sn,
				Group:  pod.Group,
				CpuIds: pod.CpuIds,
			}
			results = append(results, result)
		}
	}

	return results, nil
}

func sortAndInitNodeWithPods(nodes []types.Node, rule types.Rule) []types.NodeWithPod {

	sort.SliceStable(nodes, func(i, j int) bool {
		return util.ResourceNodeScore(nodes[i], rule) > util.ResourceNodeScore(nodes[j], rule)
	})

	nodeWithPods := make([]types.NodeWithPod, 0, len(nodes))

	for _, node := range nodes {
		nodeWithPods = append(nodeWithPods, types.NodeWithPod{
			Node: node,
			Pods: make([]types.Pod, 0),
		})
	}

	return nodeWithPods
}

func sortAndInitPods(apps []types.App) []types.Pod {

	sort.SliceStable(apps, func(i, j int) bool {
		//对比pod数量
		if apps[i].Replicas == apps[j].Replicas {

			//如果GPU相对，则对比CPU
			if apps[i].Gpu == apps[j].Gpu {
				//如果CPU相对，则对比内存
				if apps[i].Cpu == apps[j].Cpu {
					//如果内存相关，则对比disk
					if apps[i].Ram == apps[j].Ram {
						return apps[i].Disk > apps[j].Disk
					}
					return apps[i].Ram > apps[j].Ram
				}
				return apps[i].Cpu > apps[j].Cpu
			}
			return apps[i].Gpu > apps[j].Gpu

		}

		return apps[i].Replicas > apps[j].Replicas
	})

	pods := make([]types.Pod, 0)

	for _, app := range apps {
		for i := 0; i < app.Replicas; i++ {
			pods = append(pods, util.ToPod(app))
		}
	}

	fmt.Println("schedule app transform pod count: " + strconv.Itoa(len(pods)))

	return pods
}

func (schedule *CalculateSchedule) calculate(nodeWithPods []types.NodeWithPod, pods []types.Pod, rule types.Rule, allMaxInstancePerNodeLimit map[string]int) {

	forsakePods := make([]types.Pod, 0)

	for _, pod := range pods {

		forsake := true

		for i := range nodeWithPods {

			if util.RuleOverrunTimeLimit(rule, schedule.start) {
				fmt.Println("overrun time limit")
				return
			}

			if util.StaticFillOnePod(&nodeWithPods[i], pod, allMaxInstancePerNodeLimit) {
				forsake = false
				break
			}
		}

		if forsake {
			forsakePods = append(forsakePods, pod)
		}
	}

	if len(forsakePods) > 0 {
		fmt.Println("forsake pod count: " + strconv.Itoa(len(forsakePods)))
	}

}
