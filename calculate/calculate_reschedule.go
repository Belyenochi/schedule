package calculate

import (
	"django-go/pkg/django"
	"django-go/pkg/types"
	"django-go/pkg/util"
	"fmt"
	"strconv"
	"sort"
)

type CalculateReschedule struct {
	start int64 //程序启动时间戳
}

func NewReschedule(start int64) django.RescheduleInterface {
	return &CalculateReschedule{start}
}

func (reschedule *CalculateReschedule) Reschedule(nodeWithPods []types.NodeWithPod, rule types.Rule) ([]types.RescheduleResult, error) {

	nodeWithPods4CheckAgainst := types.CopyNodeWithPods(nodeWithPods)

	againstPods := searchAgainstPods(nodeWithPods4CheckAgainst, rule)

	againstPodCount := 0

	for _, pods := range againstPods {
		againstPodCount += len(pods)
	}

	fmt.Println("reschedule againstPods: " + strconv.Itoa(againstPodCount))

	return reschedule.calculate(nodeWithPods, againstPods, rule), nil
}

func (reschedule *CalculateReschedule) calculate(nodeWithPods []types.NodeWithPod, againstPods map[string][]types.Pod, rule types.Rule) []types.RescheduleResult {

	groupRuleAssociates := types.FromPods(util.ToPods(nodeWithPods))

	allMaxInstancePerNodeLimit := util.ToAllMaxInstancePerNodeLimit(rule, groupRuleAssociates)

	rescheduleResults := make([]types.RescheduleResult, 0)

	forsakePods := make([]types.Pod, 0) //装不下的容器

	var sourceSnKeys []string

	for sourceSn := range againstPods { //golang map range 有随机性，维持遍历againstPods与java版本一致
		sourceSnKeys = append(sourceSnKeys, sourceSn)
	}

	sort.Strings(sourceSnKeys)

	for _, sourceSn := range sourceSnKeys {

		pods := againstPods[sourceSn]

		for _, pod := range pods { //以打散为目的对排序后的pod、node贪心循环

			forsake := true

			for i := range nodeWithPods {

				if util.RuleOverrunTimeLimit(rule, reschedule.start) { //时间上限约束，超时跳出。
					fmt.Println("overrun time limit")
					return rescheduleResults
				}

				if sourceSn == nodeWithPods[i].Node.Sn {
					continue
				}

				if _, ok := againstPods[nodeWithPods[i].Node.Sn]; ok { //有迁移动作的宿主机不在装容器。(明显不优)
					continue
				}

				if util.StaticFillOnePod(&nodeWithPods[i], &pod, allMaxInstancePerNodeLimit) {

					forsake = false

					//动态迁移每一个stage是并行执行，每一次执行必须满足调度约束。
					rescheduleResults = append(rescheduleResults, types.RescheduleResult{
						Stage:    1,
						SourceSn: sourceSn,
						TargetSn: nodeWithPods[i].Node.Sn,
						PodSn:    pod.PodSn,
						CpuIds:   pod.CpuIds,
					})
					break
				}

			}

			if forsake { //所有的机器都不满足此容器分配
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

	//先过滤不满足资源分配的容器，nodeWithPods数据会被修改
	for sn, pods := range searchResourceAgainstPods(nodeWithPods) {
		if old, ok := result[sn]; !ok {
			result[sn] = pods
		} else {
			result[sn] = append(old, pods...)
		}
	}

	//再过滤不满足布局的容器，nodeWithPods数据会被修改
	for sn, pods := range searchLayoutAgainstPods(nodeWithPods, rule) {
		if old, ok := result[sn]; !ok {
			result[sn] = pods
		} else {
			result[sn] = append(old, pods...)
		}
	}

	//再过滤不满足cpu分配的容器，nodeWithPods数据会被修改
	for sn, pods := range searchCgroupAgainstPods(nodeWithPods) {
		if old, ok := result[sn]; !ok {
			result[sn] = pods
		} else {
			result[sn] = append(old, pods...)
		}
	}

	return result

}

//违背资源规则容器准备重新调度。<node_sn,List<Pod>>
func searchResourceAgainstPods(nodeWithPods []types.NodeWithPod) map[string][]types.Pod {

	result := make(map[string][]types.Pod, 0)

	for _, nwp := range nodeWithPods {

		againstPods := make([]types.Pod, 0)

		tempPods := make([]types.Pod, 0)

		normalPods := make([]types.Pod, 0)

		//校验资源不满足的容器
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

		//校验超过eni约束的容器
		if eniAgainstPodSize > 0 {

			againstPods = append(againstPods, tempPods[0:eniAgainstPodSize]...)

			normalPods = append(normalPods, tempPods[eniAgainstPodSize:len(tempPods)-1]...)

		} else {
			normalPods = append(normalPods, tempPods...)
		}

		nwp.Pods = normalPods //贪心判断的正常容器继续放在该机器中

		if len(againstPods) > 0 {
			result[nwp.Node.Sn] = againstPods
		}

	}

	return result
}

//违背布局规则容器准备重新调度。<node_sn,List<Pod>>
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

		nwp.Pods = normalPods //贪心判断的正常容器继续放在该机器中

		if len(againstPods) > 0 {
			result[nwp.Node.Sn] = againstPods
		}

	}

	return result

}

//违背cpu绑核分配规则容器准备重新调度。<node_sn,List<Pod>>
func searchCgroupAgainstPods(nodeWithPods []types.NodeWithPod) map[string][]types.Pod {

	result := make(map[string][]types.Pod, 0)

	for _, nwp := range nodeWithPods {

		node := nwp.Node

		//node中不存在topologies,不校验绑核。
		if len(node.Topologies) == 0 {
			continue
		}

		//node中不存在topologies,不校验绑核。
		if len(nwp.Pods) == 0 {
			continue
		}

		againstPods := make([]types.Pod, 0)

		//这台机器上重叠的cpuId分配
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

			if len(pod.CpuIds) == 0 { //没有分配cpuId的容器
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

			if socketCount > 1 { //跨socket容器

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

		nwp.Pods = normalPods //贪心判断的正常容器继续放在该机器中

		if len(againstPods) > 0 {
			result[nwp.Node.Sn] = againstPods
		}

	}

	return result

}
