package toolbox

import (
	"encoding/json"
)

func init() {
	worker_ := &Worker{
		Name: "zookeeper",
		Text: "Zookeeper",
		Work: zkWork,
	}

	AddWorker(worker_)
}

type ZookeeperBaseRequest struct {
	Path string `json:"path"`
	Data string `json:"data"`
}

func zkWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {
	var service *ZKService
	service, err = getZKService(config["address"].(string))
	if err != nil {
		return
	}

	var bs []byte
	bs, err = json.Marshal(data)
	if err != nil {
		return
	}
	request := &ZookeeperBaseRequest{}
	err = json.Unmarshal(bs, request)
	if err != nil {
		return
	}

	res = map[string]interface{}{}
	switch work {
	case "get":
		var data []byte
		data, err = service.Get(request.Path)
		if err != nil {
			return
		}
		res["data"] = string(data)
	case "save":
		var isEx bool
		isEx, err = service.Exists(request.Path)
		if err != nil {
			return
		}
		if isEx {
			err = service.SetData(request.Path, []byte(request.Data))
		} else {
			err = service.CreateIfNotExists(request.Path, []byte(request.Data))
		}
		if err != nil {
			return
		}
	case "getChildren":
		var isEx bool
		isEx, err = service.Exists(request.Path)
		if err != nil {
			return
		}
		if isEx {
			var children []string
			children, err = service.GetChildren(request.Path)
			if err != nil {
				return
			}
			res["children"] = children
		}
	case "delete":
		var isEx bool
		isEx, err = service.Exists(request.Path)
		if err != nil {
			return
		}
		if isEx {
			err = service.Delete(request.Path)
			if err != nil {
				return
			}
		}
	}
	return
}
