package main

import (
	"fmt"
	"sync"
	"time"

	"django-go/calculate"
	"django-go/pkg/django"
	"django-go/pkg/loader"
	"django-go/pkg/types"
	"django-go/pkg/util"
)

func main() {
	execute([]string{"data_test"})
}

func execute(directories []string) {
	wg := sync.WaitGroup{}
	wg.Add(2 * len(directories))

	for _, dir := range directories {

		directory := dir
		dataLoader := loader.NewLoader(directory)

		nodes, err := dataLoader.LoadNodes()
		util.MustBeTrue(err == nil, fmt.Sprintf("load nodes error, msg:%v", err))
		apps, err := dataLoader.LoadApps()
		util.MustBeTrue(err == nil, fmt.Sprintf("load apps error, msg:%v", err))
		rule, err := dataLoader.LoadRule()
		util.MustBeTrue(err == nil, fmt.Sprintf("load rule error, msg:%v", err))
		nodeWithPods, err := dataLoader.LoadNodeWithPods()
		util.MustBeTrue(err == nil, fmt.Sprintf("load node with pods error, msg:%v", err))

		go func() {
			defer wg.Done()

			start := time.Now()

			schedule := calculate.NewSchedule()
			results, err := schedule.Schedule(nodes, apps, rule)
			if err != nil {
				fmt.Print("schedule err, msg:" + err.Error())
				return
			}

			fmt.Println(fmt.Sprintf("%s | schedule source app count : %v", directory, len(apps)))
			fmt.Println(fmt.Sprintf("%s | schedule source node count : %v", directory, len(nodes)))
			fmt.Println(fmt.Sprintf("%s | schedule source total score : %v", directory, util.ScoreNodes(nodes, rule)))

			nodeWithPodsResult := toNodeWithPods(nodes, apps, results)

			fmt.Println(fmt.Sprintf("%s | schedule result use time: %v s", directory, time.Now().Sub(start).Seconds()))
			fmt.Println(fmt.Sprintf("%s | schedule result : %s", directory, util.ToJsonOrDie(nodeWithPodsResult)))
		}()

		go func() {
			defer wg.Done()

			start := time.Now()

			schedule := calculate.NewSchedule()
			results, err := schedule.Reschedule(nodeWithPods, rule)
			if err != nil {
				fmt.Println("%s | reschedule error, msg:" + err.Error())
				return
			}

			fmt.Println(fmt.Sprintf("%s | reschedule result use time: %v s", directory, time.Now().Sub(start).Seconds()))
			fmt.Println(fmt.Sprintf("%s | reschedule result : %s", directory, util.ToJsonOrDie(results)))

		}()
	}

	wg.Wait()
	fmt.Println("finish calculate")
}

func toNodeWithPods(nodes []types.Node, apps []types.App, results []django.ScheduleResult) []types.NodeWithPod {
	snNodeMap := make(map[string]types.Node, len(nodes))
	for _, node := range nodes {
		snNodeMap[node.Sn] = node
	}

	groupAppMap := make(map[string]types.App, len(apps))
	for _, app := range apps {
		groupAppMap[app.Group] = app
	}

	snResultMap := make(map[string][]django.ScheduleResult)
	for _, result := range results {
		if resultList, ok := snResultMap[result.Sn]; ok {
			resultList = append(resultList, result)
			snResultMap[result.Sn] = resultList
		} else {
			snResultMap[result.Sn] = []django.ScheduleResult{result}
		}
	}

	nodeWithPodsResult := make([]types.NodeWithPod, 0)
	for sn, results := range snResultMap {
		node := snNodeMap[sn]

		pods := make([]types.Pod, 0)
		for _, result := range results {
			app := groupAppMap[result.Group]
			pods = append(pods, util.BuildPodForApp(app))
		}

		nodeWithPodsResult = append(nodeWithPodsResult, types.NodeWithPod{Node: node, Pods: pods})
	}

	return nodeWithPodsResult
}
