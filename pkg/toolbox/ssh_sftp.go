package toolbox

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"io"
)

type SSHSftpClient struct {
	SSHClient
	isClosedSftp bool
	sftpClient   *sftp.Client
	UploadFile   chan *UploadFile
	confirmMap   map[string]chan *ConfirmInfo
}

func (this_ *SSHSftpClient) CloseSftp() {

	this_.isClosedSftp = true
	if this_.sftpClient != nil {
		err := this_.sftpClient.Close()
		if err != nil {
			this_.Logger.Error("SSH SFTP close error", zap.Error(err))
		}
	}
	this_.sftpClient = nil
	this_.CloseClient()
}

func (this_ *SSHSftpClient) initSftp() (err error) {
	if this_.isClosedSftp || this_.sftpClient == nil {
		err = this_.createSftp()
	}
	return
}

func (this_ *SSHSftpClient) createSftp() (err error) {
	SSHSftpCache[this_.Token] = this_
	err = this_.initClient()
	if err != nil {
		return
	}

	this_.sftpClient, err = sftp.NewClient(this_.sshClient)
	if err != nil {
		this_.WSWriteError("SSH FTP创建失败:" + err.Error())
		this_.CloseSftp()
		return
	}

	return
}

func (this_ *SSHSftpClient) start() (err error) {
	err = this_.initSftp()
	if err != nil {
		return
	}
	if this_.UploadFile == nil {
		this_.UploadFile = make(chan *UploadFile, 10)
	}
	go func() {
		for {
			select {
			case uploadFile := <-this_.UploadFile:
				this_.work(&SFTPRequest{
					Work:     "upload",
					WorkId:   uploadFile.WorkId,
					Dir:      uploadFile.Dir,
					Place:    uploadFile.Place,
					File:     uploadFile.File,
					FullPath: uploadFile.FullPath,
				})
			}
		}

	}()
	// 第一个协程获取用户的输入
	go func() {
		for {
			if this_.isClosedSftp {
				return
			}
			_, p, err := this_.ws.ReadMessage()
			if err != nil && err != io.EOF {
				fmt.Println("sftp ws read err:", err)
				this_.CloseSftp()
				return
			}
			//fmt.Println("sftp ws read:" + string(p))
			if len(p) > 0 {
				if this_.isClosedSftp {
					return
				}

				go func() {
					var request *SFTPRequest
					err = json.Unmarshal(p, &request)
					if err != nil {
						fmt.Println("sftp ws message to struct err:", err)
						return
					}
					this_.work(request)
				}()
			}
		}
	}()
	return
}
