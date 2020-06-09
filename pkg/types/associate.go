package types

type GraphMigrateNode struct {
	PodSn string `json:"podSn"`

	NodeSn string `json:"nodeSn"`

	CpuIDs []int `json:"cpuIds"`

	Type OperationType `json:"type"`
}

type OperationType string

const (
	EXPANSION OperationType = "EXPANSION"
	OFFLINE   OperationType = "OFFLINE"
)

type GroupRuleAssociate struct {
	Group    string `json:"group"`
	Replicas int    `json:"replicas"`
}

type PodPreAlloc struct {
	Satisfy bool  `json:"satisfy"`
	Cpus    []int `json:"cpus"`
}

var EmptySatisfy = PodPreAlloc{
	Satisfy: true,
	Cpus:    make([]int, 0),
}

var EmptyNotSatisfy = PodPreAlloc{
	Satisfy: false,
	Cpus:    make([]int, 0),
}

func FromApps(apps []App) []GroupRuleAssociate {
	groupRuleAssociates := make([]GroupRuleAssociate, 0)
	for _, app := range apps {
		groupRuleAssociates = append(groupRuleAssociates, GroupRuleAssociate{
			Group:    app.Group,
			Replicas: app.Replicas,
		})
	}
	return groupRuleAssociates
}

func FromPods(pods []Pod) []GroupRuleAssociate {

	groupCountMap := make(map[string]int)

	for _, pod := range pods {
		if _, ok := groupCountMap[pod.Group]; !ok {
			groupCountMap[pod.Group] = 1
			continue
		}
		groupCountMap[pod.Group] = groupCountMap[pod.Group] + 1
	}

	groupRuleAssociates := make([]GroupRuleAssociate, 0)
	for group, replicas := range groupCountMap {
		groupRuleAssociates = append(groupRuleAssociates, GroupRuleAssociate{
			Group:    group,
			Replicas: replicas,
		})
	}
	return groupRuleAssociates
}

func ToExpansionMigrateNode(podSn string, nodeSn string, cpuIds []int) GraphMigrateNode {
	return GraphMigrateNode{
		PodSn:  podSn,
		NodeSn: nodeSn,
		CpuIDs: cpuIds,
		Type:   EXPANSION,
	}
}

func ToOfflineMigrateNode(podSn string, nodeSn string) GraphMigrateNode {
	return GraphMigrateNode{
		PodSn:  podSn,
		NodeSn: nodeSn,
		CpuIDs: make([]int, 0),
		Type:   OFFLINE,
	}
}
