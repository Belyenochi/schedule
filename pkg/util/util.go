package util

import "encoding/json"

type String string

func (s String) ValueWithDefault(def string) string {
	if len(s) == 0 {
		return def
	}
	return string(s)
}

type StringSlice []string

func (ss StringSlice) Contain(key string) bool {
	for _, s := range ss {
		if s == key {
			return true
		}
	}
	return false
}

func ToJsonOrDie(obj interface{}) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func MustBeTrue(result bool, msg string) {
	if !result {
		panic(msg)
	}
}
