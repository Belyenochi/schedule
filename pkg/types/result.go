package types

type ScheduleResult struct {
	Sn     string `json:"sn"`
	Group  string `json:"group"`
	CpuIds []int  `json:"cpuIDs"`
}

type RescheduleResult struct {
	Stage    int    `json:"stage"`
	PodSn    string `json:"podSn"`
	SourceSn string `json:"sourceSn"`
	TargetSn string `json:"targetSn"`
	CpuIds   []int  `json:"cpuIDs"`
}

type ScoreResult struct {
	TotalScheduleScore   int `json:"totalScheduleScore"`
	TotalRescheduleScore int `json:"totalRescheduleScore"`
}