package module_toolbox

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"os"
	"path/filepath"
	"teamide/pkg/base"
)

func (this_ *ToolboxApi) extendGet(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ToolboxExtendModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = this_.ToolboxService.GetExtend(request.ExtendId)
	if err != nil {
		return
	}

	return
}

func (this_ *ToolboxApi) extendQuery(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ToolboxExtendModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.UserId == -1 {
		request.UserId = 0
	} else {
		request.UserId = requestBean.JWT.UserId
	}

	res, err = this_.ToolboxService.QueryExtends(request)
	if err != nil {
		return
	}

	return
}

func (this_ *ToolboxApi) extendSave(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ToolboxExtendModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	request.UserId = requestBean.JWT.UserId
	err = this_.ToolboxService.SaveExtend(request)
	if err != nil {
		return
	}
	res = request
	return
}

func (this_ *ToolboxApi) extendDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ToolboxExtendModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = this_.ToolboxService.DeleteExtend(request.ExtendId)
	if err != nil {
		return
	}

	return
}

func (this_ *ToolboxApi) extendLoadFile(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ToolboxExtendModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	find, err := this_.ToolboxService.GetExtend(request.ExtendId)
	if err != nil {
		return
	}
	if find == nil || find.Extend == nil {
		err = errors.New("扩展数据不存在")
		return
	}
	filePath := util.GetStringValue(find.Extend["filePath"])
	if filePath == "" {
		err = errors.New("扩展数据文件地址不存在")
		return
	}
	filePath = this_.GetFilesFile(filePath)
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return
	}
	res = string(bs)

	return
}

type ExtendSaveFile struct {
	ExtendId int64  `json:"extendId,omitempty"`
	Text     string `json:"text,omitempty"`
}

func (this_ *ToolboxApi) extendSaveFile(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &ExtendSaveFile{}
	if !base.RequestJSON(request, c) {
		return
	}

	find, err := this_.ToolboxService.GetExtend(request.ExtendId)
	if err != nil {
		return
	}
	if find == nil || find.Extend == nil {
		err = errors.New("扩展数据不存在")
		return
	}
	filePath := util.GetStringValue(find.Extend["filePath"])
	if filePath == "" {
		err = errors.New("扩展数据文件地址不存在")
		return
	}
	filePath = this_.GetFilesFile(filePath)

	dir := filepath.Dir(filePath)
	if e, _ := util.PathExists(dir); !e {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.WriteString(request.Text)

	fs, err := os.Stat(filePath)
	if err != nil {
		return
	}
	find.Extend["fileSize"] = fs.Size()
	err = this_.ToolboxService.SaveExtend(find)
	if err != nil {
		return
	}
	return
}
