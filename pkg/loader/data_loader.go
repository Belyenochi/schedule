package loader

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"django-go/pkg/types"
	"fmt"
	"django-go/pkg/util"
	"django-go/pkg/constants"
)

const (
	appDataFileName         = constants.SCHEDULE_APP_SOURCE
	nodeDataFileName        = constants.SCHEDULE_NODE_SOURCE
	nodeWithPodDataFileName = constants.RESCHEDULE_SOURCE
	ruleDataFileName        = constants.RULE_SOURCE
)

func NewLoader(dir string) Loader {
	return Loader{dir}
}

type Loader struct {
	dir string
}

func (loader Loader) loadData(target interface{}, fileName string) ([]byte, error) {

	_, currentFilePath, _, _ := runtime.Caller(1)

	dataBaseDir := strings.Replace(currentFilePath, "pkg/loader/data_loader.go", "data", 1)

	fmt.Println(fmt.Sprintf("read file from directory:%v , fileName:%v", loader.dir, fileName))

	return ioutil.ReadFile(filepath.Join(dataBaseDir, loader.dir, fileName))
}

func (loader Loader) LoadApps() ([]types.App, error) {

	apps := make([]types.App, 0)

	data, err := loader.loadData(apps, appDataFileName)

	if err == nil {

		err = json.Unmarshal(data, &apps)

		groupCount := 0

		totalReplicas := 0

		for _, app := range apps {
			groupCount++
			totalReplicas += app.Replicas
		}

		fmt.Println(fmt.Sprintf("%v | source count :【group %v, pod %v】", loader.dir, groupCount, totalReplicas))

		return apps, err
	}

	return nil, err
}

func (loader Loader) LoadNodes() ([]types.Node, error) {
	nodes := make([]types.Node, 0)
	data, err := loader.loadData(nodes, nodeDataFileName)
	if err == nil {
		err = json.Unmarshal(data, &nodes)
		fmt.Println(fmt.Sprintf("%v | source node count : %v", loader.dir, len(nodes)))
		return nodes, err
	}
	return nil, err
}

func (loader Loader) LoadNodeWithPods() ([]types.NodeWithPod, error) {
	nodeWithPods := make([]types.NodeWithPod, 0)
	data, err := loader.loadData(nodeWithPods, nodeWithPodDataFileName)
	if err == nil {
		err = json.Unmarshal(data, &nodeWithPods)
		return nodeWithPods, err
	}
	return nil, err
}

func (loader Loader) LoadRule() (types.Rule, error) {
	rule := types.Rule{}
	data, err := loader.loadData(rule, ruleDataFileName)
	if err == nil {
		err = json.Unmarshal(data, &rule)
		fmt.Println(fmt.Sprintf("%v | rule : %v", loader.dir, util.ToJsonOrDie(rule)))
		return rule, err
	}
	return types.Rule{}, err
}
