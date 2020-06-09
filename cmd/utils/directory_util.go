package utils

import (
	"strings"
	"fmt"
	"django-go/pkg/util"
	"django-go/pkg/constants"
)

/**
 * 根据args中的第一个参数或者默认目录(args为空)，确定多少个目录下的数据同时用于并行计算。
 */
func AdjustDirectorys(args []string) []string {

	var dirs []string

	if len(args) == 0 {
		dirs = []string{constants.TEST_DIRECTORY}
	} else {
		dirs = removeDuplicateElement(strings.Split(args[0], ","))
	}

	fmt.Println(fmt.Sprintf("running data directory " + util.ToJsonOrDie(dirs)))

	return dirs
}

func removeDuplicateElement(addrs []string) []string { //去掉重复元素

	result := make([]string, 0, len(addrs))

	temp := map[string]struct{}{}

	for _, item := range addrs {

		if _, ok := temp[item]; !ok {

			temp[item] = struct{}{}

			result = append(result, item)
		}
	}

	return result
}
