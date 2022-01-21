package model

import (
	"errors"
	"fmt"
	"teamide/application/base"
)

func getServiceStepRedisSetByMap(value interface{}) (step ServiceStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["redisSet"] == nil {
			return
		}
		switch data := v["redisSet"].(type) {
		case map[string]interface{}:
			v["redisSet"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to redis set error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to service step redis set error"))
	}
	if err != nil {
		return
	}
	step = &ServiceStepRedisSet{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func getServiceStepRedisGetByMap(value interface{}) (step ServiceStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["redisGet"] == nil {
			return
		}
		switch data := v["redisGet"].(type) {
		case map[string]interface{}:
			v["redisGet"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to redis get error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to service step redis get error"))
	}
	if err != nil {
		return
	}
	step = &ServiceStepRedisGet{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func getServiceStepRedisDelByMap(value interface{}) (step ServiceStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["redisDel"] == nil {
			return
		}
		switch data := v["redisDel"].(type) {
		case map[string]interface{}:
			v["redisDel"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to redis del error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to service step redis del error"))
	}
	if err != nil {
		return
	}
	step = &ServiceStepRedisDel{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func getServiceStepRedisExpireByMap(value interface{}) (step ServiceStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["redisExpire"] == nil {
			return
		}
		switch data := v["redisExpire"].(type) {
		case map[string]interface{}:
			v["redisExpire"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to redis expire error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to service step redis expire error"))
	}
	if err != nil {
		return
	}
	step = &ServiceStepRedisExpire{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func getServiceStepRedisExpireatByMap(value interface{}) (step ServiceStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["redisExpireat"] == nil {
			return
		}
		switch data := v["redisExpireat"].(type) {
		case map[string]interface{}:
			v["redisExpireat"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to redis expireat error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to service step redis expireat error"))
	}
	if err != nil {
		return
	}
	step = &ServiceStepRedisExpireat{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}
