package filework

import (
	"mime/multipart"
	"strings"
)

func GetService(place string, placeId string) (service Service, err error) {
	switch place {
	case "local":
		service = &LocalService{}
	}
	return
}

func Create(workerId string, place string, placeId string, path string, isDir bool) (file *FileInfo, err error) {
	progress := newProgress(workerId, place, placeId, "create")
	progress.Data["path"] = path
	progress.Data["isDir"] = isDir
	defer func() { progress.end(err) }()

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}

	file, err = service.Create(path, isDir)
	return
}

func CallAction(progressId string, action string) (err error) {
	progress := getProgress(progressId)
	if progress != nil {
		progress.callAction(action)
	}
	return
}

func File(workerId string, place string, placeId string, path string) (file *FileInfo, err error) {

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	file, err = service.File(path)
	return
}

func Files(workerId string, place string, placeId string, dir string) (parentPath string, files []*FileInfo, err error) {

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	parentPath, files, err = service.Files(dir)
	return
}

func Read(workerId string, place string, placeId string, path string) (bytes []byte, err error) {

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	bytes, err = service.Read(path)
	return
}

func Write(workerId string, place string, placeId string, path string, bytes []byte) (err error) {
	progress := newProgress(workerId, place, placeId, "write")
	progress.Data["path"] = path
	defer func() { progress.end(err) }()

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	err = service.Write(path, bytes)
	return
}

func Rename(workerId string, place string, placeId string, oldPath string, newPath string) (file *FileInfo, err error) {
	progress := newProgress(workerId, place, placeId, "rename")
	progress.Data["oldPath"] = oldPath
	progress.Data["newPath"] = newPath

	defer func() { progress.end(err) }()

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	file, err = service.Rename(oldPath, newPath)
	return
}

func Remove(workerId string, place string, placeId string, path string) (err error) {
	progress := newProgress(workerId, place, placeId, "remove")
	progress.Data["path"] = path
	progress.Data["fileCount"] = 0
	progress.Data["removeCount"] = 0
	defer func() { progress.end(err) }()

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	err = service.Remove(path, func(fileCount int, removeCount int) {
		progress.Data["fileCount"] = fileCount
		progress.Data["removeCount"] = removeCount
	})
	return
}

func Move(workerId string, place string, placeId string, oldPath string, newPath string) (err error) {
	progress := newProgress(workerId, place, placeId, "move")
	progress.Data["oldPath"] = oldPath
	progress.Data["newPath"] = newPath
	defer func() { progress.end(err) }()

	toService, err := GetService(place, placeId)
	if err != nil {
		return
	}

	err = toService.Move(oldPath, newPath)
	return
}

func Copy(workerId string, place string, placeId string, path string, fromPlace string, fromPlaceId string, fromPath string) {
	var err error
	progress := newProgress(workerId, place, placeId, "copy")
	progress.Data["path"] = path
	progress.Data["fromPlace"] = fromPlace
	progress.Data["fromPlaceId"] = fromPlaceId
	progress.Data["fromPath"] = fromPath

	defer func() { progress.end(err) }()

	toService, err := GetService(place, placeId)
	if err != nil {
		return
	}

	fromService, err := GetService(fromPlace, fromPlaceId)
	if err != nil {
		return
	}

	err = toService.Copy(path, fromService, fromPath, func(fileCount int, fileSize int64) {

	})
	return
}

func Upload(workerId string, place string, placeId string, dir string, fullPath string, fileList []*multipart.FileHeader) (fileInfoList []*FileInfo, err error) {

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}

	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	if strings.HasPrefix(fullPath, "/") {
		fullPath = fullPath[1:]
	}

	upload := func(one *multipart.FileHeader) (fileInfo *FileInfo, err error) {

		path := dir + one.Filename
		if len(fullPath) > 0 {
			path = dir + fullPath
		}

		progress := newProgress(workerId, place, placeId, "upload")
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

		fileInfo, err = service.WriteByReader(path, openF, func(readSize int64, writeSize int64) {
			progress.Data["successSize"] = writeSize
		})

		return
	}

	for _, one := range fileList {
		var fileInfo *FileInfo
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
