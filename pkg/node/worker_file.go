package node

import (
	"teamide/pkg/filework"
)

type FileInfo struct {
	Name     string `json:"name,omitempty"`
	IsDir    bool   `json:"isDir,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Place    string `json:"place,omitempty"`
	Path     string `json:"path,omitempty"`
	ModTime  int64  `json:"modTime,omitempty"`
	FileMode string `json:"fileMode,omitempty"`
}

func (this_ *Worker) workExist(lineNodeIdList []string, path string) (exist bool, err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodFileExist, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				Path: path,
			},
		})
		if e != nil {
			return
		}

		if res != nil && res.FileWorkData != nil {
			exist = res.FileWorkData.Exist
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	exist, err = filework.NewLocalService().Exist(path)

	return
}

func (this_ *Worker) workFile(lineNodeIdList []string, path string) (file *filework.FileInfo, err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodFileFile, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				Path: path,
			},
		})
		if e != nil {
			return
		}

		if res != nil && res.FileWorkData != nil {
			file = res.FileWorkData.File
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	file, err = filework.NewLocalService().File(path)

	return
}

func (this_ *Worker) workFiles(lineNodeIdList []string, dir string) (path string, fileList []*filework.FileInfo, err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodFileFiles, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				Dir: dir,
			},
		})
		if e != nil {
			return
		}

		if res != nil && res.FileWorkData != nil {
			fileList = res.FileWorkData.FileList
			path = res.FileWorkData.Path
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	path, fileList, err = filework.NewLocalService().Files(dir)

	return
}

func (this_ *Worker) workFileCreate(lineNodeIdList []string, path string, isDir bool) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodFileCreate, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				Path:  path,
				IsDir: isDir,
			},
		})
		if e != nil {
			return
		}

		if res != nil && res.FileWorkData != nil {
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	err = filework.NewLocalService().Create(path, isDir)
	return
}

func (this_ *Worker) workFileRename(lineNodeIdList []string, oldPath string, newPath string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodFileRename, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				OldPath: oldPath,
				NewPath: newPath,
			},
		})
		if e != nil {
			return
		}

		if res != nil && res.FileWorkData != nil {
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	err = filework.NewLocalService().Rename(oldPath, newPath)
	return
}

func (this_ *Worker) workFileMove(lineNodeIdList []string, oldPath string, newPath string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodFileMove, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				OldPath: oldPath,
				NewPath: newPath,
			},
		})
		if e != nil {
			return
		}

		if res != nil && res.FileWorkData != nil {
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	err = filework.NewLocalService().Move(oldPath, newPath)
	return
}

func (this_ *Worker) workFileRemove(lineNodeIdList []string, path string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodFileRemove, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				Path: path,
			},
		})
		if e != nil {
			return
		}

		if res != nil && res.FileWorkData != nil {
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	return
}

func (this_ *Worker) workFileRead(lineNodeIdList []string, path string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodFileRead, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				Path: path,
			},
		})
		if e != nil {
			return
		}

		if res != nil && res.FileWorkData != nil {
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	return
}
