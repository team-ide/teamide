package redis

import (
	"server/base"
)

type getRequest struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
	Key     string `json:"key"`
}

type getResponse struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type ValueInfo struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func getWork(req interface{}) (res interface{}, err error) {
	request := &getRequest{}
	response := &getResponse{}
	err = base.ToBean(req.([]byte), request)
	if err != nil {
		return
	}
	var service Service
	service, err = getService(request.Address, request.Auth)
	if err != nil {
		return
	}
	var valueInfo ValueInfo
	valueInfo, err = service.Get(request.Key)
	if err != nil {
		return
	}
	response.Type = valueInfo.Type
	response.Value = valueInfo.Value
	res = response
	return
}

type keysRequest struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
	Pattern string `json:"pattern"`
	Size    int    `json:"size"`
}

type keysResponse struct {
	Keys  []string `json:"keys"`
	Count int      `json:"count"`
}

func keysWork(req interface{}) (res interface{}, err error) {
	request := &keysRequest{}
	response := &keysResponse{}
	err = base.ToBean(req.([]byte), request)
	if err != nil {
		return
	}
	var service Service
	service, err = getService(request.Address, request.Auth)
	if err != nil {
		return
	}
	var count int
	var keys []string
	count, keys, err = service.Keys(request.Pattern, request.Size)
	if err != nil {
		return
	}
	response.Keys = keys
	response.Count = count
	res = response
	return
}

type doRequest struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Index   int64  `json:"index"`
	Count   int64  `json:"count"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}

type doResponse struct {
}

func doWork(req interface{}) (res interface{}, err error) {
	request := &doRequest{}
	response := &doResponse{}
	err = base.ToBean(req.([]byte), request)
	if err != nil {
		return
	}
	var service Service
	service, err = getService(request.Address, request.Auth)
	if err != nil {
		return
	}
	if request.Type == "set" {
		err = service.Set(request.Key, request.Value)
	} else if request.Type == "sadd" {
		err = service.Sadd(request.Key, request.Value)
	} else if request.Type == "srem" {
		err = service.Srem(request.Key, request.Value)
	} else if request.Type == "lpush" {
		err = service.Lpush(request.Key, request.Value)
	} else if request.Type == "rpush" {
		err = service.Rpush(request.Key, request.Value)
	} else if request.Type == "lset" {
		err = service.Lset(request.Key, request.Index, request.Value)
	} else if request.Type == "lrem" {
		err = service.Lrem(request.Key, request.Count, request.Value)
	} else if request.Type == "hset" {
		err = service.Hset(request.Key, request.Field, request.Value)
	} else if request.Type == "hdel" {
		err = service.Hdel(request.Key, request.Field)
	}
	if err != nil {
		return
	}
	res = response
	return
}

type deleteRequest struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
	Key     string `json:"key"`
}

type deleteResponse struct {
	Count int `json:"count"`
}

func deleteWork(req interface{}) (res interface{}, err error) {
	request := &deleteRequest{}
	response := &deleteResponse{}
	err = base.ToBean(req.([]byte), request)
	if err != nil {
		return
	}
	var service Service
	service, err = getService(request.Address, request.Auth)
	if err != nil {
		return
	}
	var count int
	count, err = service.Del(request.Key)
	response.Count = count
	res = response
	if err != nil {
		return
	}
	return
}

type deletePatternRequest struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
	Pattern string `json:"pattern"`
}

type deletePatternResponse struct {
	Count int `json:"count"`
}

func deletePatternWork(req interface{}) (res interface{}, err error) {
	request := &deletePatternRequest{}
	response := &deletePatternResponse{}
	err = base.ToBean(req.([]byte), request)
	if err != nil {
		return
	}
	var service Service
	service, err = getService(request.Address, request.Auth)
	if err != nil {
		return
	}
	var count int
	count, err = service.DelPattern(request.Pattern)
	response.Count = count
	res = response
	if err != nil {
		return
	}
	return
}
