package toolbox

import (
	"encoding/json"
	"teamide/pkg/form"
	"teamide/pkg/zookeeper"
)

func zookeeperWorker() *Worker {
	worker_ := &Worker{
		Name: "zookeeper",
		Text: "Zookeeper",
		Work: zkWork,
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "连接地址（127.0.0.1:2181）", Name: "address", DefaultValue: "127.0.0.1:2181",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
			},
		},
	}

	return worker_
}

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

func zkWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

	var zkConfig zookeeper.Config
	var configBS []byte
	configBS, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(configBS, &zkConfig)
	if err != nil {
		return
	}

	var service *zookeeper.ZKService
	service, err = getZKService(zkConfig)
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

			var children []map[string]interface{}

			var names []string
			names, err = service.GetChildren(request.Path)
			if err != nil {
				return
			}
			for _, name := range names {
				var one = map[string]interface{}{}
				one["name"] = name

				//var nameSubs []string
				//nameSubs, _ = service.GetChildren(request.Path + "/" + name)
				//one["hasChildren"] = len(nameSubs) > 0

				children = append(children, one)
			}
			res["children"] = children
		}
	case "hasChildren":
		var isEx bool
		isEx, err = service.Exists(request.Path)
		if err != nil {
			return
		}
		if isEx {

			var names []string
			names, err = service.GetChildren(request.Path)
			if err != nil {
				return
			}
			res["hasChildren"] = len(names) > 0
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
