package util

import "django-go/pkg/types"

func BuildPodForApp(app types.App) types.Pod {
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
	resourceSum := 0
	for _, pod := range pods {
		resourceSum += podResourceValue(pod, resource)
	}
	return resourceSum
}

func podResourceValue(pod types.Pod, resource types.Resource) int {
	switch resource {
	case types.Disk:
		return pod.Disk
	case types.RAM:
		return pod.Ram
	case types.CPU:
		return pod.Cpu
	default:
		return pod.Gpu
	}
}
