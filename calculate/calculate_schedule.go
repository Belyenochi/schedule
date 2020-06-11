package calculate

import (
	"fmt"
	"sort"
	"strconv"
	"django-go/pkg/django"
	"django-go/pkg/types"
	"django-go/pkg/util"
)

type CalculateSchedule struct {
	start int64 //程序启动时间戳
}

func NewSchedule(start int64) django.ScheduleInterface {
	return &CalculateSchedule{start}
}

func (schedule *CalculateSchedule) Schedule(nodes []types.Node, apps []types.App, rule types.Rule) ([]types.ScheduleResult, error) {

	nodeWithPods := sortAndInitNodeWithPods(nodes, rule)

	allPods := sortAndInitPods(apps)

	allMaxInstancePerNodeLimit := util.ToAllMaxInstancePerNodeLimit(rule, types.FromApps(apps))

	schedule.calculate(nodeWithPods, allPods, rule, allMaxInstancePerNodeLimit)

	results := make([]types.ScheduleResult, 0)
	for _, nwp := range nodeWithPods {
		for _, pod := range nwp.Pods {
			result := types.ScheduleResult{
				Sn:     nwp.Node.Sn,
				Group:  pod.Group,
				CpuIds: pod.CpuIds,
			}
			results = append(results, result)
		}
	}

	return results, nil
}

/**
 * 通过node按照规则分数从大到小排序, 然后nodes信息转化为能够容乃下pod对象的NodeWithPod对象，与node对象一对一构建。
 *
 * @param nodes 宿主机列表
 * @param rule  规则对象
 * @return 能够容乃下pod对象的NodeWithPod对象
 */
func sortAndInitNodeWithPods(nodes []types.Node, rule types.Rule) []types.NodeWithPod {

	sort.SliceStable(nodes, func(i, j int) bool {
		return util.ResourceNodeScore(nodes[i], rule) > util.ResourceNodeScore(nodes[j], rule)
	})

	nodeWithPods := make([]types.NodeWithPod, 0, len(nodes))

	for _, node := range nodes {
		nodeWithPods = append(nodeWithPods, types.NodeWithPod{
			Node: node,
			Pods: make([]types.Pod, 0),
		})
	}

	return nodeWithPods
}

/**
 * 通过apps列表对象转化为实际要扩容的pod列表
 *
 * @param apps 应用列表
 * @return 通过app信息构建的pods列表
 */
func sortAndInitPods(apps []types.App) []types.Pod {

	sort.SliceStable(apps, func(i, j int) bool {//排序后的app列表

		if apps[i].Replicas == apps[j].Replicas { //app对应的pod数量从多到少排序

			if apps[i].Gpu == apps[j].Gpu { //按gpu、cpu、内存、磁盘从大到小排序。

				if apps[i].Cpu == apps[j].Cpu {

					if apps[i].Ram == apps[j].Ram {
						return apps[i].Disk > apps[j].Disk
					}
					return apps[i].Ram > apps[j].Ram
				}
				return apps[i].Cpu > apps[j].Cpu
			}
			return apps[i].Gpu > apps[j].Gpu

		}

		return apps[i].Replicas > apps[j].Replicas
	})

	pods := make([]types.Pod, 0)

	for _, app := range apps {
		for i := 0; i < app.Replicas; i++ {
			pods = append(pods, util.ToPod(app))
		}
	}

	fmt.Println("schedule app transform pod count: " + strconv.Itoa(len(pods)))

	return pods
}

func (schedule *CalculateSchedule) calculate(nodeWithPods []types.NodeWithPod, pods []types.Pod, rule types.Rule, allMaxInstancePerNodeLimit map[string]int) {

	forsakePods := make([]types.Pod, 0)//装不下的容器

	for _, pod := range pods {//以打散为目的对排序后的pod、node贪心循环

		forsake := true//遗弃标识

		for i := range nodeWithPods {

			if util.RuleOverrunTimeLimit(rule, schedule.start) {//时间上限约束，超时跳出。
				fmt.Println("overrun time limit")
				return
			}

			if util.StaticFillOnePod(&nodeWithPods[i], &pod, allMaxInstancePerNodeLimit) {
				forsake = false
				break
			}
		}

		if forsake {//所有的机器都不满足此容器分配
			forsakePods = append(forsakePods, pod)
		}
	}

	if len(forsakePods) > 0 {
		fmt.Println("forsake pod count: " + strconv.Itoa(len(forsakePods)))
	}

}
