package types

type Rule struct {
	TimeLimitInMins             int                          `json:"timeLimitInMins"`
	DefaultMaxInstancePerHost   int                          `json:"defaultMaxInstancePerHost"`
	GroupMaxInstancePerHosts    []GroupMaxInstancePerHost    `json:"groupMaxInstancePerHosts"`
	ReplicasMaxInstancePerHosts []ReplicasMaxInstancePerHost `json:"replicasMaxInstancePerHosts"`
	NodeResourceWeights         []ResourceWeight             `json:"nodeResourceWeights"`
}

type GroupMaxInstancePerHost struct {
	Group              string `json:"group"`
	MaxInstancePerHost int    `json:"maxInstancePerHost"`
	Compactness        bool   `json:"compactness"`
}

type Restrain string

const (
	LE Restrain = "le"
	GE Restrain = "ge"
)

var AllRestrains = []string{string(LE), string(GE)}

type ReplicasMaxInstancePerHost struct {
	Replicas           int      `json:"replicas"`
	Restrain           Restrain `json:"restrain"`
	MaxInstancePerHost int      `json:"maxInstancePerHost"`
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
	Resource         Resource `json:"resource"`
	Weight           int      `json:"weight"`
	MachineModelName string   `json:"machineModelName"`
}
