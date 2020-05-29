package loader

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"django-go/pkg/types"
)

const (
	appDataFileName         = "schedule.app.source"
	nodeDataFileName        = "schedule.node.source"
	nodeWithPodDataFileName = "reschedule.source"
	ruleDataFileName        = "rule.source"
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
	return ioutil.ReadFile(filepath.Join(dataBaseDir, loader.dir, fileName))
}

func (loader Loader) LoadApps() ([]types.App, error) {
	apps := make([]types.App, 0)
	data, err := loader.loadData(apps, appDataFileName)
	if err == nil {
		err = json.Unmarshal(data, &apps)
		return apps, err
	}
	return nil, err
}

func (loader Loader) LoadNodes() ([]types.Node, error) {
	nodes := make([]types.Node, 0)
	data, err := loader.loadData(nodes, nodeDataFileName)
	if err == nil {
		err = json.Unmarshal(data, &nodes)
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
		return rule, err
	}
	return types.Rule{}, err
}
