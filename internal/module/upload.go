package module

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path"
	"teamide/internal/base"
	"teamide/pkg/util"
	"time"
)

type UploadResponse struct {
	Files []*UploadFile `json:"files,omitempty"`
}

type UploadFile struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
	Size int64  `json:"size,omitempty"`
}

func (this_ *Api) apiUpload(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	response := &UploadResponse{}

	place := c.PostForm("place")
	if place == "" {
		err = errors.New("上传位置信息获取失败")
		return
	}
	uploadFile, err := c.FormFile("file")
	if err != nil {
		return
	}
	nowTime := time.Now()
	filePath := fmt.Sprintf("%s/%d/%d/%d/%s", place, nowTime.Year(), nowTime.Month(), nowTime.Day(), util.UUID())

	fileDir := this_.GetFilesFile(filePath)

	exist, err := util.PathExists(fileDir)
	if err != nil {
		return
	}
	if !exist {
		err = os.MkdirAll(fileDir, 0777)
	}

	openFile, err := uploadFile.Open()
	if err != nil {
		return
	}
	defer openFile.Close()
	fileName := uploadFile.Filename
	//获取文件名称带后缀
	fileNameWithSuffix := path.Base(uploadFile.Filename)
	//获取文件的后缀(文件类型)
	fileType := path.Ext(fileNameWithSuffix)
	filePath += "/" + fileName

	fileSavePath := this_.GetFilesFile(filePath)

	createFile, err := os.Create(fileSavePath)
	if err != nil {
		return
	}
	defer func() { _ = createFile.Close() }()

	_, err = io.Copy(createFile, openFile)
	if err != nil {
		return
	}

	response.Files = append(response.Files, &UploadFile{
		Name: fileName,
		Type: fileType,
		Path: filePath,
	})
	res = response
	return
}
