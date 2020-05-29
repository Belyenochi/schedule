package util

import "django-go/pkg/types"

func NodeWithPodsToPods(nodeWithPods []types.NodeWithPod) []types.Pod {
	allPods := make([]types.Pod, 0)
	for _, node := range nodeWithPods {
		if len(node.Pods) > 0 {
			allPods = append(allPods, node.Pods...)
		}
	}
	return allPods
}

func NodeWithPodsToNodes(nodeWithPods []types.NodeWithPod) []types.Node {
	allNodes := make([]types.Node, 0)
	for _, nodeWithPod := range nodeWithPods {
		allNodes = append(allNodes, nodeWithPod.Node)
	}
	return allNodes
}

func NodesTotalResource(nodes []types.Node, resource types.Resource) int {
	resourceSum := 0
	for _, node := range nodes {
		resourceSum += nodeResourceValue(node, resource)
	}
	return resourceSum
}

func nodeResourceValue(node types.Node, resource types.Resource) int {
	switch resource {
	case types.CPU:
		return node.Cpu
	case types.RAM:
		return node.Ram
	case types.Disk:
		return node.Disk
	default:
		return node.Gpu
	}
}
