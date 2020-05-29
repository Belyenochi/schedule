package types

type App struct {
	AppName  string `json:"appName"`
	Group    string `json:"group"`
	Gpu      int    `json:"gpu"`
	Cpu      int    `json:"cpu"`
	Ram      int    `json:"ram"`
	Disk     int    `json:"disk"`
	Replicas int    `json:"replicas"`
}

type Node struct {
	Sn               string     `json:"sn"`
	MachineModelName string     `json:"machineModelName"`
	Gpu              int        `json:"gpu"`
	Cpu              int        `json:"cpu"`
	Ram              int        `json:"ram"`
	Disk             int        `json:"disk"`
	Eni              int        `json:"eni"`
	Topologies       []Topology `json:"topologies"`
}

type Topology struct {
	Socket int `json:"socket"`
	Core   int `json:"core"`
	Cpu    int `json:"cpu"`
}

type Pod struct {
	AppName string `json:"appName"`
	Group   string `json:"group"`
	Gpu     int    `json:"gpu"`
	Cpu     int    `json:"cpu"`
	Ram     int    `json:"ram"`
	Disk    int    `json:"disk"`
	CpuIds  []int  `json:"cpuIds"`
}

type NodeWithPod struct {
	Node Node  `json:"node"`
	Pods []Pod `json:"pods"`
}
