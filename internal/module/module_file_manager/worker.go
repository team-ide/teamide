package module_file_manager

import (
	"errors"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/filework"
	"teamide/pkg/ssh"
)

var (
	toolboxService *module_toolbox.ToolboxService
	nodeContext    *module_node.NodeContext
)

func GetService(fileWorkerKey string, place string, placeId string) (service filework.Service, err error) {
	switch place {
	case "local":
		service = &filework.LocalService{}
	case "ssh":
		if placeId == "" {
			err = errors.New("SSH配置不能为空")
			return
		}
		var id int64
		id, err = strconv.ParseInt(placeId, 10, 64)
		if err != nil {
			return
		}
		var tD *module_toolbox.ToolboxModel
		tD, err = toolboxService.Get(id)
		if err != nil {
			return
		}
		if tD == nil || tD.Option == "" {
			err = errors.New("SSH[" + placeId + "]配置不存在")
			return
		}

		var config *ssh.Config
		config, err = toolboxService.GetSSHConfig(tD.Option)

		service = ssh.CreateOrGetClient(fileWorkerKey, config)
	case "node":
		if placeId == "" {
			err = errors.New("node配置不能为空")
			return
		}
		service = module_node.NewFileService(placeId, nodeContext)
	}
	if service == nil {
		err = errors.New("[" + place + "]文件服务不存在")
		return
	}
	return
}

func Create(workerId string, fileWorkerKey string, place string, placeId string, path string, isDir bool) (file *filework.FileInfo, err error) {
	progress := newProgress(workerId, place, placeId, "create")
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["path"] = path
	progress.Data["isDir"] = isDir
	defer func() { progress.end(err) }()

	service, err := GetService(fileWorkerKey, place, placeId)
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

func CallAction(progressId string, action string) (err error) {
	progress := getProgress(progressId)
	if progress != nil {
		progress.callAction(action)
	}
	return
}

func File(_ string, fileWorkerKey string, place string, placeId string, path string) (file *filework.FileInfo, err error) {

	service, err := GetService(fileWorkerKey, place, placeId)
	if err != nil {
		return
	}
	file, err = service.File(path)
	return
}

func Files(_ string, fileWorkerKey string, place string, placeId string, dir string) (parentPath string, files []*filework.FileInfo, err error) {

	service, err := GetService(fileWorkerKey, place, placeId)
	if err != nil {
		return
	}
	parentPath, files, err = service.Files(dir)
	return
}

func Read(workerId string, fileWorkerKey string, place string, placeId string, path string, writer io.Writer) (file *filework.FileInfo, err error) {
	progress := newProgress(workerId, place, placeId, "read")
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["path"] = path
	progress.Data["readSize"] = 0
	progress.Data["writeSize"] = 0
	progress.Data["successSize"] = 0
	defer func() { progress.end(err) }()

	service, err := GetService(fileWorkerKey, place, placeId)
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
	})
	return
}

func Write(workerId string, fileWorkerKey string, place string, placeId string, path string, reader io.Reader, size int) (file *filework.FileInfo, err error) {
	progress := newProgress(workerId, place, placeId, "write")
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["path"] = path
	progress.Data["readSize"] = 0
	progress.Data["writeSize"] = 0
	progress.Data["size"] = size
	defer func() { progress.end(err) }()

	service, err := GetService(fileWorkerKey, place, placeId)
	if err != nil {
		return
	}

	err = service.Write(path, reader, func(readSize int64, writeSize int64) {
		progress.Data["readSize"] = readSize
		progress.Data["writeSize"] = writeSize
		progress.Data["successSize"] = writeSize
	})
	if err != nil {
		return
	}
	file, err = service.File(path)

	return
}

