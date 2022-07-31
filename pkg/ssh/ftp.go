package ssh

import (
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"mime/multipart"
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
	WorkerId   string
	Config     Config
	sshClient  *ssh.Client
	UploadFile chan *UploadFile
	confirmMap map[string]chan *util.FileConfirmInfo
}

func (this_ *SftpClient) Start() {
	this_.listenUpload()
}
func (this_ *SftpClient) listenUpload() {
	if this_.UploadFile == nil {
		this_.UploadFile = make(chan *UploadFile, 10)

		go func() {
			for {
				select {
				case uploadFile := <-this_.UploadFile:
					_, _ = this_.Work(&SFTPRequest{
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
	}
	return
}
func (this_ *SftpClient) newSftp() (sftpClient *sftp.Client, err error) {
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
