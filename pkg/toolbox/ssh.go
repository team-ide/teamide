package toolbox

import (
	"encoding/json"
	"teamide/pkg/application/base"
	"teamide/pkg/form"
)

func init() {
	worker_ := &Worker{
		Name: "ssh",
		Text: "SSH",
		Work: sshWork,
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{Label: "类型", Name: "type", Type: "select", DefaultValue: "tcp",
					Options: []*form.Option{
						{Text: "TCP", Value: "tcp"},
					},
					Rules: []*form.Rule{
						{Required: true, Message: "SSH类型不能为空"},
					},
				},
				{Label: "连接地址（127.0.0.1:22）", Name: "host", DefaultValue: "127.0.0.1:22",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "Username", Name: "username"},
				{Label: "Password", Name: "password"},
				{Label: "PublicKey", Name: "publicKey", Type: "file", Placeholder: "请上传PublicKey文件"},
			},
		},
	}

	AddWorker(worker_)
}

type SSHConfig struct {
	Type      string `json:"type"`
	Address   string `json:"address"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	PublicKey string `json:"publicKey"`
}

type SSHBaseRequest struct {
	Token string `json:"token"`
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
