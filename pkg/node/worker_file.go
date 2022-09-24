package node

import (
	"errors"
	"go.uber.org/zap"
	"os"
	"teamide/pkg/filework"
	"teamide/pkg/util"
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
	if err != nil || send {
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
	if err != nil || send {
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
	if err != nil || send {
		return
	}

	path, fileList, err = filework.NewLocalService().Files(dir)

	return
}

func (this_ *Worker) workFileCreate(lineNodeIdList []string, path string, isDir bool) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodFileCreate, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				Path:  path,
				IsDir: isDir,
			},
		})
		if e != nil {
			return
		}
		return
	})
	if err != nil || send {
		return
	}

	err = filework.NewLocalService().Create(path, isDir)
	return
}

func (this_ *Worker) workFileRename(lineNodeIdList []string, oldPath string, newPath string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodFileRename, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				OldPath: oldPath,
				NewPath: newPath,
			},
		})
		if e != nil {
			return
		}
		return
	})
	if err != nil || send {
		return
	}

	err = filework.NewLocalService().Rename(oldPath, newPath)
	return
}

func (this_ *Worker) workFileMove(lineNodeIdList []string, oldPath string, newPath string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodFileMove, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				OldPath: oldPath,
				NewPath: newPath,
			},
		})
		if e != nil {
			return
		}
		return
	})
	if err != nil || send {
		return
	}

	err = filework.NewLocalService().Move(oldPath, newPath)
	return
}

func (this_ *Worker) workFileRemove(lineNodeIdList []string, path string) (fileCount int, removeCount int, err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodFileRemove, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				Path: path,
			},
		})
		if e != nil {
			return
		}
		return
	})
	if err != nil || send {
		return
	}

	err = filework.NewLocalService().Remove(path, func(fileCount_ int, removeCount_ int) {
		fileCount = fileCount_
		removeCount = removeCount_
	})

	return
}

func (this_ *Worker) workFileRead(lineNodeIdList []string, path string, sendKey string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodFileRead, &Message{
			LineNodeIdList: lineNodeIdList,
			SendKey:        sendKey,
			FileWorkData: &FileWorkData{
				Path: path,
			},
		})
		if e != nil {
			return
		}

		return
	})
	if err != nil || send {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}

	go func() {
		defer func() { _ = f.Close() }()

		var line []string
		for i := len(lineNodeIdList) - 1; i >= 0; i-- {
			line = append(line, lineNodeIdList[i])
		}

		err = this_.workSend(line, sendKey, f.Read)
		if err != nil {
			Logger.Error("file read send error", zap.Error(err))
		}
	}()

	return
}

func (this_ *Worker) workFileWrite(lineNodeIdList []string, path string) (sendKey string, err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodFileWrite, &Message{
			LineNodeIdList: lineNodeIdList,
			FileWorkData: &FileWorkData{
				Path: path,
			},
		})
		if e != nil {
			return
		}

		if res != nil {
			sendKey = res.SendKey
		}

		return
	})
	if err != nil || send {
		return
	}

	sendKey = util.UUID()

	var file *os.File
	this_.addOnBytesCache(sendKey, &OnBytes{
		start: func() (err error) {
			file, err = os.Create(path)
			return
		},
		on: func(buf []byte) (err error) {
			if file == nil {
				err = errors.New("文件[" + path + "]未打开")
				return
			}
			_, err = file.Write(buf)
			return
		},
		end: func() (err error) {
			if file != nil {
				_ = file.Close()
			}
			return
		},
	})

	return
}

func (this_ *Worker) workSend(lineNodeIdList []string, key string, read func(p []byte) (n int, err error)) (err error) {

	err = this_.workSendBytesStart(lineNodeIdList, key)
	if err != nil {
		return
	}

	var buf = make([]byte, 1024*32)
	err = util.ReadByFunc(read, buf, func(n int) (e error) {
		//Logger.Info("workSend read", zap.Any("key", key), zap.Any("n", n), zap.Any("str", string(buf[:n])))
		if n > 0 {
			e = this_.workSendBytes(lineNodeIdList, key, buf[:n])
		}
		return
	})
	if err != nil {
		return
	}

	err = this_.workSendBytesEnd(lineNodeIdList, key)
	if err != nil {
		return
	}

	return
}

func (this_ *Worker) workSendBytesStart(lineNodeIdList []string, key string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, key, func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodSendBytesStart, &Message{
			LineNodeIdList: lineNodeIdList,
			SendKey:        key,
		})
		if e != nil {
			return
		}

		return
	})
	if err != nil || send {
		return
	}

	onBytes := this_.getOnBytesCache(key)
	if onBytes == nil {
		err = errors.New("流读取器[" + key + "]不存在")
		return
	}

	err = onBytes.start()

	return
}

func (this_ *Worker) workSendBytesEnd(lineNodeIdList []string, key string) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, key, func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodSendBytesEnd, &Message{
			LineNodeIdList: lineNodeIdList,
			SendKey:        key,
		})
		if e != nil {
			return
		}

		return
	})
	if err != nil || send {
		return
	}

	onBytes := this_.getOnBytesCache(key)
	if onBytes == nil {
		err = errors.New("流读取器[" + key + "]不存在")
		return
	}

	err = onBytes.end()
	this_.removeOnBytesCache(key)

	return
}

func (this_ *Worker) workSendBytes(lineNodeIdList []string, key string, buf []byte) (err error) {
	send, err := this_.sendToNext(lineNodeIdList, key, func(listener *MessageListener) (e error) {
		_, e = this_.Call(listener, methodSendBytes, &Message{
			LineNodeIdList: lineNodeIdList,
			SendKey:        key,
			HasBytes:       true,
			Bytes:          buf,
		})
		if e != nil {
			return
		}

		return
	})
	if err != nil || send {
		return
	}

	onBytes := this_.getOnBytesCache(key)
	if onBytes == nil {
		err = errors.New("流读取器[" + key + "]不存在")
		return
	}

	err = onBytes.on(buf)

	return
}
