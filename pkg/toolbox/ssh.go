package toolbox

import (
	"encoding/json"
	"teamide/pkg/application/base"
)

func init() {
	worker_ := &Worker{
		Name: "ssh",
		Text: "SSH",
		Work: sshWork,
	}

	AddWorker(worker_)
}

type SSHConfig struct {
	Type     string `json:"type"`
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type SSHBaseRequest struct {
}

var (
	sshTokenCache = map[string]*SSHConfig{}
)

func sshWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

	var sshConfig *SSHConfig
	var bs []byte
	bs, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &sshConfig)
	if err != nil {
		return
	}

	bs, err = json.Marshal(data)
	if err != nil {
		return
	}
	request := &SSHBaseRequest{}
	err = json.Unmarshal(bs, request)
	if err != nil {
		return
	}

	res = map[string]interface{}{}
	switch work {
	case "createToken":
		var token string = base.GenerateUUID()
		sshTokenCache[token] = sshConfig
		res["token"] = token
	}
	return
}
