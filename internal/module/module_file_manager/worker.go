package module_file_manager

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"teamide/internal/context"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/filework"
	"teamide/pkg/ssh"
)

func NewWorker(toolboxService_ *module_toolbox.ToolboxService, nodeService_ *module_node.NodeService) *worker {
	return &worker{
		ServerContext:  toolboxService_.ServerContext,
		toolboxService: toolboxService_,
		nodeService:    nodeService_,
	}
}

type worker struct {
	*context.ServerContext
	toolboxService *module_toolbox.ToolboxService
	nodeService    *module_node.NodeService
}

func (this_ *worker) GetService(fileWorkerKey string, param *BaseParam) (service filework.Service, err error) {
	switch param.Place {
	case "local":
		service = filework.NewLocalService()
	case "ssh":
		if param.PlaceId == "" {
			err = errors.New("SSH配置不能为空")
			return
		}
		var id int64
		id, err = strconv.ParseInt(param.PlaceId, 10, 64)
		if err != nil {
			return
		}
		var tD *module_toolbox.ToolboxModel
		tD, err = this_.toolboxService.Get(id)
		if err != nil {
			return
		}
		if tD == nil || tD.Option == "" {
			err = errors.New("SSH[" + param.PlaceId + "]配置不存在")
			return
		}

		var config *ssh.Config
		config, err = this_.toolboxService.GetSSHConfig(tD.Option)

		service = ssh.CreateOrGetClient(fileWorkerKey, config)
	case "node":
		if param.PlaceId == "" {
			err = errors.New("node配置不能为空")
			return
		}
		service = module_node.NewFileService(param.PlaceId, this_.nodeService)
	}
	if service == nil {
		err = errors.New("[" + param.Place + "]文件服务不存在")
		return
	}
	return
}

func (this_ *worker) Close(workerId string) {
	progressList := getProgressList(workerId)
	for _, one := range progressList {
		one.closeCallAction()
	}

	return
}

func (this_ *worker) Create(param *BaseParam, fileWorkerKey string, path string, isDir bool) (file *filework.FileInfo, err error) {
	progress := newProgress(param, "create", func() {

	})
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["path"] = path
	progress.Data["isDir"] = isDir

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		progress.end(err)
	}()

	service, err := this_.GetService(fileWorkerKey, param)
	if err != nil {
		return
	}

	err = service.Create(path, isDir)
	if err != nil {
		return
	}

	file, err = service.File(path)
	return
}

func (this_ *worker) CallAction(progressId string, action string) (err error) {
	progress := getProgress(progressId)
	if progress != nil {
		progress.callAction(action)
	}
	return
}

func (this_ *worker) CallStop(progressId string) (err error) {
	progress := getProgress(progressId)
	if progress != nil {
		progress.callStopped = true
	}
	return
}

func (this_ *worker) File(param *BaseParam, fileWorkerKey string, path string) (file *filework.FileInfo, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	service, err := this_.GetService(fileWorkerKey, param)
	if err != nil {
		return
	}
	file, err = service.File(path)
	return
}

func (this_ *worker) Files(param *BaseParam, fileWorkerKey string, dir string) (parentPath string, files []*filework.FileInfo, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	service, err := this_.GetService(fileWorkerKey, param)
	if err != nil {
		return
	}
	parentPath, files, err = service.Files(dir)
	return
}

func (this_ *worker) Read(param *BaseParam, fileWorkerKey string, path string, writer io.Writer) (file *filework.FileInfo, err error) {
	var false_ = false
	var callStop *bool = &false_
	progress := newProgress(param, "read", func() {
		*callStop = true
	})
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["path"] = path
	progress.Data["readSize"] = 0
	progress.Data["writeSize"] = 0
	progress.Data["successSize"] = 0
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		progress.end(err)
	}()

	service, err := this_.GetService(fileWorkerKey, param)
	if err != nil {
		return
	}

	file, err = service.File(path)
	if err != nil {
		return
	}
	progress.Data["size"] = file.Size

	err = service.Read(path, writer, func(readSize int64, writeSize int64) {
		progress.Data["readSize"] = readSize
		progress.Data["writeSize"] = writeSize
		progress.Data["successSize"] = writeSize
	}, callStop)
	return
}

