package types

type Rule struct {
	TimeLimitInMins             int                          `json:"timeLimitInMins"`
	DefaultMaxInstancePerNode   int                          `json:"defaultMaxInstancePerNode"`
	ScoreWeight                 ScoreWeight                  `json:"scoreWeight"`
	ReplicasMaxInstancePerNodes []ReplicasMaxInstancePerNode `json:"replicasMaxInstancePerNodes"`
	NodeResourceWeights         []ResourceWeight             `json:"nodeResourceWeights"`
}

type ScoreWeight struct {
	MigratePod               int `json:"migratePod"`
	GroupMoreInstancePerNode int `json:"groupMoreInstancePerNode"`
	SocketCross              int `json:"socketCross"`
	CoreBind                 int `json:"coreBind"`
	SensitiveCpuBind         int `json:"sensitiveCpuBind"`
}

type Restrain string

const (
	LE Restrain = "le"//小于或等于
	GE Restrain = "ge"//大于或等于
)

var AllRestrains = []string{string(LE), string(GE)}

type ReplicasMaxInstancePerNode struct {
	Replicas           int      `json:"replicas"`
	Restrain           Restrain `json:"restrain"`
	MaxInstancePerNode int      `json:"maxInstancePerNode"`
}

type Resource string

const (
	GPU  Resource = "GPU"
	CPU  Resource = "CPU"
	RAM  Resource = "RAM"
	Disk Resource = "Disk"
)

var AllResources = []Resource{GPU, CPU, RAM, Disk}

type ResourceWeight struct {
	Resource      Resource `json:"resource"`
	Weight        int      `json:"weight"`
	NodeModelName string   `json:"nodeModelName"`
}

func (sw ScoreWeight) Copy() ScoreWeight {
	return ScoreWeight{
		MigratePod:               sw.MigratePod,
		GroupMoreInstancePerNode: sw.GroupMoreInstancePerNode,
		SocketCross:              sw.SocketCross,
		CoreBind:                 sw.CoreBind,
		SensitiveCpuBind:         sw.SensitiveCpuBind,
	}
}

func (r ReplicasMaxInstancePerNode) Copy() ReplicasMaxInstancePerNode {
	return ReplicasMaxInstancePerNode{
		Replicas:           r.Replicas,
		Restrain:           r.Restrain,
		MaxInstancePerNode: r.MaxInstancePerNode,
	}
}
func (rw ResourceWeight) Copy() ResourceWeight {
	return ResourceWeight{
		Resource:      rw.Resource,
		Weight:        rw.Weight,
		NodeModelName: rw.NodeModelName,
	}
}

func copyReplicasMaxInstancePerNodes(replicasMaxInstancePerNodes []ReplicasMaxInstancePerNode) []ReplicasMaxInstancePerNode {

	result := make([]ReplicasMaxInstancePerNode, 0)

	for _, r := range replicasMaxInstancePerNodes {
		result = append(result, r.Copy())
	}

	return result
}

func copyNodeResourceWeights(nodeResourceWeights []ResourceWeight) []ResourceWeight {

	result := make([]ResourceWeight, 0)

	for _, n := range nodeResourceWeights {
		result = append(result, n.Copy())
	}

	return result
}

func (rule Rule) Copy() Rule {
	return Rule{
		TimeLimitInMins:             rule.TimeLimitInMins,
		DefaultMaxInstancePerNode:   rule.DefaultMaxInstancePerNode,
		ScoreWeight:                 rule.ScoreWeight.Copy(),
		ReplicasMaxInstancePerNodes: copyReplicasMaxInstancePerNodes(rule.ReplicasMaxInstancePerNodes),
		NodeResourceWeights:         copyNodeResourceWeights(rule.NodeResourceWeights),
	}
}
