package node

import (
	"io"
	"os"
	"teamide/pkg/filework"
)

func (this_ *Server) FileWorkExist(lineNodeIdList []string, path string) (exist bool, err error) {
	exist, err = this_.workExist(lineNodeIdList, path)
	return
}

func (this_ *Server) FileWorkCreate(lineNodeIdList []string, path string, isDir bool) (err error) {
	err = this_.workFileCreate(lineNodeIdList, path, isDir)
	return
}

func (this_ *Server) FileWorkRename(lineNodeIdList []string, oldPath string, newPath string) (err error) {
	err = this_.workFileRename(lineNodeIdList, oldPath, newPath)
	return
}

func (this_ *Server) FileWorkMove(lineNodeIdList []string, oldPath string, newPath string) (err error) {
	err = this_.workFileMove(lineNodeIdList, oldPath, newPath)
	return
}

func (this_ *Server) FileWorkFiles(lineNodeIdList []string, dir string) (parentPath string, files []*filework.FileInfo, err error) {
	parentPath, files, err = this_.workFiles(lineNodeIdList, dir)
	return
}

func (this_ *Server) FileWorkFile(lineNodeIdList []string, path string) (file *filework.FileInfo, err error) {
	file, err = this_.workFile(lineNodeIdList, path)
	return
}

func (this_ *Server) FileWorkWrite(lineNodeIdList []string, path string, reader io.Reader, onDo func(readSize int64, writeSize int64)) (err error) {
	return
}

func (this_ *Server) FileWorkRead(lineNodeIdList []string, path string, writer io.Writer, onDo func(readSize int64, writeSize int64)) (err error) {
	return
}

func (this_ *Server) FileWorkRemove(lineNodeIdList []string, path string, onDo func(fileCount int, removeCount int)) (err error) {
	var fileCount int
	var removeCount int

	err = this_.fileRemove(lineNodeIdList, path, func() {
		fileCount++
		onDo(fileCount, removeCount)
	}, func() {
		removeCount++
		onDo(fileCount, removeCount)
	})

	return
}

func (this_ *Server) fileRemove(lineNodeIdList []string, path string, onLoad func(), onRemove func()) (err error) {
	var isDir bool

	var info *filework.FileInfo
	info, err = this_.FileWorkFile(lineNodeIdList, path)
	if err != nil {
		return
	}
	isDir = info.IsDir

	onLoad()
	if isDir {
		var ds []*filework.FileInfo
		_, ds, err = this_.FileWorkFiles(lineNodeIdList, path)
		if err != nil {
			return
		}

		for _, d := range ds {
			if d.Name == ".." {
				continue
			}
			err = this_.fileRemove(lineNodeIdList, path+"/"+d.Name, onLoad, onRemove)
			if err != nil {
				return
			}
		}
	}
	err = os.Remove(path)
	if err != nil {
		return
	}
	onRemove()

	return
}

func (this_ *Server) FileWorkCount(lineNodeIdList []string, path string, onDo func(fileCount int)) (fileCount int, err error) {
	return
}

func (this_ *Server) FileWorkCountSize(lineNodeIdList []string, path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error) {
	return
}
