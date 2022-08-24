package ssh

import (
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"mime/multipart"
	"sync"
	"teamide/pkg/util"
)

type UploadFile struct {
	Dir      string
	Place    string
	WorkId   string
	FullPath string
	File     *multipart.FileHeader
	WorkerId string
}

type SftpClient struct {
	WorkerId    string
	Config      Config
	sshClient   *ssh.Client
	confirmMap  map[string]chan *util.FileConfirmInfo
	newSftpLock sync.Mutex
}

func (this_ *SftpClient) Start() {
}

func (this_ *SftpClient) AddUpload(uploadFile *UploadFile) {
	if uploadFile == nil {
		return
	}
	go func() {
		_, _ = this_.Work(&SFTPRequest{
			Work:     "upload",
			WorkId:   uploadFile.WorkId,
			Dir:      uploadFile.Dir,
			Place:    uploadFile.Place,
			File:     uploadFile.File,
			FullPath: uploadFile.FullPath,
		})
	}()
	return
}
func (this_ *SftpClient) newSftp() (sftpClient *sftp.Client, err error) {
	this_.newSftpLock.Lock()
	defer this_.newSftpLock.Unlock()

	err = this_.initClient()
	if err != nil {
		return
	}

	sftpClient, err = sftp.NewClient(this_.sshClient)
	if err != nil {
		return
	}

	return
}

func (this_ *SftpClient) Close() {
	this_.closeClient()
	return
}

func (this_ *SftpClient) closeClient() {
	if this_.sshClient != nil {
		_ = this_.sshClient.Close()
		this_.sshClient = nil
	}
	return
}
func (this_ *SftpClient) initClient() (err error) {
	if this_.sshClient == nil {
		err = this_.createClient()
	}
	return
}

func (this_ *SftpClient) createClient() (err error) {

	if this_.sshClient, err = NewClient(this_.Config); err != nil {
		util.Logger.Error("createClient error", zap.Error(err))
		return
	}
	go func() {
		err = this_.sshClient.Wait()
		this_.Close()
	}()
	return
}
