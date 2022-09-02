package filework

func GetService(place string, placeId string) (service Service, err error) {
	switch place {
	case "local":
		service = &LocalService{}
	}
	return
}

func File(place string, placeId string, path string) (file *FileInfo, err error) {

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	file, err = service.File(path)
	return
}

func Files(place string, placeId string, dir string) (files []*FileInfo, err error) {

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	files, err = service.Files(dir)
	return
}

func Read(place string, placeId string, path string) (bytes []byte, err error) {

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	bytes, err = service.Read(path)
	return
}

func Rename(place string, placeId string, oldName string, newName string) (err error) {

	service, err := GetService(place, placeId)
	if err != nil {
		return
	}
	err = service.Rename(oldName, newName)
	return
}
