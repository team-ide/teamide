package filework

import "io"

type FTPService struct {
}

func (this_ *FTPService) Exist(path string) (exist bool, err error) {
	return
}

func (this_ *FTPService) Create(path string, isDir bool) (file *FileInfo, err error) {
	return
}

func (this_ *FTPService) Write(path string, bytes []byte) (err error) {
	return
}

func (this_ *FTPService) WriteByReader(path string, reader io.Reader, onDo func(readSize int64, writeSize int64)) (file *FileInfo, err error) {
	return
}

func (this_ *FTPService) Read(path string) (bytes []byte, err error) {
	return
}

func (this_ *FTPService) Rename(oldPath string, newPath string) (file *FileInfo, err error) {
	return
}

func (this_ *FTPService) Move(oldPath string, newPath string) (err error) {
	return
}

func (this_ *FTPService) Copy(path string, fromService Service, fromPath string, onDo func(fileCount int, fileSize int64)) (err error) {
	return
}

func (this_ *FTPService) Remove(path string, onDo func(fileCount int, removeCount int)) (err error) {
	return
}

func (this_ *FTPService) Count(path string, onDo func(fileCount int)) (fileCount int, err error) {
	return
}

func (this_ *FTPService) CountSize(path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error) {
	return
}

func (this_ *FTPService) Files(path string) (parentPath string, files []*FileInfo, err error) {
	return
}

func (this_ *FTPService) File(path string) (file *FileInfo, err error) {
	return
}
