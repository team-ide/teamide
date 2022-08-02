package toolbox

import (
	"encoding/json"
	"teamide/pkg/zookeeper"
)

func getZKService(zkConfig zookeeper.Config) (res *zookeeper.ZKService, err error) {
	key := "zookeeper-" + zkConfig.Address
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s *zookeeper.ZKService
		s, err = zookeeper.CreateZKService(zkConfig)
		if err != nil {
			return
		}
		_, err = s.Exists("/")
		if err != nil {
			return
		}
		res = s
		return
	})
	if err != nil {
		return
	}
	res = service.(*zookeeper.ZKService)
	return
}

type ZookeeperBaseRequest struct {
	Path string `json:"path"`
	Data string `json:"data"`
}

func ZKWork(work string, config *zookeeper.Config, data map[string]interface{}) (res map[string]interface{}, err error) {

	var service *zookeeper.ZKService
	service, err = getZKService(*config)
	if err != nil {
		return
	}

	dataBS, err := json.Marshal(data)
	if err != nil {
		return
	}
	request := &ZookeeperBaseRequest{}
	err = json.Unmarshal(dataBS, request)
	if err != nil {
		return
	}

	res = map[string]interface{}{}
	switch work {
	case "get":
		var data []byte
		var statInfo *zookeeper.StatInfo
		data, statInfo, err = service.Get(request.Path)
		if err != nil {
			return
		}
		res["data"] = string(data)
		res["stat"] = statInfo
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

			var children []map[string]interface{}

			var names []string
			names, err = service.GetChildren(request.Path)
			if err != nil {
				return
			}
			for _, name := range names {
				var one = map[string]interface{}{}
				one["name"] = name

				childrenPath := "/" + name
				if request.Path != "/" {
					childrenPath = request.Path + childrenPath
				}
				var statInfo *zookeeper.StatInfo
				statInfo, err = service.Stat(childrenPath)
				if err != nil {
					return
				}
				one["hasChildren"] = statInfo.NumChildren > 0

				children = append(children, one)
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
