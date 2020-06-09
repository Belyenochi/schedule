package util

import (
	"django-go/pkg/types"
	"django-go/pkg/constants"
	"fmt"
	"errors"
)

/**
 * rule -> nodeResourceWeights数据转化为<资源,<机型,权重>>。
 */
func toResourceScoreMap(rule types.Rule) map[types.Resource]map[string]int {

	resourceMap := make(map[types.Resource]map[string]int)

	for _, rw := range rule.NodeResourceWeights {

		mmn := String(rw.NodeModelName).ValueWithDefault(constants.SCORE_EMPTY_NODE_MODEL_NAME)

		if _, ok := resourceMap[rw.Resource]; ok {
			resourceMap[rw.Resource][mmn] = rw.Weight
		} else {
			resourceMap[rw.Resource] = map[string]int{mmn: rw.Weight}
		}
	}

	return resourceMap
}

/**
 * 返回所有应用<分组,分组对应的最大可堆叠实例数>
 */
func ToAllMaxInstancePerNodeLimit(rule types.Rule, groupRuleAssociates []types.GroupRuleAssociate) map[string]int {

	allMaxInstancePerNodeLimit := make(map[string]int)

	for _, replicasMaxInstancePerNode := range rule.ReplicasMaxInstancePerNodes {

		restrain := replicasMaxInstancePerNode.Restrain

		if !StringSlice(types.AllRestrains).Contain(string(restrain)) {
			fmt.Println(errors.New("ReplicasMaxInstancePerNode rule error"))
			continue
		}

		replicas := replicasMaxInstancePerNode.Replicas

		maxInstancePerHost := replicasMaxInstancePerNode.MaxInstancePerNode

		if restrain == types.GE {//此规则用意为应用实例数大于或等于(GE)replicas的应用，最大堆叠数量为maxInstancePerNode。
			for _, groupRuleAssociate := range groupRuleAssociates {
				if groupRuleAssociate.Replicas >= replicas {
					allMaxInstancePerNodeLimit[groupRuleAssociate.Group] = maxInstancePerHost
				}
			}
			continue
		}

		if restrain == types.LE {
			for _, groupRuleAssociate := range groupRuleAssociates {
				if groupRuleAssociate.Replicas <= replicas {
					if _, ok := allMaxInstancePerNodeLimit[groupRuleAssociate.Group]; !ok {
						allMaxInstancePerNodeLimit[groupRuleAssociate.Group] = maxInstancePerHost
						continue
					}
					allMaxInstancePerNodeLimit[groupRuleAssociate.Group] = Min(allMaxInstancePerNodeLimit[groupRuleAssociate.Group], maxInstancePerHost)
				}
			}
		}

	}

	for _, groupRuleAssociate := range groupRuleAssociates {
		if _, ok := allMaxInstancePerNodeLimit[groupRuleAssociate.Group]; !ok {
			allMaxInstancePerNodeLimit[groupRuleAssociate.Group] = rule.DefaultMaxInstancePerNode
		}
	}

	return allMaxInstancePerNodeLimit
}
