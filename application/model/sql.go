package model

import (
	"errors"
	"fmt"
	"teamide/application/base"
)

func getColumnMapsByValue(value interface{}) (columns []map[string]interface{}, err error) {
	if value == nil {
		return
	}
	values, valuesOk := value.([]interface{})
	if !valuesOk {
		return
	}
	for _, valuesOne := range values {

		switch v := valuesOne.(type) {
		case map[string]interface{}:
			newV := v
			if len(v) == 1 {
				for mapKey, mapValue := range v {
					if mapValue == nil {
						continue
					}
					newV = map[string]interface{}{}
					newV["name"] = mapKey
					vMap, vMapOk := mapValue.(map[string]interface{})
					if vMapOk {
						base.ToBean([]byte(base.ToJSON(vMap)), &newV)
					} else {
						newV["value"] = fmt.Sprint(mapValue)
					}
				}
			}
			// _, find := newV["ignoreEmpty"]
			// if !find {
			// 	newV["ignoreEmpty"] = true
			// }
			columns = append(columns, newV)
		case string:
			columns = append(columns, map[string]interface{}{
				"name":        v,
				"ignoreEmpty": true,
			})
		default:
			err = errors.New(fmt.Sprint("[", v, "] to column error"))
		}
	}
	return
}

func getWhereMapsByValue(value interface{}) (wheres []map[string]interface{}, err error) {
	if value == nil {
		return
	}
	values, valuesOk := value.([]interface{})
	if !valuesOk {
		return
	}
	for _, valuesOne := range values {
		switch v := valuesOne.(type) {
		case map[string]interface{}:
			wheres = append(wheres, v)
		case string:
			wheres = append(wheres, map[string]interface{}{
				"custom":    true,
				"customSql": v,
			})
		default:
			err = errors.New(fmt.Sprint("[", v, "] to where error"))
		}
	}
	return
}

func getActionStepSqlSelectByMap(value interface{}) (step ActionStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["sqlSelect"] == nil {
			return
		}
		switch data := v["sqlSelect"].(type) {
		case map[string]interface{}:
			data["columns"], err = getColumnMapsByValue(data["columns"])
			if err != nil {
				return
			}
			data["wheres"], err = getWhereMapsByValue(data["wheres"])
			v["sqlSelect"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to sql select error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to action step sql select error"))
	}
	if err != nil {
		return
	}
	step = &ActionStepSqlSelect{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func getActionStepSqlInsertByMap(value interface{}) (step ActionStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["sqlInsert"] == nil {
			return
		}
		switch data := v["sqlInsert"].(type) {
		case map[string]interface{}:
			data["columns"], err = getColumnMapsByValue(data["columns"])
			v["sqlInsert"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to sql insert error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to action step sql insert error"))
	}
	if err != nil {
		return
	}
	step = &ActionStepSqlInsert{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func getActionStepSqlUpdateByMap(value interface{}) (step ActionStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["sqlUpdate"] == nil {
			return
		}
		switch data := v["sqlUpdate"].(type) {
		case map[string]interface{}:
			data["columns"], err = getColumnMapsByValue(data["columns"])
			if err != nil {
				return
			}
			data["wheres"], err = getWhereMapsByValue(data["wheres"])
			v["sqlUpdate"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to sql update error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to action step sql update error"))
	}
	if err != nil {
		return
	}
	step = &ActionStepSqlUpdate{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func getActionStepSqlDeleteByMap(value interface{}) (step ActionStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["sqlDelete"] == nil {
			return
		}
		switch data := v["sqlDelete"].(type) {
		case map[string]interface{}:
			data["wheres"], err = getWhereMapsByValue(data["wheres"])
			v["sqlDelete"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to sql delete error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to action step sql delete error"))
	}
	if err != nil {
		return
	}
	step = &ActionStepSqlDelete{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}
