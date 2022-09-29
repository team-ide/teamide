package model

import (
	"encoding/json"
	"teamide/pkg/util"
)

func toList(value interface{}, newOne func(str string) interface{}, addOne func(one interface{})) {
	if value == nil {
		return
	}
	list, listOk := value.([]interface{})
	if listOk {
		for _, v := range list {
			one := toOne(v, newOne)
			addOne(one)
		}
	} else {
		one := toOne(value, newOne)
		addOne(one)
	}
	return
}

func toOne(value interface{}, newOne func(str string) interface{}) (res interface{}) {
	vMap, vMapOk := value.(map[string]interface{})

	if vMapOk {
		res = newOne("")
		bs, _ := json.Marshal(vMap)
		_ = json.Unmarshal(bs, res)

	} else {
		str, _ := util.GetStringValue(value)
		res = newOne(str)
	}
	return
}
