package node

import (
	"io"
	"teamide/pkg/filework"
)

func (this_ *Server) FileWorkExist(lineNodeIdList []string, path string) (exist bool, err error) {
	return
}

func (this_ *Server) FileWorkCreate(lineNodeIdList []string, path string, isDir bool) (err error) {
	return
}

func (this_ *Server) FileWorkWrite(lineNodeIdList []string, path string, reader io.Reader, onDo func(readSize int64, writeSize int64)) (err error) {
	return
}

func (this_ *Server) FileWorkRead(lineNodeIdList []string, path string, writer io.Writer, onDo func(readSize int64, writeSize int64)) (err error) {
	return
}

func (this_ *Server) FileWorkRename(lineNodeIdList []string, oldPath string, newPath string) (err error) {
	return
}

func (this_ *Server) FileWorkMove(lineNodeIdList []string, oldPath string, newPath string) (err error) {
	return
}

func (this_ *Server) FileWorkRemove(lineNodeIdList []string, path string, onDo func(fileCount int, removeCount int)) (err error) {
	return
}

func (this_ *Server) FileWorkCount(lineNodeIdList []string, path string, onDo func(fileCount int)) (fileCount int, err error) {
	return
}

func (this_ *Server) FileWorkCountSize(lineNodeIdList []string, path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error) {
	return
}

func (this_ *Server) FileWorkFiles(lineNodeIdList []string, dir string) (parentPath string, files []*filework.FileInfo, err error) {
	return
}

func (this_ *Server) FileWorkFile(lineNodeIdList []string, path string) (file *filework.FileInfo, err error) {
	return
}
