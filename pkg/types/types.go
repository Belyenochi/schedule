package types

type App struct {
	AppName      string `json:"appName"`
	Group        string `json:"group"`
	CpuSensitive bool   `json:"cpuSensitive"`
	Gpu          int    `json:"gpu"`
	Cpu          int    `json:"cpu"`
	Ram          int    `json:"ram"`
	Disk         int    `json:"disk"`
	Replicas     int    `json:"replicas"`
}

type Node struct {
	Sn            string     `json:"sn"`
	NodeModelName string     `json:"nodeModelName"`
	Gpu           int        `json:"gpu"`
	Cpu           int        `json:"cpu"`
	Ram           int        `json:"ram"`
	Disk          int        `json:"disk"`
	Eni           int        `json:"eni"`
	Topologies    []Topology `json:"topologies"`
}

type Topology struct {
	Socket int `json:"socket"`
	Core   int `json:"core"`
	Cpu    int `json:"cpu"`
}

type Pod struct {
	PodSn   string `json:"podSn"`
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

func (app App) Copy() App {
	return App{
		AppName:      app.AppName,
		Group:        app.Group,
		CpuSensitive: app.CpuSensitive,
		Gpu:          app.Gpu,
		Cpu:          app.Cpu,
		Ram:          app.Ram,
		Disk:         app.Disk,
		Replicas:     app.Replicas,
	}
}

func (node Node) Copy() Node {
	return Node{
		Sn:            node.Sn,
		NodeModelName: node.NodeModelName,
		Gpu:           node.Gpu,
		Cpu:           node.Cpu,
		Ram:           node.Ram,
		Disk:          node.Disk,
		Eni:           node.Eni,
		Topologies:    copyTopologies(node.Topologies),
	}
}

func (topology Topology) Copy() Topology {
	return Topology{
		Socket: topology.Socket,
		Core:   topology.Core,
		Cpu:    topology.Cpu,
	}
}

func (pod Pod) Copy() Pod {
	return Pod{
		PodSn:   pod.PodSn,
		AppName: pod.AppName,
		Group:   pod.Group,
		Gpu:     pod.Gpu,
		Cpu:     pod.Cpu,
		Ram:     pod.Ram,
		Disk:    pod.Disk,
		CpuIds:  pod.CpuIds,
	}
}

func (nwp NodeWithPod) Copy() NodeWithPod {
	return NodeWithPod{
		Node: nwp.Node.Copy(),
		Pods: copyPods(nwp.Pods),
	}
}

func CopyApps(apps []App) []App {

	result := make([]App, 0)

	for _, app := range apps {
		result = append(result, app.Copy())
	}

	return result
}

func CopyNodes(nodes []Node) []Node {

	result := make([]Node, 0)

	for _, node := range nodes {
		result = append(result, node.Copy())
	}

	return result
}

func CopyNodeWithPods(nwps []NodeWithPod) []NodeWithPod {

	result := make([]NodeWithPod, 0)

	for _, nwp := range nwps {
		result = append(result, nwp.Copy())
	}

	return result
}


func copyPods(pods []Pod) []Pod {

	result := make([]Pod, 0)

	for _, pod := range pods {
		result = append(result, pod.Copy())
	}

	return result
}

func copyTopologies(topologies []Topology) []Topology {

	result := make([]Topology, 0)

	for _, topology := range topologies {
		result = append(result, topology.Copy())
	}

	return result
}

func (node Node) Value(resource Resource) int {
	if resource == GPU {
		return node.Gpu
	} else if resource == CPU {
		return node.Cpu
	} else if resource == RAM {
		return node.Ram
	} else {
		return node.Disk
	}
}

func (pod Pod) Value(resource Resource) int {
	if resource == GPU {
		return pod.Gpu
	} else if resource == CPU {
		return pod.Cpu
	} else if resource == RAM {
		return pod.Ram
	} else {
		return pod.Disk
	}
}