func (this_ *worker) Write(param *BaseParam, fileWorkerKey string, path string, reader io.Reader, size int) (file *filework.FileInfo, err error) {
	var false_ = false
	var callStop *bool = &false_
	progress := newProgress(param, "write", func() {
		*callStop = true
	})
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["path"] = path
	progress.Data["readSize"] = 0
	progress.Data["writeSize"] = 0
	progress.Data["size"] = size

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		progress.end(err)
	}()

	service, err := this_.GetService(fileWorkerKey, param)
	if err != nil {
		return
	}

	err = service.Write(path, reader, func(readSize int64, writeSize int64) {
		progress.Data["readSize"] = readSize
		progress.Data["writeSize"] = writeSize
		progress.Data["successSize"] = writeSize
	}, callStop)
	if err != nil {
		return
	}
	file, err = service.File(path)

	return
}

func (this_ *worker) Rename(param *BaseParam, fileWorkerKey string, oldPath string, newPath string) (file *filework.FileInfo, err error) {
	progress := newProgress(param, "rename", func() {

	})
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["oldPath"] = oldPath
	progress.Data["newPath"] = newPath

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		progress.end(err)
	}()

	service, err := this_.GetService(fileWorkerKey, param)
	if err != nil {
		return
	}
	err = service.Rename(oldPath, newPath)
	if err != nil {
		return
	}

	file, err = service.File(newPath)
	return
}

func (this_ *worker) Remove(param *BaseParam, fileWorkerKey string, path string) (err error) {
	progress := newProgress(param, "remove", func() {

	})
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["path"] = path
	progress.Data["fileCount"] = 0
	progress.Data["removeCount"] = 0

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		progress.end(err)
	}()

	service, err := this_.GetService(fileWorkerKey, param)
	if err != nil {
		return
	}
	err = service.Remove(path, func(fileCount int, removeCount int) {
		progress.Data["fileCount"] = fileCount
		progress.Data["removeCount"] = removeCount
	})
	return
}

func (this_ *worker) Move(param *BaseParam, fileWorkerKey string, oldPath string, newPath string) (err error) {
	progress := newProgress(param, "move", func() {

	})
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["oldPath"] = oldPath
	progress.Data["newPath"] = newPath

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		progress.end(err)
	}()

	toService, err := this_.GetService(fileWorkerKey, param)
	if err != nil {
		return
	}

	err = toService.Move(oldPath, newPath)
	return
}

func (this_ *worker) Copy(param *BaseParam, fileWorkerKey string, path string, fromFileWorkerKey string, fromPlace string, fromPlaceId string, fromPath string) {
	var err error
	callStop := new(bool)
	progress := newProgress(param, "copy", func() {
		*callStop = true
	})
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["path"] = path
	progress.Data["fromFileWorkerKey"] = fromFileWorkerKey
	progress.Data["fromPlace"] = fromPlace
	progress.Data["fromPlaceId"] = fromPlaceId
	progress.Data["fromPath"] = fromPath
	progress.Data["sameFile"] = false

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		progress.end(err)
	}()

	toService, err := this_.GetService(fileWorkerKey, param)
	if err != nil {
		return
	}

	fromService, err := this_.GetService(fromFileWorkerKey, &BaseParam{
		Place:   fromPlace,
		PlaceId: fromPlaceId,
	})
	if err != nil {
		return
	}

	fromFile, err := fromService.File(fromPath)
	if err != nil {
		return
	}
	//var writer io.WriteCloser
	var reader io.ReadCloser
	if !fromFile.IsDir {

		var exist bool
		var toMd5 string
		var fromMd5 string
		exist, toMd5, err = toService.ExistAndMd5(path)

		if exist {
			if toMd5 != "" {
				exist, fromMd5, err = fromService.ExistAndMd5(fromPath)
				if err != nil {
					return
				}
				if fromMd5 == toMd5 {
					progress.Data["sameFile"] = true
					return
				}
			}

			var action string
			action, err = progress.waitAction("文件["+path+"]已存在，是否覆盖？",
				[]*Action{
					newAction("是", "yes", "color-green"),
					newAction("否", "no", "color-orange"),
				})
			if err != nil {
				return
			}
			if action != "yes" {
				return
			}
		}

		reader, err = fromService.OpenReader(fromPath)
		if err != nil {
			err = errors.New("get reader error:" + err.Error())
			return
		}
		defer func() { _ = reader.Close() }()

		err = toService.Write(path, reader, func(readSize int64, writeSize int64) {
			progress.Data["readSize"] = readSize
			progress.Data["writeSize"] = writeSize
			progress.Data["successSize"] = writeSize
		}, callStop)
		if err != nil {
			return
		}
		progress.Data["fileInfo"], _ = this_.File(param, fileWorkerKey, path)

	} else {
		err = errors.New("暂不支持移动文件夹")
	}

	return
}
