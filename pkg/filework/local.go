package filework

import (
	"os"
	"path/filepath"
	"teamide/pkg/util"
)

type LocalService struct {
}

func (this_ *LocalService) Exist(path string) (exist bool, err error) {
	return
}

func (this_ *LocalService) Write(path string, onDo func(fileCount *int, fileSize *int64), confirmInfo *ConfirmInfo) (err error) {
	return
}

func (this_ *LocalService) Read(path string) (bytes []byte, err error) {
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

func (this_ *LocalService) File(path string) (file *FileInfo, err error) {
	return
}

func getFileInfo(path string) (fileInfo *FileInfo, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}
		return
	}

	fileInfo = &FileInfo{
		Name:     filepath.VolumeName(path),
		Path:     path,
		IsDir:    true,
		ModTime:  util.GetTimeTime(stat.ModTime()),
		FileMode: stat.Mode().String(),
		Size:     stat.Size(),
	}
	return
}
