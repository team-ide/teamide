package toolbox

import (
	"encoding/json"
	"errors"
	"teamide/pkg/ssh"
	"teamide/pkg/util"
)

type SSHBaseRequest struct {
	Token string `json:"token"`
	*ssh.SFTPRequest
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
		ssh.TokenCache[token] = config
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
