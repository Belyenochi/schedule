package util

import "django-go/pkg/types"

/**
 * 通过app中的信息构建一个Pod对象
 */
func ToPod(app types.App) types.Pod {
	return types.Pod{
		AppName: app.AppName,
		Group:   app.Group,
		Gpu:     app.Gpu,
		Cpu:     app.Cpu,
		Ram:     app.Ram,
		Disk:    app.Disk,
	}
}

/**
 * 计算一批Pod在某一项资源的总数
 */
func PodsTotalResource(pods []types.Pod, resource types.Resource) int {
	sum := 0
	for _, pod := range pods {
		sum += pod.Value(resource)
	}
	return sum
}
