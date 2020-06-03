package calculate

import (
	"django-go/pkg/django"
	"django-go/pkg/types"
	"django-go/pkg/util"
	"fmt"
	"strconv"
)

type CalculateReschedule struct {
	start int64
}

func NewReschedule(start int64) django.RescheduleInterface {
	return &CalculateReschedule{start}
}

func (reschedule *CalculateReschedule) Reschedule(nodeWithPods []types.NodeWithPod, rule types.Rule) ([]types.RescheduleResult, error) {

	nodeWithPods4CheckAgainst := types.CopyNodeWithPods(nodeWithPods)

	againstPods := searchAgainstPods(nodeWithPods4CheckAgainst, rule)

	return reschedule.calculate(nodeWithPods, againstPods, rule), nil
}

func (reschedule *CalculateReschedule) calculate(nodeWithPods []types.NodeWithPod, againstPods map[string][]types.Pod, rule types.Rule) []types.RescheduleResult {

	groupRuleAssociates := types.FromPods(util.ToPods(nodeWithPods))

	allMaxInstancePerNodeLimit := util.ToAllMaxInstancePerNodeLimit(rule, groupRuleAssociates)

	rescheduleResults := make([]types.RescheduleResult, 0)

	forsakePods := make([]types.Pod, 0)

	for sourceSn, pods := range againstPods {

		for _, pod := range pods {

			forsake := true

			for _, nwp := range nodeWithPods {

				if util.RuleOverrunTimeLimit(rule, reschedule.start) {
					fmt.Println("overrun time limit")
					return rescheduleResults
				}

				if sourceSn == nwp.Node.Sn {
					continue
				}

				if _, ok := againstPods[nwp.Node.Sn]; ok {
					continue
				}

				if util.StaticFillOnePod(&nwp, pod, allMaxInstancePerNodeLimit) {
					forsake = false

					rescheduleResults = append(rescheduleResults, types.RescheduleResult{
						Stage:    1,
						SourceSn: sourceSn,
						TargetSn: nwp.Node.Sn,
						PodSn:    pod.PodSn,
						CpuIds:   pod.CpuIds,
					})
					break
				}

			}

			if forsake {
				forsakePods = append(forsakePods, pod)
			}

		}

	}

	if len(forsakePods) > 0 {
		fmt.Println("forsake pod count: " + strconv.Itoa(len(forsakePods)))
	}

	return rescheduleResults
}

func searchAgainstPods(nodeWithPods []types.NodeWithPod, rule types.Rule) map[string][]types.Pod {

	result := make(map[string][]types.Pod, 0)

	for sn, pods := range searchResourceAgainstPods(nodeWithPods) {
		if old, ok := result[sn]; !ok {
			result[sn] = pods
		} else {
			result[sn] = append(old, pods...)
		}
	}

	for sn, pods := range searchLayoutAgainstPods(nodeWithPods, rule) {
		if old, ok := result[sn]; !ok {
			result[sn] = pods
		} else {
			result[sn] = append(old, pods...)
		}
	}

	for sn, pods := range searchCgroupAgainstPods(nodeWithPods) {
		if old, ok := result[sn]; !ok {
			result[sn] = pods
		} else {
			result[sn] = append(old, pods...)
		}
	}

	return result

}

func searchResourceAgainstPods(nodeWithPods []types.NodeWithPod) map[string][]types.Pod {

	result := make(map[string][]types.Pod, 0)

	for _, nwp := range nodeWithPods {

		againstPods := make([]types.Pod, 0)

		tempPods := make([]types.Pod, 0)

		normalPods := make([]types.Pod, 0)

		for _, pod := range nwp.Pods {

			against := false

			for _, resource := range types.AllResources {

				nodeResource := nwp.Node.Value(resource)

				supposePods := make([]types.Pod, len(tempPods)+1)

				supposePods = append(supposePods, tempPods...)

				supposePods = append(supposePods, pod)

				podsResource := util.PodsTotalResource(supposePods, resource)

				if nodeResource < podsResource {
					againstPods = append(againstPods, pod)
					against = true
					break
				}

			}

			if !against {
				tempPods = append(tempPods, pod)
			}

		}

		eniAgainstPodSize := len(tempPods) - nwp.Node.Eni

		if eniAgainstPodSize > 0 {

			againstPods = append(againstPods, tempPods[0:eniAgainstPodSize]...)

			normalPods = append(normalPods, tempPods[eniAgainstPodSize:len(tempPods)-1]...)

		} else {
			normalPods = append(normalPods, tempPods...)
		}

		nwp.Pods = normalPods

		if len(againstPods) > 0 {
			result[nwp.Node.Sn] = againstPods
		}

	}

	return result
}

