package toolbox

import (
	"encoding/json"
	"teamide/pkg/form"
)

func otherWorker() *Worker {
	worker_ := &Worker{
		Name: "other",
		Text: "其它",
		Work: otherWork,
		ConfigForm: &form.Form{
			Fields: []*form.Field{},
		},
	}

	return worker_
}

type OtherRequest struct {
}

func otherWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

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
