package main

import (
	"fmt"
	"sync"
	"time"

	"django-go/pkg/loader"
	"django-go/pkg/types"
	"django-go/pkg/util"
	"django-go/cmd/utils"
	"os"
	"django-go/pkg/store"
	"django-go/calculate"
)

/**
 * 运行静态布局和动态迁移功能
 * 使用其他目录数据，可以硬编码args = new String[] {"data_2"}做测试，或通过IDEA传参的方式启动main方法;
 */
func main() {

	start := time.Now()

	directorys := utils.AdjustDirectorys(os.Args[1:])

	execute(start, directorys)

	fmt.Println(fmt.Sprintf("finish calculate, total use time : %v/s", time.Now().Sub(start).Seconds()))

	utils.Score()

	time.Sleep(2 * time.Hour) //2小时避免程序退出

}

func execute(start time.Time, directories []string) {

	wg := sync.WaitGroup{}

	wg.Add(2 * len(directories))

	maxTimeLimitInMins := 0

	for _, dir := range directories {

		directory := dir

		dataLoader := loader.NewLoader(directory)

		rule, err := dataLoader.LoadRule() //加载调度规则数据'rule.source'
		util.MustBeTrue(err == nil, fmt.Sprintf("load rule error, msg:%v", err))

		nodes, err := dataLoader.LoadNodes() //加载静态布局机器数据'schedule.node.source'
		util.MustBeTrue(err == nil, fmt.Sprintf("load nodes error, msg:%v", err))

		apps, err := dataLoader.LoadApps() //加载静态布局应用数据'schedule.app.source'
		util.MustBeTrue(err == nil, fmt.Sprintf("load apps error, msg:%v", err))

		nodeWithPods, err := dataLoader.LoadNodeWithPods() //加载动态迁移原始数据'reschedule.source'
		util.MustBeTrue(err == nil, fmt.Sprintf("load node with pods error, msg:%v", err))

		maxTimeLimitInMins = util.Max(maxTimeLimitInMins, rule.TimeLimitInMins)

		go func() {

			defer wg.Done()

			schedule(directory, start.Unix(), rule, nodes, apps)

		}()

		go func() {

			defer wg.Done()

			reschedule(directory, start.Unix(), rule, nodeWithPods)

		}()
	}

	wg.Wait()
}

func schedule(directory string, start int64, rule types.Rule, nodes []types.Node, apps []types.App) {

	schedule := calculate.NewSchedule(start)

	fmt.Println(fmt.Sprintf("%s | schedule source total score : %v", directory, util.ResourceNodesScore(nodes, rule)))

	//执行静态布局功能
	results, err := schedule.Schedule(types.CopyNodes(nodes), types.CopyApps(apps), rule.Copy())

	if err != nil {
		fmt.Println("schedule err, msg:" + err.Error())
		return
	}

	store.StoreSchedule(results, directory)

	nodeWithPods := util.ResultToNodeWithPods(nodes, apps, results)

	statistic := util.ScheduleStatisticFrom(directory, nodeWithPods, rule, types.FromApps(apps))

	statistic.Log("schedule result")
}

func reschedule(directory string, start int64, rule types.Rule, nodeWithPods []types.NodeWithPod) {

	groupRuleAssociates := types.FromPods(util.ToPods(nodeWithPods))

	statistic := util.ScheduleStatisticFrom(directory, nodeWithPods, rule, groupRuleAssociates)

	statistic.Log("reschedule source")

	reschedule := calculate.NewReschedule(start)

	//执行动态迁移功能
	results, err := reschedule.Reschedule(types.CopyNodeWithPods(nodeWithPods), rule.Copy())

	if err != nil {
		fmt.Println("%s | reschedule error, msg:" + err.Error())
		return
	}

	store.StoreReschedule(results, directory)

}
