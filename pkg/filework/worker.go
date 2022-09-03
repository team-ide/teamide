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
	progress := newProgress(workerId, "create")
	progress.Data["place"] = place
	progress.Data["placeId"] = placeId
	progress.Data["path"] = path
	progress.Data["isDir"] = isDir
	defer progress.end(err)

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
	progress := newProgress(workerId, "write")
	progress.Data["place"] = place
	progress.Data["placeId"] = placeId
	progress.Data["path"] = path
	defer progress.end(err)

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	err = service.Write(path, bytes)
	return
}

func Rename(workerId string, place string, placeId string, oldPath string, newPath string) (file *FileInfo, err error) {
	progress := newProgress(workerId, "rename")
	progress.Data["place"] = place
	progress.Data["placeId"] = placeId
	progress.Data["oldPath"] = oldPath
	progress.Data["newPath"] = newPath
	defer progress.end(err)

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	file, err = service.Rename(oldPath, newPath)
	return
}

func Remove(workerId string, place string, placeId string, path string) (err error) {
	progress := newProgress(workerId, "remove")
	progress.Data["place"] = place
	progress.Data["placeId"] = placeId
	progress.Data["path"] = path
	progress.Data["fileCount"] = 0
	progress.Data["removeCount"] = 0
	defer progress.end(err)

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

func Copy(workerId string, fromPlace string, fromPlaceId string, fromPath string, toPlace string, toPlaceId string, toPath string) {
	var err error
	progress := newProgress(workerId, "copy")
	progress.Data["fromPlace"] = fromPlace
	progress.Data["fromPlaceId"] = fromPlaceId
	progress.Data["fromPath"] = fromPath
	progress.Data["toPlace"] = toPlace
	progress.Data["toPlaceId"] = toPlaceId
	progress.Data["toPath"] = toPath
	defer progress.end(err)

	fromService, err := GetService(fromPlace, fromPlaceId)
	if err != nil {
		return
	}

	toService, err := GetService(toPlace, toPlaceId)
	if err != nil {
		return
	}
	err = toService.Copy(fromPath, fromService, toPath, func(fileCount int, fileSize int64) {

	})
	return
}

func Upload(workerId string, place string, placeId string, dir string, fullPath string, fileList []*multipart.FileHeader) {
	var err error
	progress := newProgress(workerId, "copy")
	progress.Data["place"] = place
	progress.Data["placeId"] = placeId
	progress.Data["dir"] = dir
	progress.Data["fullPath"] = fullPath
	defer progress.end(err)

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}

	if len(fullPath) > 0 {
		dir += fullPath
	}
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	for _, one := range fileList {
		path := dir + "/" + one.Filename

		var exist bool
		exist, err = service.Exist(path)

		if exist {

		}
	}

	return
}
