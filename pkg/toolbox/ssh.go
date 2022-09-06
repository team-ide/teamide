package toolbox

import (
	"encoding/json"
	"teamide/pkg/ssh"
	"teamide/pkg/util"
)

type SSHBaseRequest struct {
	Token    string `json:"token"`
	WorkerId string `json:"workerId,omitempty"`
}

func SSHWork(work string, config *ssh.Config, data map[string]interface{}) (res map[string]interface{}, err error) {

	dataBS, err := json.Marshal(data)
	if err != nil {
		return
	}
	request := &SSHBaseRequest{}
	err = json.Unmarshal(dataBS, request)
	if err != nil {
		return
	}

	res = map[string]interface{}{}
	switch work {
	case "createToken":
		var token = util.UUID()
		ssh.AddTokenCache(token, config)
		res["token"] = token
	case "close":
	}
	return
}
