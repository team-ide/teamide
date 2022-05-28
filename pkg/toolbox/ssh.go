package toolbox

import (
	"encoding/json"
	"errors"
	"teamide/pkg/form"
	"teamide/pkg/ssh"
	"teamide/pkg/util"
)

func sshWorker() *Worker {
	worker_ := &Worker{
		Name: "ssh",
		Text: "SSH",
		Work: sshWork,
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "类型", Name: "type", Type: "select", DefaultValue: "tcp",
					Options: []*form.Option{
						{Text: "TCP", Value: "tcp"},
					},
					Rules: []*form.Rule{
						{Required: true, Message: "SSH类型不能为空"},
					},
				},
				{
					Label: "连接地址（127.0.0.1:22）", Name: "address", DefaultValue: "127.0.0.1:22",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
				{Label: "Username", Name: "username"},
				{Label: "Password", Name: "password", Type: "password"},
				{Label: "PublicKey", Name: "publicKey", Type: "file", Placeholder: "请上传PublicKey文件"},
			},
		},
	}

	return worker_
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
	*ssh.SFTPRequest
}

func sshWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

	var sshConfig *ssh.Config
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
		var token = util.GenerateUUID()
		ssh.TokenCache[token] = sshConfig
		res["token"] = token
	case "readText":
		var token = request.Token
		client := ssh.SftpCache[token]
		if client == nil {
			err = errors.New("FTP会话丢失")
			return
		}
		var response *ssh.SFTPResponse
		if request.Place == "local" {
			response, err = client.LocalReadText(request.SFTPRequest)
		} else if request.Place == "remote" {
			response, err = client.RemoteReadText(request.SFTPRequest)
		}
		if err != nil {
			return
		}
		res["response"] = response
	case "saveText":
		var token = request.Token
		client := ssh.SftpCache[token]
		if client == nil {
			err = errors.New("FTP会话丢失")
			return
		}
		var response *ssh.SFTPResponse
		if request.Place == "local" {
			response, err = client.LocalSaveText(request.SFTPRequest)
		} else if request.Place == "remote" {
			response, err = client.RemoteSaveText(request.SFTPRequest)
		}
		if err != nil {
			return
		}
		res["response"] = response
	}
	return
}
