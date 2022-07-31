package toolbox

import (
	"encoding/json"
	"sync"
	"teamide/pkg/ssh"
	"teamide/pkg/util"
)

type SSHBaseRequest struct {
	Token string `json:"token"`
	*ssh.SFTPRequest
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
		ssh.TokenCache[token] = config
		res["token"] = token
	case "readText":
		client := createOrGetSftpClient(request.WorkerId, config)
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
		client := createOrGetSftpClient(request.WorkerId, config)
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
	case "ftpWork":
		client := createOrGetSftpClient(request.WorkerId, config)
		var response *ssh.SFTPResponse
		response, err = client.Work(request.SFTPRequest)
		if err != nil {
			return
		}
		res["response"] = response
	case "close":
		closeSftpClient(request.WorkerId)
	}
	return
}

var (
	sftpCache     = make(map[string]*ssh.SftpClient)
	sftpCacheLock sync.Mutex
)

func createOrGetSftpClient(workerId string, config *ssh.Config) (res *ssh.SftpClient) {
	sftpCacheLock.Lock()
	defer sftpCacheLock.Unlock()
	res, ok := sftpCache[workerId]
	if !ok {
		res = &ssh.SftpClient{
			WorkerId: workerId,
			Config:   *config,
		}
		res.Start()
		sftpCache[workerId] = res
	}
	return
}
func GetSftpClient(workerId string) (res *ssh.SftpClient) {
	sftpCacheLock.Lock()
	defer sftpCacheLock.Unlock()
	res = sftpCache[workerId]
	return
}

func closeSftpClient(workerId string) (res *ssh.SftpClient) {
	sftpCacheLock.Lock()
	defer sftpCacheLock.Unlock()
	res, ok := sftpCache[workerId]
	if ok {
		delete(sftpCache, workerId)
		res.Close()
	}
	return
}