func Rename(workerId string, fileWorkerKey string, place string, placeId string, oldPath string, newPath string) (file *filework.FileInfo, err error) {
	progress := newProgress(workerId, place, placeId, "rename")
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["oldPath"] = oldPath
	progress.Data["newPath"] = newPath

	defer func() { progress.end(err) }()

	service, err := GetService(fileWorkerKey, place, placeId)
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

func Remove(workerId string, fileWorkerKey string, place string, placeId string, path string) (err error) {
	progress := newProgress(workerId, place, placeId, "remove")
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["path"] = path
	progress.Data["fileCount"] = 0
	progress.Data["removeCount"] = 0
	defer func() { progress.end(err) }()

	service, err := GetService(fileWorkerKey, place, placeId)
	if err != nil {
		return
	}
	err = service.Remove(path, func(fileCount int, removeCount int) {
		progress.Data["fileCount"] = fileCount
		progress.Data["removeCount"] = removeCount
	})
	return
}

func Move(workerId string, fileWorkerKey string, place string, placeId string, oldPath string, newPath string) (err error) {
	progress := newProgress(workerId, place, placeId, "move")
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["oldPath"] = oldPath
	progress.Data["newPath"] = newPath
	defer func() { progress.end(err) }()

	toService, err := GetService(fileWorkerKey, place, placeId)
	if err != nil {
		return
	}

	err = toService.Move(oldPath, newPath)
	return
}

func Copy(workerId string, fileWorkerKey string, place string, placeId string, path string, fromPlace string, fromPlaceId string, fromPath string) {
	var err error
	progress := newProgress(workerId, place, placeId, "copy")
	progress.Data["fileWorkerKey"] = fileWorkerKey
	progress.Data["path"] = path
	progress.Data["fromPlace"] = fromPlace
	progress.Data["fromPlaceId"] = fromPlaceId
	progress.Data["fromPath"] = fromPath

	defer func() { progress.end(err) }()

	//toService, err := GetService(place, placeId)
	//if err != nil {
	//	return
	//}
	//
	//fromService, err := GetService(fromPlace, fromPlaceId)
	//if err != nil {
	//	return
	//}

	//err = toService.Copy(path, fromService, fromPath, func(fileCount int, fileSize int64) {
	//
	//})
	return
}

func Upload(workerId string, fileWorkerKey string, place string, placeId string, dir string, fullPath string, fileList []*multipart.FileHeader) (fileInfoList []*filework.FileInfo, err error) {

	service, err := GetService(fileWorkerKey, place, placeId)
	if err != nil {
		return
	}

	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	if strings.HasPrefix(fullPath, "/") {
		fullPath = fullPath[1:]
	}

	upload := func(one *multipart.FileHeader) (fileInfo *filework.FileInfo, err error) {

		path := dir + one.Filename
		if len(fullPath) > 0 {
			path = dir + fullPath
		}

		progress := newProgress(workerId, place, placeId, "upload")
		progress.Data["fileWorkerKey"] = fileWorkerKey
		progress.Data["dir"] = dir
		progress.Data["fullPath"] = fullPath
		progress.Data["filename"] = one.Filename
		progress.Data["path"] = path
		progress.Data["size"] = one.Size
		progress.Data["successSize"] = 0

		defer func() { progress.end(err) }()

		var exist bool
		exist, err = service.Exist(path)

		if exist {
			action := progress.waitAction("文件["+path+"]已存在，是否覆盖？",
				[]*Action{
					newAction("是", "yes", "color-green"),
					newAction("否", "no", "color-orange"),
				})
			if action == "no" {
				return
			}
		}
		var openF multipart.File
		openF, err = one.Open()
		if err != nil {
			return
		}

		err = service.Write(path, openF, func(readSize int64, writeSize int64) {
			progress.Data["successSize"] = writeSize
		})
		if err != nil {
			return
		}

		fileInfo, err = service.File(path)
		return
	}

	for _, one := range fileList {
		var fileInfo *filework.FileInfo
		fileInfo, err = upload(one)
		if err != nil {
			return
		}
		if fileInfo != nil {
			fileInfoList = append(fileInfoList, fileInfo)
		}
	}

	return
}
