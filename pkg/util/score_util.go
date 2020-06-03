package util

import (
	"django-go/pkg/types"
	"django-go/pkg/constants"
	"errors"
	"fmt"
)

//golang不支持方法重载，命名与java有区别
func ResourceNodeScore(node types.Node, rule types.Rule) int {
	return ResourceNodesScore([]types.Node{node}, rule)
}

func ResourceNodesScore(nodes []types.Node, rule types.Rule) int {

	scoreSum := 0

	scoreMap := toResourceScoreMap(rule)

	for _, node := range nodes {
		scoreSum += resourceScore(scoreMap, node)
	}

	return scoreSum
}

func ScoreNodeWithPods(nodeWithPods []types.NodeWithPod, rule types.Rule, groupRuleAssociates []types.GroupRuleAssociate) int {

	if !dataSatisfy(nodeWithPods, groupRuleAssociates) {
		return constants.INVALID_SCORE
	}

	if !resourceSatisfy(nodeWithPods) {
		return constants.INVALID_SCORE
	}

	layoutScore := layoutScore(nodeWithPods, rule, groupRuleAssociates)

	if layoutScore < 0 {
		return constants.INVALID_SCORE
	}

	cgroupScore := cgroupScore(nodeWithPods, rule)

	if cgroupScore < 0 {
		return constants.INVALID_SCORE
	}

	resourceScore := ResourceNodesScore(ToNodes(nodeWithPods), rule)

	fmt.Println(fmt.Sprintf("layoutScore:%v , cgroupScore:%v , resourceScore:%v ", layoutScore, cgroupScore, resourceScore))

	return resourceScore + layoutScore + cgroupScore
}

func cgroupScore(nodeWithPods []types.NodeWithPod, rule types.Rule) int {

	sw := rule.ScoreWeight

	totalSocketCrossCount := 0
	totalCoreBindCount := 0
	totalSensitiveCpuBindCount := 0

	for _, nwp := range nodeWithPods {

		node := nwp.Node

		if len(node.Topologies) == 0 {
			continue
		}

		pods := nwp.Pods

		for _, pod := range pods {
			if len(pod.CpuIds) == 0 {
				fmt.Println(errors.New(fmt.Sprintf("go against cgroup | node:%v have pod no cpuIDs", node.Sn)))
				return constants.INVALID_SCORE
			}

			if len(pod.CpuIds) != pod.Cpu {
				fmt.Println(errors.New(fmt.Sprintf("go against cgroup | node:%v have pod cpuID size unequal to pod cpu", node.Sn)))
				return constants.INVALID_SCORE
			}
		}

		cpuIDCountAgainstMap := make(map[int]int)

		for key, value := range CpuIDCountMap(nwp) {
			if value > 1 {
				cpuIDCountAgainstMap[key] = value
			}
		}

		if len(cpuIDCountAgainstMap) > 0 {
			fmt.Println(errors.New(fmt.Sprintf("go against cgroup | node:%v cpuIds overlap: %v", node.Sn, ToJsonOrDie(cpuIDCountAgainstMap))))
			return constants.INVALID_SCORE
		}

		cpuToSocket := CpuToSocket(node)

		for _, pod := range pods {
			for _, cpuId := range pod.CpuIds {
				if _, ok := cpuToSocket[cpuId]; !ok { //socket不存在表示此容器cpuId不是来自node的topology
					fmt.Println(errors.New(fmt.Sprintf("go against cgroup | node:%v cpuId :%v invalid", node.Sn, cpuId)))
					return constants.INVALID_SCORE
				}
			}
		}

		for _, pod := range pods {

			socketCountMap := make(map[int]int)

			for _, cpuId := range pod.CpuIds {

				if socket, ok := cpuToSocket[cpuId]; !ok {
					socketCountMap[socket] = 1
				} else {
					socketCountMap[socket] = socketCountMap[socket] + 1
				}

			}

			if len(socketCountMap) > 1 {

				defaultKey := -1

				defaultValue := -1

				for key, value := range socketCountMap {

					if value > defaultValue {
						defaultKey = key
						defaultValue = value
					}

				}

				delete(socketCountMap, defaultKey)

				for _, value := range socketCountMap {
					totalSocketCrossCount += value
				}

			}

		}

		cpuToCore := CpuToCore(node)

		for _, pod := range pods {

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

			for _, value := range sameCoreMap {
				totalCoreBindCount += value - 1
			}

		}

		//TODO 已经很复杂了，如果排名拉不开差距在增加sensitiveCpuBind数据校验
	}

	return totalSocketCrossCount*sw.SocketCross + totalCoreBindCount*sw.CoreBind + totalSensitiveCpuBindCount*sw.SensitiveCpuBind

}

