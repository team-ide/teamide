package filework

import "io"

type NodeService struct {
}

func (this_ *NodeService) Exist(path string) (exist bool, err error) {
	return
}

func (this_ *NodeService) Create(path string, isDir bool) (file *FileInfo, err error) {
	return
}

func (this_ *NodeService) Write(path string, bytes []byte) (err error) {
	return
}

func (this_ *NodeService) WriteByReader(path string, reader io.Reader, onDo func(readSize int64, writeSize int64)) (file *FileInfo, err error) {
	return
}

func (this_ *NodeService) Read(path string) (bytes []byte, err error) {
	return
}

func (this_ *NodeService) Rename(oldPath string, newPath string) (file *FileInfo, err error) {
	return
}

func (this_ *NodeService) Move(oldPath string, newPath string) (err error) {
	return
}

func (this_ *NodeService) Copy(path string, fromService Service, fromPath string, onDo func(fileCount int, fileSize int64)) (err error) {
	return
}

func (this_ *NodeService) Remove(path string, onDo func(fileCount int, removeCount int)) (err error) {
	return
}

func (this_ *NodeService) Count(path string, onDo func(fileCount int)) (fileCount int, err error) {
	return
}

func (this_ *NodeService) CountSize(path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error) {
	return
}

func (this_ *NodeService) Files(path string) (parentPath string, files []*FileInfo, err error) {

	return
}

func (this_ *NodeService) File(path string) (file *FileInfo, err error) {
	return
}
