package node

import (
	"io"
	"sync"
	"teamide/pkg/filework"
	"teamide/pkg/util"
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

	sendKey, err := this_.workFileWrite(lineNodeIdList, path)
	if err != nil {
		return
	}

	err = this_.workSendBytesStart(lineNodeIdList, sendKey)

	if err != nil {
		return
	}

	var readSize int64
	var writeSize int64
	var buf = make([]byte, 1024*32)
	err = util.Read(reader, buf, func(n int) (e error) {
		readSize += int64(n)
		onDo(readSize, writeSize)
		e = this_.workSendBytes(lineNodeIdList, sendKey, buf[:n])
		writeSize += int64(n)
		onDo(readSize, writeSize)
		return
	})
	err = this_.workSendBytesEnd(lineNodeIdList, sendKey)

	if err != nil {
		return
	}

	return
}

func (this_ *Server) FileWorkRead(lineNodeIdList []string, path string, writer io.Writer, onDo func(readSize int64, writeSize int64)) (err error) {

	sendKey := util.UUID()

	var waitGroupForStop sync.WaitGroup
	waitGroupForStop.Add(1)
	var readSize int64
	var writeSize int64
	this_.addOnBytesCache(sendKey, &OnBytes{
		start: func() (err error) {
			return
		},
		on: func(buf []byte) (err error) {
			n := len(buf)
			readSize += int64(n)
			onDo(readSize, writeSize)
			err = util.Write(writer, buf, func(n int) (e error) {
				writeSize += int64(n)
				onDo(readSize, writeSize)
				return
			})
			return
		},
		end: func() (err error) {
			waitGroupForStop.Done()
			return
		},
	})
	err = this_.workFileRead(lineNodeIdList, path, sendKey)
	if err != nil {
		this_.removeOnBytesCache(sendKey)
		return
	}

	waitGroupForStop.Wait()
	return
}

func (this_ *Server) FileWorkRemove(lineNodeIdList []string, path string, onDo func(fileCount int, removeCount int)) (err error) {
	fileCount, removeCount, err := this_.workFileRemove(lineNodeIdList, path)
	onDo(fileCount, removeCount)
	if err != nil {
		return
	}
	return
}

func (this_ *Server) FileWorkCount(lineNodeIdList []string, path string, onDo func(fileCount int)) (fileCount int, err error) {
	return
}

func (this_ *Server) FileWorkCountSize(lineNodeIdList []string, path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error) {
	return
}
