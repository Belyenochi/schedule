package util

import (
	"django-go/pkg/types"
)

func PodSize(nwps []types.NodeWithPod) int {
	sum := 0
	for _, nwp := range nwps {
		sum += len(nwp.Pods)
	}
	return sum
}

func ToPods(nwps []types.NodeWithPod) []types.Pod {
	result := make([]types.Pod, 0)
	for _, nwp := range nwps {
		if len(nwp.Pods) > 0 {
			result = append(result, nwp.Pods...)
		}
	}
	return result
}

func ToPodMap(nwps []types.NodeWithPod) map[string]types.Pod {
	result := make(map[string]types.Pod)
	for _, pod := range ToPods(nwps) {
		result[pod.PodSn] = pod
	}
	return result
}

func ToNodes(nwps []types.NodeWithPod) []types.Node {
	result := make([]types.Node, 0)
	for _, nwp := range nwps {
		result = append(result, nwp.Node)
	}
	return result
}

func ToNodeMap(nwps []types.NodeWithPod) map[string]types.Node {
	result := make(map[string]types.Node)
	for _, node := range ToNodes(nwps) {
		result[node.Sn] = node
	}
	return result
}

func SurplusResource(nwp types.NodeWithPod, resource types.Resource) int {
	return nwp.Node.Value(resource) - PodsTotalResource(nwp.Pods, resource)
}

func Allocation(nwps []types.NodeWithPod, resource types.Resource) float64 {

	nodeResource := NodesTotalResource(ToNodes(nwps), resource)

	if nodeResource == 0 {
		return 0
	}

	podsResource := PodsTotalResource(ToPods(nwps), resource)

	return 100 * float64(podsResource) / float64(nodeResource)

}

func GroupCountPreNodeMap(nwp types.NodeWithPod) map[string]int {
	result := make(map[string]int)
	for _, pod := range nwp.Pods {
		if _, ok := result[pod.Group]; !ok {
			result[pod.Group] = 1
			continue
		}
		result[pod.Group] = result[pod.Group] + 1
	}
	return result
}

func CpuIDCountMap(nwp types.NodeWithPod) map[int]int {
	result := make(map[int]int)
	for _, pod := range nwp.Pods {
		for _, cpu := range pod.CpuIds {
			if _, ok := result[cpu]; !ok {
				result[cpu] = 1
				continue
			}
			result[cpu] = result[cpu] + 1
		}
	}
	return result
}

func ResultToNodeWithPods(nodes []types.Node, apps []types.App, result []types.ScheduleResult) []types.NodeWithPod {

	nodeMap := make(map[string]types.Node)

	for _, node := range nodes {
		nodeMap[node.Sn] = node
	}

	appMap := make(map[string]types.App)

	for _, app := range apps {
		appMap[app.Group] = app
	}

	nwps := make([]types.NodeWithPod, 0)

	resultSnMap := make(map[string][]types.ScheduleResult)

	for _, sr := range result {
		resultSnMap[sr.Sn] = append(resultSnMap[sr.Sn], sr)
	}

	for sn, list := range resultSnMap {

		node := nodeMap[sn]

		pods := make([]types.Pod, 0)

		for _, r := range list {

			pod := ToPod(appMap[r.Group])

			pod.CpuIds = r.CpuIds

			pods = append(pods, pod)
		}

		nwps = append(nwps, types.NodeWithPod{
			Node: node,
			Pods: pods,
		})

	}

	return nwps

}
