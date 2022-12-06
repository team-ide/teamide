package module_toolbox

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-dialect/worker"
	"io"
	"net/http"
	"net/url"
	"os"
	"teamide/internal/base"
	"teamide/pkg/util"
)

type WorkRequest struct {
	ToolboxId   int64                  `json:"toolboxId,omitempty"`
	ToolboxType string                 `json:"toolboxType,omitempty"`
	Work        string                 `json:"work,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

func (this_ *ToolboxApi) work(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &WorkRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = this_.ToolboxService.Work(requestBean, request.ToolboxId, request.ToolboxType, request.Work, request.Data)
	if err != nil {
		return
	}

	return
}

func (this_ *ToolboxApi) databaseDownload(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	data := map[string]string{}
	err = c.Bind(&data)
	if err != nil {
		return
	}

	taskId := data["taskId"]
	if taskId == "" {
		err = errors.New("taskId获取失败")
		return
	}

	task := worker.GetTask(taskId)
	if task == nil {
		err = errors.New("任务不存在")
		return
	}
	if task.Extend == nil || task.Extend["downloadPath"] == "" {
		err = errors.New("任务导出文件丢失")
		return
	}
	tempDir, err := util.GetTempDir()
	if err != nil {
		return
	}

	path := tempDir + task.Extend["downloadPath"].(string)
	exists, err := util.PathExists(path)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New("文件不存在")
		return
	}
	var fileName string
	var fileSize int64
	ff, err := os.Lstat(path)
	if err != nil {
		return
	}
	var fileInfo *os.File
	if ff.IsDir() {
		exists, err = util.PathExists(path + ".zip")
		if err != nil {
			return
		}
		if !exists {
			err = util.Zip(path, path+".zip")
			if err != nil {
				return
			}
		}
		ff, err = os.Lstat(path + ".zip")
		if err != nil {
			return
		}
		fileInfo, err = os.Open(path + ".zip")
		if err != nil {
			return
		}
	} else {
		fileInfo, err = os.Open(path)
		if err != nil {
			return
		}
	}
	fileName = ff.Name()
	fileSize = ff.Size()

	defer func() {
		_ = fileInfo.Close()
	}()

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", fmt.Sprint(fileSize))
	c.Header("download-file-name", fileName)

	_, err = io.Copy(c.Writer, fileInfo)
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	res = base.HttpNotResponse
	return
}
