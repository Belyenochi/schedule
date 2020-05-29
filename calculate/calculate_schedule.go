package calculate

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"django-go/pkg/django"
	"django-go/pkg/types"
	"django-go/pkg/util"
)

type CalculateSchedule struct {
	start int64
}

func NewSchedule() django.ScheduleInterface {
	return &CalculateSchedule{}
}

func (schedule *CalculateSchedule) Schedule(nodes []types.Node, apps []types.App, rule types.Rule) ([]django.ScheduleResult, error) {
	schedule.start = time.Now().Unix()

	nodeWithPods := buildNodeWithPods(nodes, rule)

	allPods := buildPods(apps)

	allMaxInstancePerHostLimit := buildAllMaxInstancePerHostLimit(rule, apps)

	schedule.calculate(nodeWithPods, allPods, rule, allMaxInstancePerHostLimit)

	return buildScheduleResult(nodeWithPods), nil
}

func buildScheduleResult(nwps []types.NodeWithPod) []django.ScheduleResult {
	results := make([]django.ScheduleResult, 0)
	for _, nwp := range nwps {
		for _, pod := range nwp.Pods {
			result := django.ScheduleResult{
				Sn:     nwp.Node.Sn,
				Group:  pod.Group,
				CpuIds: pod.CpuIds,
			}
			results = append(results, result)
		}
	}

	return results
}

func (schedule *CalculateSchedule) Reschedule(nodeWithPods []types.NodeWithPod, rule types.Rule) ([]django.RescheduleResult, error) {
	return nil, nil
}

func (schedule *CalculateSchedule) ruleTimeout(rule types.Rule) bool {
	return time.Now().Unix()-schedule.start > int64(rule.TimeLimitInMins*60*1000)
}

func buildNodeWithPods(nodes []types.Node, rule types.Rule) []types.NodeWithPod {
	sort.Slice(nodes, func(i, j int) bool {
		return util.ScoreNode(nodes[i], rule) < util.ScoreNode(nodes[j], rule)
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

func buildPods(apps []types.App) []types.Pod {
	sort.Slice(apps, func(i, j int) bool {
		//对比pod数量
		if apps[i].Replicas == apps[j].Replicas {
			//如果CPU相对，则对比内存
			if apps[i].Cpu == apps[j].Cpu {
				//如果内存相关，则对比disk
				if apps[i].Ram == apps[j].Ram {
					return apps[i].Disk < apps[j].Disk
				}
				return apps[i].Ram < apps[j].Ram
			}
			return apps[i].Cpu < apps[j].Cpu
		}
		return apps[i].Replicas < apps[j].Replicas
	})

	pods := make([]types.Pod, 0)
	for _, app := range apps {
		for i := 0; i < app.Replicas; i++ {
			pods = append(pods, util.BuildPodForApp(app))
		}
	}

	return pods
}

func buildAllMaxInstancePerHostLimit(rule types.Rule, apps []types.App) map[string]int {
	allMaxInstancePerHostLimit := make(map[string]int)

	for _, replicasMaxInstancePerHost := range rule.ReplicasMaxInstancePerHosts {
		restrain := replicasMaxInstancePerHost.Restrain
		if !util.StringSlice(types.AllRestrains).Contain(string(restrain)) {
			continue
		}

		replicas := replicasMaxInstancePerHost.Replicas
		maxInstancePerHost := replicasMaxInstancePerHost.MaxInstancePerHost

		for _, app := range apps {
			if (restrain == types.GE && app.Replicas >= replicas) ||
				(restrain == types.LE && app.Replicas <= replicas) {
				allMaxInstancePerHostLimit[app.Group] = maxInstancePerHost
			}
		}
	}

	for _, app := range apps {
		if _, ok := allMaxInstancePerHostLimit[app.Group]; !ok {
			allMaxInstancePerHostLimit[app.Group] = rule.DefaultMaxInstancePerHost
		}
	}

	return allMaxInstancePerHostLimit
}

func (schedule *CalculateSchedule) calculate(nodeWithPods []types.NodeWithPod, pods []types.Pod, rule types.Rule, podPerHostLimit map[string]int) {
	forsakePods := make([]types.Pod, 0)

	for _, pod := range pods {
		forsake := true

		for i := range nodeWithPods {
			if schedule.ruleTimeout(rule) {
				fmt.Println("rule timeout")
				return
			}

			if nodeFillOnePod(nodeWithPods[i], pod, podPerHostLimit) {
				forsake = false
				nodeWithPods[i].Pods = append(nodeWithPods[i].Pods, pod)
				break
			}
		}

		if forsake {
			forsakePods = append(forsakePods, pod)
		}
	}

	fmt.Println("forsake pod count: " + strconv.Itoa(len(forsakePods)))
}

func violateMaxInstancePerHost(nwp types.NodeWithPod, pod types.Pod, podPerHostLimit map[string]int) bool {
	groupPodCount := 0
	for _, p := range nwp.Pods {
		if pod.Group == p.Group {
			groupPodCount++
		}
	}

	return groupPodCount > podPerHostLimit[pod.Group]
}

func nodeFillOnePod(nwp types.NodeWithPod, pod types.Pod, podPerHostLimit map[string]int) bool {
	return nwp.Node.Gpu >= util.PodsTotalResource(nwp.Pods, types.GPU) &&
		nwp.Node.Cpu >= util.PodsTotalResource(nwp.Pods, types.CPU) &&
		nwp.Node.Disk >= util.PodsTotalResource(nwp.Pods, types.Disk) &&
		nwp.Node.Ram >= util.PodsTotalResource(nwp.Pods, types.RAM) &&
		nwp.Node.Eni >= len(nwp.Pods)+1 &&
		!violateMaxInstancePerHost(nwp, pod, podPerHostLimit)
}
