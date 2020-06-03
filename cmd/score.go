package main

import (
	"time"
	"django-go/cmd/utils"
	"os"
	"fmt"
	"django-go/pkg/loader"
	"django-go/pkg/util"
	"django-go/pkg/types"
	"django-go/pkg/store"
	"django-go/pkg/constants"
)

func main() {
	start := time.Now()

	score(utils.AdjustDirectorys(os.Args[1:]))

	fmt.Println(fmt.Sprintf("finish score, total use time : %v/s", time.Now().Sub(start).Seconds()))
}

func score(directorys []string) {

	totalScheduleScore := 0

	totalRescheduleScore := 0

	for _, directory := range directorys {

		dataLoader := loader.NewLoader(directory)

		rule, err := dataLoader.LoadRule()
		util.MustBeTrue(err == nil, fmt.Sprintf("load rule error, msg:%v", err))

		nodes, err := dataLoader.LoadNodes()
		util.MustBeTrue(err == nil, fmt.Sprintf("load nodes error, msg:%v", err))

		apps, err := dataLoader.LoadApps()
		util.MustBeTrue(err == nil, fmt.Sprintf("load apps error, msg:%v", err))

		nodeWithPods, err := dataLoader.LoadNodeWithPods()
		util.MustBeTrue(err == nil, fmt.Sprintf("load node with pods error, msg:%v", err))

		scheduleResults, err := store.LoadScheduleResults(directory)
		util.MustBeTrue(err == nil, fmt.Sprintf("load schedule result error, msg:%v", err))

		rescheduleResults, err := store.LoadRescheduleResults(directory)
		util.MustBeTrue(err == nil, fmt.Sprintf("load reschedule result error, msg:%v", err))

		scheduleScore := scheduleScore(directory, scheduleResults, rule, nodes, apps)

		if scheduleScore == constants.INVALID_SCORE || totalScheduleScore == constants.INVALID_SCORE {
			totalScheduleScore = constants.INVALID_SCORE
		} else {
			totalScheduleScore += scheduleScore
		}

		rescheduleScore := rescheduleScore(directory, rescheduleResults, rule, nodeWithPods); //评测当前目录下动态迁移功能

		if rescheduleScore == constants.INVALID_SCORE || totalRescheduleScore == constants.INVALID_SCORE {
			totalRescheduleScore = constants.INVALID_SCORE
		} else {
			totalRescheduleScore += rescheduleScore
		}

	}

	fmt.Println(fmt.Sprintf("total schedule score : %v , reschedule score : %v", totalScheduleScore, totalRescheduleScore))
	fmt.Println(fmt.Sprintf("ScoreResult:%v", util.ToJsonOrDie(types.ScoreResult{
		TotalScheduleScore:   totalScheduleScore,
		TotalRescheduleScore: totalRescheduleScore,
	})))

}

func scheduleScore(directory string, scheduleResults []types.ScheduleResult, rule types.Rule, nodes []types.Node, apps []types.App) int {

	if len(scheduleResults) == 0 {
		return constants.INVALID_SCORE
	}

	nodeWithPods := util.ResultToNodeWithPods(nodes, apps, scheduleResults)

	scheduleScore := util.ScoreNodeWithPods(nodeWithPods, rule, types.FromApps(apps))

	fmt.Println(fmt.Sprintf("%v | schedule result total score : %v", directory, scheduleScore))

	return scheduleScore

}

func rescheduleScore(directory string, rescheduleResults []types.RescheduleResult, rule types.Rule, nodeWithPods []types.NodeWithPod) int {

	if len(rescheduleResults) == 0 {
		return constants.INVALID_SCORE
	}

	rescheduleScore := util.ScoreReschedule(rescheduleResults, rule, nodeWithPods)

	fmt.Println(fmt.Sprintf("%v | reschedule result total score : %v", directory, rescheduleScore))

	return rescheduleScore

}
