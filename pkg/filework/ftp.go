package filework

import "io"

type FTPService struct {
}

func (this_ *FTPService) Exist(path string) (exist bool, err error) {
	return
}

func (this_ *FTPService) Write(path string, reader io.Reader, onDo func(fileCount *int, fileSize *int64), confirmInfo *ConfirmInfo) (err error) {
	return
}

func (this_ *FTPService) Read(path string, writer io.Writer, onDo func(fileCount *int, fileSize *int64)) (err error) {
	return
}

func (this_ *FTPService) Rename(oldName string, newName string) (err error) {
	return
}

func (this_ *FTPService) Move(fromPath string, fromService Service, toPath string, onDo func(fileCount int, fileSize int64), confirmInfo *ConfirmInfo) (err error) {
	return
}

func (this_ *FTPService) Copy(fromPath string, fromService Service, toPath string, onDo func(fileCount int, fileSize int64), confirmInfo *ConfirmInfo) (err error) {
	return
}

func (this_ *FTPService) Remove(path string, onDo func(fileCount int)) (err error) {
	return
}

func (this_ *FTPService) Count(path string, onDo func(fileCount int)) (fileCount int, err error) {
	return
}

func (this_ *FTPService) CountSize(path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error) {
	return
}

func (this_ *FTPService) Files(path string) (files []*FileInfo, err error) {
	return
}
