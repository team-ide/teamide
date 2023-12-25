package module_node

import (
	"errors"
	"io"
	"teamide/pkg/filework"
	"teamide/pkg/node"
)

func NewFileService(nodeId string, nodeService *NodeService) *fileService {
	return &fileService{
		nodeId:      nodeId,
		nodeService: nodeService,
	}
}

type fileService struct {
	nodeId      string
	nodeLine    []string
	nodeService *NodeService
}

func (this_ *fileService) getServer() (server *node.Server, err error) {
	if this_.nodeService.GetContext() == nil {
		err = errors.New("node上下文未初始化")
		return
	}
	server = this_.nodeService.GetContext().GetServer()

	nodeLine := this_.nodeService.GetContext().GetNodeLineTo(this_.nodeId)

	if len(nodeLine) == 0 {
		err = errors.New("无法连接到节点[" + this_.nodeId + "]")
		return
	}
	this_.nodeLine = nodeLine

	return
}

func (this_ *fileService) Exist(path string) (exist bool, err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	exist, err = server.FileWorkExist(this_.nodeLine, path)
	return
}

func (this_ *fileService) ExistAndMd5(path string) (exist bool, md5 string, err error) {
	exist, err = this_.Exist(path)
	if err != nil {
		return
	}
	if !exist {
		return
	}
	return
}

func (this_ *fileService) Create(path string, isDir bool) (err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	err = server.FileWorkCreate(this_.nodeLine, path, isDir)
	return
}

func (this_ *fileService) Write(path string, reader io.Reader, onDo func(readSize int64, writeSize int64), callStop *bool) (err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	err = server.FileWorkWrite(this_.nodeLine, path, reader, onDo, callStop)
	return
}

func (this_ *fileService) Read(path string, writer io.Writer, onDo func(readSize int64, writeSize int64), callStop *bool) (err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	err = server.FileWorkRead(this_.nodeLine, path, writer, onDo, callStop)
	return
}

func (this_ *fileService) Rename(oldPath string, newPath string) (err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	err = server.FileWorkRename(this_.nodeLine, oldPath, newPath)
	return
}

func (this_ *fileService) Move(oldPath string, newPath string) (err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	err = server.FileWorkMove(this_.nodeLine, oldPath, newPath)
	return
}

func (this_ *fileService) Remove(path string, onDo func(fileCount int, removeCount int)) (err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	err = server.FileWorkRemove(this_.nodeLine, path, onDo)
	return
}

func (this_ *fileService) Count(path string, onDo func(fileCount int)) (fileCount int, err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	fileCount, err = server.FileWorkCount(this_.nodeLine, path, onDo)
	return
}

func (this_ *fileService) CountSize(path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	fileCount, fileSize, err = server.FileWorkCountSize(this_.nodeLine, path, onDo)
	return
}

func (this_ *fileService) Files(dir string) (parentPath string, files []*filework.FileInfo, err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	parentPath, files, err = server.FileWorkFiles(this_.nodeLine, dir)
	return
}

func (this_ *fileService) File(path string) (file *filework.FileInfo, err error) {
	var server *node.Server
	server, err = this_.getServer()
	if err != nil {
		return
	}

	file, err = server.FileWorkFile(this_.nodeLine, path)
	return
}

func (this_ *fileService) OpenReader(path string) (reader io.ReadCloser, err error) {
	err = errors.New("节点暂不支持该功能")
	return
}
func (this_ *fileService) OpenWriter(path string) (writer io.WriteCloser, err error) {
	err = errors.New("节点暂不支持该功能")
	return
}
