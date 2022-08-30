package filework

import "io"

type LocalService struct {
}

func (this_ *LocalService) Exist(path string) (exist bool, err error) {
	return
}

func (this_ *LocalService) Write(path string, reader io.Reader, onDo func(fileCount *int, fileSize *int64), confirmInfo *ConfirmInfo) (err error) {
	return
}

func (this_ *LocalService) Read(path string, writer io.Writer, onDo func(fileCount *int, fileSize *int64)) (err error) {
	return
}

func (this_ *LocalService) Rename(oldName string, newName string) (err error) {
	return
}

func (this_ *LocalService) Move(fromPath string, fromService Service, toPath string, onDo func(fileCount int, fileSize int64), confirmInfo *ConfirmInfo) (err error) {
	return
}

func (this_ *LocalService) Copy(fromPath string, fromService Service, toPath string, onDo func(fileCount int, fileSize int64), confirmInfo *ConfirmInfo) (err error) {
	return
}

func (this_ *LocalService) Remove(path string, onDo func(fileCount int)) (err error) {
	return
}

func (this_ *LocalService) Count(path string, onDo func(fileCount int)) (fileCount int, err error) {
	return
}

func (this_ *LocalService) CountSize(path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error) {
	return
}

func (this_ *LocalService) Files(path string) (files []*FileInfo, err error) {
	return
}
