package toolbox

import (
	"encoding/json"
)

type OtherConfig struct {
	Type string `json:"type"`
}

type OtherRequest struct {
}

func OtherWork(work string, config *OtherConfig, data map[string]interface{}) (res map[string]interface{}, err error) {

	dataBS, err := json.Marshal(data)
	if err != nil {
		return
	}
	request := &OtherRequest{}
	err = json.Unmarshal(dataBS, request)
	if err != nil {
		return
	}

	res = map[string]interface{}{}
	switch work {
	case "export":
		res["request"] = request
	}
	return
}
