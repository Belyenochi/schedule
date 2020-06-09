package util

import "django-go/pkg/types"

/**
 * 计算一批宿主机某一项资源的总数
 */
func NodesTotalResource(nodes []types.Node, resource types.Resource) int {
	sum := 0
	for _, node := range nodes {
		sum += node.Value(resource)
	}
	return sum
}

/**
 * 根据node topologie信息映射<cpu,socket>
 */
func CpuToSocket(node types.Node) map[int]int {
	result := make(map[int]int, 0)
	for _, topologie := range node.Topologies {
		result[topologie.Cpu] = topologie.Socket
	}
	return result
}

/**
 * 根据node topologie信息映射<cpu,core>
 */
func CpuToCore(node types.Node) map[int]int {
	result := make(map[int]int, 0)
	for _, topologie := range node.Topologies {
		result[topologie.Cpu] = topologie.Core
	}
	return result
}
