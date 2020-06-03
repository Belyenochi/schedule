package util

import (
	"django-go/pkg/types"
	"time"
	"django-go/pkg/constants"
)

func StaticFillOnePod(nwp *types.NodeWithPod, pod types.Pod, allMaxInstancePerNodeLimit map[string]int) bool {

	if !ResourceFillOnePod(*nwp, pod) {
		return false
	}

	if !LayoutFillOnePod(allMaxInstancePerNodeLimit, *nwp, pod) {
		return false
	}

	podPreAlloc := CgroupFillOnePod(*nwp, pod)

	if !podPreAlloc.Satisfy {
		return false
	}

	if len(podPreAlloc.Cpus) > 0 {
		pod.CpuIds = podPreAlloc.Cpus
	}

	nwp.Pods = append(nwp.Pods, pod)

	return true
}

func ResourceFillOnePod(nwp types.NodeWithPod, pod types.Pod) bool {

	node := nwp.Node

	pods := nwp.Pods

	if SurplusResource(nwp, types.GPU) < pod.Gpu {
		return false
	}
	if SurplusResource(nwp, types.CPU) < pod.Cpu {
		return false
	}
	if SurplusResource(nwp, types.RAM) < pod.Ram {
		return false
	}
	if SurplusResource(nwp, types.Disk) < pod.Disk {
		return false
	}

	return node.Eni-len(pods) >= 1

}

func LayoutFillOnePod(allMaxInstancePerNodeLimit map[string]int, nwp types.NodeWithPod, verifyPod types.Pod) bool {

	maxInstancePerNodeLimit := allMaxInstancePerNodeLimit[verifyPod.Group]

	count := 0

	for _, pod := range nwp.Pods {
		if pod.Group == verifyPod.Group {
			count++
		}
	}

	return count+1 <= maxInstancePerNodeLimit

}

func CgroupFillOnePod(nwp types.NodeWithPod, pod types.Pod) types.PodPreAlloc {

	node := nwp.Node

	if len(node.Topologies) == 0 {
		return types.EmptySatisfy
	}

	pods := nwp.Pods

	usedCpuMap := make(map[int]bool)

	for _, pod := range pods {
		for _, cpu := range pod.CpuIds {
			usedCpuMap[cpu] = true
		}
	}

	useableTopologies := make([]types.Topology, 0)

	for _, topology := range node.Topologies {
		if _, ok := usedCpuMap[topology.Cpu]; !ok {
			useableTopologies = append(useableTopologies, topology)
		}
	}

	if len(useableTopologies) < pod.Cpu {
		return types.EmptyNotSatisfy
	}

	socketMap := make(map[int][]types.Topology)

	for _, topology := range useableTopologies {
		if _, ok := socketMap[topology.Socket]; !ok {
			socketMap[topology.Socket] = make([]types.Topology, 0)
		}
		socketMap[topology.Socket] = append(socketMap[topology.Socket], topology)
	}

	for _, socketTopologys := range socketMap {

		if len(socketTopologys) < pod.Cpu {
			continue
		}

		coreMap := make(map[int][]types.Topology)

		for _, topology := range socketTopologys {

			if _, ok := coreMap[topology.Core]; !ok {
				coreMap[topology.Core] = make([]types.Topology, 0)
			}
			coreMap[topology.Core] = append(coreMap[topology.Core], topology)

		}

		if len(coreMap) < pod.Cpu {
			continue
		}

		cpus := make([]int, 0)

		for _, coreTopologys := range coreMap {
			cpus = append(cpus, coreTopologys[0].Cpu)

			if len(cpus) == pod.Cpu {
				break
			}
		}

		return types.PodPreAlloc{
			Satisfy: true,
			Cpus:    cpus,
		}

	}

	return types.EmptyNotSatisfy

}

func RuleOverrunTimeLimit(rule types.Rule, start int64) bool {
	return time.Now().Unix()-start > int64(rule.TimeLimitInMins*constants.MILLISECONDS_4_ONE_MIN)
}

func ScheduleStatisticFrom(directory string, nodeWithPods []types.NodeWithPod, rule types.Rule, groupRuleAssociates []types.GroupRuleAssociate) types.ScheduleStatistic {
	return types.ScheduleStatistic{
		Directory:      directory,
		NodeCount:      len(nodeWithPods),
		PodCount:       PodSize(nodeWithPods),
		GpuAllocation:  Allocation(nodeWithPods, types.GPU),
		CpuAllocation:  Allocation(nodeWithPods, types.CPU),
		RamAllocation:  Allocation(nodeWithPods, types.RAM),
		DiskAllocation: Allocation(nodeWithPods, types.Disk),
		Score:          ScoreNodeWithPods(nodeWithPods, rule, groupRuleAssociates),
	}
}
