package types

import (
	"fmt"
)

type ScheduleStatistic struct {
	Directory      string  `json:"directory"`
	NodeCount      int     `json:"nodeCount"`
	PodCount       int     `json:"podCount"`
	GpuAllocation  float64 `json:"gpuAllocation"`
	CpuAllocation  float64 `json:"cpuAllocation"`
	RamAllocation  float64 `json:"ramAllocation"`
	DiskAllocation float64 `json:"diskAllocation"`
	Score          int     `json:"score"`
}
func (statistic *ScheduleStatistic) Log(label string) {

	fmt.Println(fmt.Sprintf("schedule statistic | dir: %s label: %v | count:【node %v, pod %v】", statistic.Directory, label, statistic.NodeCount, statistic.PodCount))

	fmt.Println(fmt.Sprintf("schedule statistic | dir: %s label: %v | allocation:【gpu %f%%, cpu %f%%, ram %f%%, disk %f%%】", statistic.Directory, label, statistic.GpuAllocation, statistic.CpuAllocation, statistic.RamAllocation, statistic.DiskAllocation))

	fmt.Println(fmt.Sprintf("schedule statistic | dir: %s label: %v | total score: %v", statistic.Directory, label, statistic.Score))

}
