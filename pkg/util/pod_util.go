package util

import "django-go/pkg/types"

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

func PodsTotalResource(pods []types.Pod, resource types.Resource) int {
	sum := 0
	for _, pod := range pods {
		sum += pod.Value(resource)
	}
	return sum
}