func searchLayoutAgainstPods(nodeWithPods []types.NodeWithPod, rule types.Rule) map[string][]types.Pod {

	result := make(map[string][]types.Pod, 0)

	groupRuleAssociates := types.FromPods(util.ToPods(nodeWithPods))

	maxInstancePerNodes := util.ToAllMaxInstancePerNodeLimit(rule, groupRuleAssociates)

	for _, nwp := range nodeWithPods {

		groupCountPreNodeMap := make(map[string]int)

		againstPods := make([]types.Pod, 0)

		normalPods := make([]types.Pod, 0)

		for _, pod := range nwp.Pods {

			maxInstancePerNode := maxInstancePerNodes[pod.Group]

			oldValue := 0

			if value, ok := groupCountPreNodeMap[pod.Group]; ok {
				oldValue = value
			}

			if oldValue == maxInstancePerNode {
				againstPods = append(againstPods, pod)
				continue
			}

			groupCountPreNodeMap[pod.Group] = oldValue + 1

			normalPods = append(normalPods, pod)

		}

		nwp.Pods = normalPods

		if len(againstPods) > 0 {
			result[nwp.Node.Sn] = againstPods
		}

	}

	return result

}

func searchCgroupAgainstPods(nodeWithPods []types.NodeWithPod) map[string][]types.Pod {

	result := make(map[string][]types.Pod, 0)

	for _, nwp := range nodeWithPods {

		node := nwp.Node

		if len(node.Topologies) == 0 {
			continue
		}

		if len(nwp.Pods) == 0 {
			continue
		}

		againstPods := make([]types.Pod, 0)

		againstCpuIds := make(map[int]bool)

		for key, value := range util.CpuIDCountMap(nwp) {
			if value > 1 {
				againstCpuIds[key] = true
			}
		}

		cpuToSocket := util.CpuToSocket(node)

		cpuToCore := util.CpuToCore(node)

		normalPods := make([]types.Pod, 0)

		for _, pod := range nwp.Pods {

			if len(pod.CpuIds) == 0 {
				againstPods = append(againstPods, pod)
				continue
			}

			tempCpuIds := make(map[int]bool)

		outter:
			for cpuId := range againstCpuIds {

				for _, podCpuId := range pod.CpuIds {

					if podCpuId == cpuId {
						continue outter
					}
				}

				tempCpuIds[cpuId] = true

			}

			if len(againstCpuIds) != len(tempCpuIds) {

				againstCpuIds = tempCpuIds

				againstPods = append(againstPods, pod)

				continue
			}

			socketCountMap := make(map[int]bool)

			for _, cpuId := range pod.CpuIds {
				socketCountMap[cpuToSocket[cpuId]] = true
			}

			socketCount := len(socketCountMap)

			if socketCount > 1 {

				againstPods = append(againstPods, pod)

				continue
			}

			coreCountMap := make(map[int]int)

			for _, cpuId := range pod.CpuIds {

				if core, ok := cpuToCore[cpuId]; !ok {
					coreCountMap[core] = 1
				} else {
					coreCountMap[core] = coreCountMap[core] + 1
				}

			}

			sameCoreMap := make(map[int]int)

			for key, value := range coreCountMap {
				if value > 1 {
					sameCoreMap[key] = value
				}
			}

			if len(sameCoreMap) > 0 {

				againstPods = append(againstPods, pod)

				continue
			}

			//TODO 已经很复杂了，如果排名拉不开差距在增加sensitiveCpuBind数据校验

			normalPods = append(normalPods, pod)

		}

		nwp.Pods = normalPods

		if len(againstPods) > 0 {
			result[nwp.Node.Sn] = againstPods
		}

	}

	return result

}
