package util

import "django-go/pkg/types"

func NodesTotalResource(nodes []types.Node, resource types.Resource) int {
	sum := 0
	for _, node := range nodes {
		sum += node.Value(resource)
	}
	return sum
}

func CpuToSocket(node types.Node) map[int]int {
	result := make(map[int]int, 0)
	for _, topologie := range node.Topologies {
		result[topologie.Cpu] = topologie.Socket
	}
	return result
}

func CpuToCore(node types.Node) map[int]int {
	result := make(map[int]int, 0)
	for _, topologie := range node.Topologies {
		result[topologie.Cpu] = topologie.Core
	}
	return result
}
