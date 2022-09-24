package filework

import "io"

type FileInfo struct {
	Name     string `json:"name,omitempty"`
	IsDir    bool   `json:"isDir,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Path     string `json:"path,omitempty"`
	ModTime  int64  `json:"modTime,omitempty"`
	FileMode string `json:"fileMode,omitempty"`
}

type Service interface {
	Exist(path string) (exist bool, err error)
	Create(path string, isDir bool) (err error)
	Write(path string, reader io.Reader, onDo func(readSize int64, writeSize int64), callStop *bool) (err error)
	Read(path string, writer io.Writer, onDo func(readSize int64, writeSize int64), callStop *bool) (err error)
	Rename(oldPath string, newPath string) (err error)
	Move(oldPath string, newPath string) (err error)
	Remove(path string, onDo func(fileCount int, removeCount int)) (err error)
	Count(path string, onDo func(fileCount int)) (fileCount int, err error)
	CountSize(path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error)
	Files(dir string) (parentPath string, files []*FileInfo, err error)
	File(path string) (file *FileInfo, err error)
}
