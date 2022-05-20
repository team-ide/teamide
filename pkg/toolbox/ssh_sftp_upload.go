package toolbox

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func SFTPUpload(c *gin.Context) (res interface{}, err error) {
	token := c.PostForm("token")
	if token == "" {
		err = errors.New("token获取失败")
		return
	}
	dir := c.PostForm("dir")
	if dir == "" {
		err = errors.New("dir获取失败")
		return
	}
	place := c.PostForm("place")
	if place == "" {
		err = errors.New("place获取失败")
		return
	}
	workId := c.PostForm("workId")
	if workId == "" {
		err = errors.New("workId获取失败")
		return
	}
	client := SSHSftpCache[token]
	if client == nil {
		err = errors.New("FTP会话丢失")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		return
	}

	uploadFile := &UploadFile{
		Dir:      dir,
		Place:    place,
		WorkId:   workId,
		File:     file,
		FullPath: c.PostForm("fullPath"),
	}
	client.UploadFile <- uploadFile

	return
}

func SFTPDownload(data map[string]string, c *gin.Context) (err error) {

	token := data["token"]
	if token == "" {
		err = errors.New("token获取失败")
		return
	}
	place := data["place"]
	if place == "" {
		err = errors.New("place获取失败")
		return
	}
	path := data["path"]
	if path == "" {
		err = errors.New("path获取失败")
		return
	}
	client := SSHSftpCache[token]
	if client == nil {
		err = errors.New("SSH会话丢失")
		return
	}
	if place == "local" {
		err = client.localDownload(c, path)
	} else if place == "remote" {
		err = client.remoteDownload(c, path)
	}

	return
}
