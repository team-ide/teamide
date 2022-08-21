package node

import (
	"errors"
	"os"
	"sort"
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

func (this_ *Worker) workFiles(lineNodeIdList []string, dir string) (resDir string, fileList []*FileInfo, err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodNetProxyNewConn, &Message{
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
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	fileList = append(fileList, &FileInfo{
		Name:  "..",
		IsDir: true,
		Place: "local",
	})
	if dir == "" {
		dir, err = os.UserHomeDir()
		if err != nil {
			return
		}
	}

	dir = util.FormatPath(dir)
	if err != nil {
		return
	}
	resDir = dir

	fileInfo, err := os.Lstat(dir)
	if err != nil {
		if err == os.ErrNotExist {
			err = nil
			return
		}
		return
	}

	if !fileInfo.IsDir() {
		err = errors.New("路径[" + dir + "]不是目录")
		return
	}

	dirFiles, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var dirNames []string
	var fileNames []string

	fMap := map[string]os.DirEntry{}
	for _, f := range dirFiles {
		fName := f.Name()
		fMap[fName] = f
		if f.IsDir() {
			dirNames = append(dirNames, fName)
		} else {
			fileNames = append(fileNames, fName)
		}
	}

	sort.Strings(dirNames)
	sort.Strings(fileNames)

	for _, one := range dirNames {
		f := fMap[one]
		var fi os.FileInfo
		fi, err = f.Info()
		if err != nil {
			return
		}
		ModTime := fi.ModTime()
		fileList = append(fileList, &FileInfo{
			Name:     one,
			IsDir:    true,
			Place:    "local",
			ModTime:  util.GetTimeTime(ModTime),
			FileMode: fi.Mode().String(),
		})
	}
	for _, one := range fileNames {
		f := fMap[one]
		var fi os.FileInfo
		fi, err = f.Info()
		if err != nil {
			return
		}
		ModTime := fi.ModTime()
		fileList = append(fileList, &FileInfo{
			Name:     one,
			Size:     fi.Size(),
			Place:    "local",
			ModTime:  util.GetTimeTime(ModTime),
			FileMode: fi.Mode().String(),
		})
	}
	return
}

func (this_ *Worker) workFileRemove(lineNodeIdList []string, dir string) (resDir string, fileList []*FileInfo, err error) {
	send, err := this_.sendToNext(lineNodeIdList, "", func(listener *MessageListener) (e error) {
		res, e := this_.Call(listener, methodNetProxyNewConn, &Message{
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
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	fileList = append(fileList, &FileInfo{
		Name:  "..",
		IsDir: true,
		Place: "local",
	})
	if dir == "" {
		dir, err = os.UserHomeDir()
		if err != nil {
			return
		}
	}

	dir = util.FormatPath(dir)
	if err != nil {
		return
	}
	resDir = dir

	fileInfo, err := os.Lstat(dir)
	if err != nil {
		if err == os.ErrNotExist {
			err = nil
			return
		}
		return
	}

	if !fileInfo.IsDir() {
		err = errors.New("路径[" + dir + "]不是目录")
		return
	}

	dirFiles, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var dirNames []string
	var fileNames []string

	fMap := map[string]os.DirEntry{}
	for _, f := range dirFiles {
		fName := f.Name()
		fMap[fName] = f
		if f.IsDir() {
			dirNames = append(dirNames, fName)
		} else {
			fileNames = append(fileNames, fName)
		}
	}

	sort.Strings(dirNames)
	sort.Strings(fileNames)

	for _, one := range dirNames {
		f := fMap[one]
		var fi os.FileInfo
		fi, err = f.Info()
		if err != nil {
			return
		}
		ModTime := fi.ModTime()
		fileList = append(fileList, &FileInfo{
			Name:     one,
			IsDir:    true,
			Place:    "local",
			ModTime:  util.GetTimeTime(ModTime),
			FileMode: fi.Mode().String(),
		})
	}
	for _, one := range fileNames {
		f := fMap[one]
		var fi os.FileInfo
		fi, err = f.Info()
		if err != nil {
			return
		}
		ModTime := fi.ModTime()
		fileList = append(fileList, &FileInfo{
			Name:     one,
			Size:     fi.Size(),
			Place:    "local",
			ModTime:  util.GetTimeTime(ModTime),
			FileMode: fi.Mode().String(),
		})
	}
	return
}

func (this_ *Worker) workFileRead(lineNodeIdList []string, path string) (size int64, text string, err error) {
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
			text = res.FileWorkData.Text
		}
		return
	})
	if err != nil {
		return
	}
	if send {
		return
	}

	size, text, err = util.FileRead(os.Lstat, util.LocalFileOpen, path, 1024*1024*10)
	return
}

func (this_ *Worker) workFileRename(lineNodeIdList []string, isNew bool, isDir bool, oldPath string, newPath string) (path string, err error) {
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
	if isNew {
		if isDir {
			err = os.MkdirAll(newPath, os.ModePerm)
		} else {
			var f *os.File
			f, err = os.Create(newPath)
			defer func() {
				_ = f.Close()
			}()
		}
		if err != nil {
			return
		}

		return
	}

	_, err = os.Lstat(oldPath)
	if err != nil {
		return
	}
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return
	}
	return
}