func dataSatisfy(nodeWithPods []types.NodeWithPod, groupRuleAssociates []types.GroupRuleAssociate) bool {

	sourceGroupReplicas := make(map[string]int, len(groupRuleAssociates))

	for _, groupRuleAssociate := range groupRuleAssociates {
		sourceGroupReplicas[groupRuleAssociate.Group] = groupRuleAssociate.Replicas
	}

	pods := ToPods(nodeWithPods)

	groupCountMap := make(map[string]int)

	for _, pod := range pods {
		if _, ok := groupCountMap[pod.Group]; !ok {
			groupCountMap[pod.Group] = 1
			continue
		}
		groupCountMap[pod.Group] = groupCountMap[pod.Group] + 1
	}

	keys := make([]string, 0, len(sourceGroupReplicas))
	for k := range sourceGroupReplicas {
		keys = append(keys, k)
	}

	for _, group := range keys {
		if sourceGroupReplicas[group] != groupCountMap[group] {
			return false
		}
	}

	return true

}

func resourceSatisfy(nodeWithPods []types.NodeWithPod) bool {

	for _, nwp := range nodeWithPods {

		node := nwp.Node

		pods := nwp.Pods

		if node.Eni < len(pods) {
			fmt.Println(errors.New(fmt.Sprintf("go against eni alloc | node:%s , eni: %v ,podSize: %v", node.Sn, node.Eni, len(pods))))
			return false
		}

		for _, resource := range types.AllResources {

			nodeResource := node.Value(resource)

			podsResource := PodsTotalResource(pods, resource)

			if nodeResource < podsResource {
				fmt.Println(errors.New(fmt.Sprintf("go against resource alloc | node:%v ,resource: %v ,nodeResource:%v , podsResource:%v",
					node.Sn, resource, nodeResource, podsResource)))
				return false
			}
		}

	}

	return true

}

func layoutScore(nodeWithPods []types.NodeWithPod, rule types.Rule, groupRuleAssociates []types.GroupRuleAssociate) int {

	totalGroupMoreInstancePerNodeCount := 0

	maxInstancePerNodes := ToAllMaxInstancePerNodeLimit(rule, groupRuleAssociates)

	for _, nwp := range nodeWithPods {

		groupCountPreNodeMap := GroupCountPreNodeMap(nwp)

		for key, value := range groupCountPreNodeMap {

			if value > maxInstancePerNodes[key] { //当前机器和应用分组超过允许堆叠的数量，当前布局无效。

				fmt.Println(errors.New(fmt.Sprintf("go against layout | node:%v , group:%v , count:%v , maxInstancePerNode:%v",
					nwp.Node.Sn, key, value, maxInstancePerNodes[key])))

				return constants.INVALID_SCORE
			}

			if value > 1 { //当前宿主机下，若此应用分组布局多于一个。
				totalGroupMoreInstancePerNodeCount += value - 1
			}

		}

	}

	return totalGroupMoreInstancePerNodeCount * rule.ScoreWeight.GroupMoreInstancePerNode

}

func resourceScore(scoreMap map[types.Resource]map[string]int, node types.Node) int {

	sumScore := 0

	for _, r := range types.AllResources {

		wm, ok := scoreMap[r]

		if !ok {
			wm = make(map[string]int)
		}

		if w, ok := wm[node.NodeModelName]; ok {
			sumScore += w * node.Value(r)
			continue
		}

		if w, ok := wm[constants.SCORE_EMPTY_NODE_MODEL_NAME]; ok {
			sumScore += w * node.Value(r)
		}

	}

	return sumScore
}

func ScoreReschedule(resultList []types.RescheduleResult, rule types.Rule, sourceList []types.NodeWithPod) int {
	return 0
}
