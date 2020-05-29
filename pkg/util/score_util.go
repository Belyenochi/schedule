package util

import "django-go/pkg/types"

const (
	scoreEmptyName = "SOCRE_EMPTY_SMNAME"
)

func ScoreNodes(nodes []types.Node, rule types.Rule) int {
	scoreSum := 0

	scoreMap := toScoreMap(rule)

	for _, node := range nodes {
		scoreSum += score(scoreMap, node)
	}
	return scoreSum
}

func ScoreNode(node types.Node, rule types.Rule) int {
	return ScoreNodes([]types.Node{node}, rule)
}

func ScoreNodeWithPods(nodeWithPods []types.NodeWithPod, rule types.Rule) int {
	nodes := make([]types.Node, 0, len(nodeWithPods))
	for _, nwp := range nodeWithPods {
		nodes = append(nodes, nwp.Node)
	}

	return ScoreNodes(nodes, rule)
}

func score(scoreMap map[types.Resource]map[string]int, node types.Node) int {
	weightSum := 0
	mmn := String(node.MachineModelName).ValueWithDefault(scoreEmptyName)
	for _, r := range types.AllResources {
		mw, ok := scoreMap[r]
		if !ok {
			continue
		}
		if w, ok := mw[mmn]; ok {
			weightSum += w * nodeResourceValue(node, r)
		}
	}
	return weightSum
}

//rule -> nodeResourceWeights数据转化为<资源,<机型,权重>>。
func toScoreMap(rule types.Rule) map[types.Resource]map[string]int {
	resourceMap := make(map[types.Resource]map[string]int)
	for _, rw := range rule.NodeResourceWeights {
		mmn := String(rw.MachineModelName).ValueWithDefault(scoreEmptyName)
		if _, ok := resourceMap[rw.Resource]; ok {
			resourceMap[rw.Resource][mmn] = rw.Weight
		} else {
			resourceMap[rw.Resource] = map[string]int{mmn: rw.Weight}
		}
	}

	return resourceMap
}
