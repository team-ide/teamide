package filework

import (
	"errors"
	"io"
	"os"
	"sort"
	"strings"
	"teamide/pkg/util"
)

type LocalService struct {
}

func (this_ *LocalService) Exist(path string) (exist bool, err error) {

	exist, err = util.PathExists(path)

	return
}

func (this_ *LocalService) Create(path string, isDir bool) (file *FileInfo, err error) {
	path = util.FormatPath(path)
	exist, err := util.PathExists(path)
	if err != nil {
		return
	}
	if exist {
		err = errors.New("路径[" + path + "]已存在")
		return
	}

	if isDir {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return
		}
	} else {
		var f *os.File
		f, err = os.Create(path)
		if err != nil {
			return
		}
		defer func() { _ = f.Close() }()
	}
	file, err = this_.File(path)
	return
}

func (this_ *LocalService) Write(path string, bytes []byte) (err error) {
	return
}

func (this_ *LocalService) WriteByReader(path string, reader io.Reader, onDo func(readSize int64, writeSize int64)) (file *FileInfo, err error) {
	path = util.FormatPath(path)

	pathDir := path[0:strings.LastIndex(path, "/")]

	_, err = os.Stat(pathDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(pathDir, os.ModePerm)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	var readSize int64
	var writeSize int64
	for {
		bs := make([]byte, 1024)
		var n int
		n, err = reader.Read(bs)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
		}
		readSize += int64(n)
		onDo(readSize, writeSize)
		n, err = f.Write(bs[:n])
		if err != nil {
			break
		}
		writeSize += int64(n)
		onDo(readSize, writeSize)
	}

	if err != nil {
		return
	}
	file, err = this_.File(path)
	return
}

func (this_ *LocalService) Read(path string) (bytes []byte, err error) {
	return
}

func (this_ *LocalService) Rename(oldPath string, newPath string) (file *FileInfo, err error) {
	oldPath = util.FormatPath(oldPath)
	newPath = util.FormatPath(newPath)

	exist, err := util.PathExists(oldPath)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("路径[" + oldPath + "]不存在")
		return
	}

	exist, err = util.PathExists(newPath)
	if err != nil {
		return
	}
	if exist {
		err = errors.New("路径[" + newPath + "]已存在")
		return
	}

	err = os.Rename(oldPath, newPath)
	if err != nil {
		return
	}
	file, err = this_.File(newPath)
	return
}

func (this_ *LocalService) Move(oldPath string, newPath string) (err error) {
	oldPath = util.FormatPath(oldPath)
	newPath = util.FormatPath(newPath)

	exist, err := util.PathExists(oldPath)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("路径[" + oldPath + "]不存在")
		return
	}

	exist, err = util.PathExists(newPath)
	if err != nil {
		return
	}
	if exist {
		err = errors.New("路径[" + newPath + "]已存在")
		return
	}

	err = os.Rename(oldPath, newPath)
	if err != nil {
		return
	}
	return
}

func (this_ *LocalService) Copy(path string, fromService Service, fromPath string, onDo func(fileCount int, fileSize int64)) (err error) {
	return
}

func (this_ *LocalService) Remove(path string, onDo func(fileCount int, removeCount int)) (err error) {
	var fileCount int
	var removeCount int

	err = removeFile(path, func() {
		fileCount++
		onDo(fileCount, removeCount)
	}, func() {
		removeCount++
		onDo(fileCount, removeCount)
	})

	return
}
func removeFile(path string, onLoad func(), onRemove func()) (err error) {
	var isDir bool

	var info os.FileInfo
	info, err = os.Stat(path)
	if err != nil {
		return
	}
	isDir = info.IsDir()

	onLoad()
	if isDir {
		var ds []os.DirEntry
		ds, err = os.ReadDir(path)
		if err != nil {
			return
		}

		for _, d := range ds {
			err = removeFile(path+"/"+d.Name(), onLoad, onRemove)
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
func (this_ *LocalService) Count(path string, onDo func(fileCount int)) (fileCount int, err error) {
	return
}

func (this_ *LocalService) CountSize(path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error) {
	return
}

func (this_ *LocalService) Files(dir string) (parentPath string, files []*FileInfo, err error) {
	parentPath = dir
	if parentPath == "" {
		parentPath, err = os.UserHomeDir()
		if err != nil {
			return
		}
	}
	parentPath = util.FormatPath(parentPath)
	if !strings.HasSuffix(parentPath, "/") {
		parentPath += "/"
	}

	files = []*FileInfo{
		{
			Name:  "..",
			Path:  parentPath + "..",
			IsDir: true,
		},
	}

	fileInfo, err := os.Stat(parentPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.New("路径[" + parentPath + "]不存在")
			return
		}
		return
	}

	if !fileInfo.IsDir() {
		err = errors.New("路径[" + parentPath + "]不是目录")
		return
	}

	fs, err := util.LocalLoadFiles(parentPath)
	if err != nil {
		return
	}
	var dirNames []string
	var fileNames []string

	fMap := map[string]os.FileInfo{}
	for _, f := range fs {
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
		fileOne := getFileInfoByStat(parentPath+one, fMap[one])
		files = append(files, fileOne)
	}
	for _, one := range fileNames {
		fileOne := getFileInfoByStat(parentPath+one, fMap[one])
		files = append(files, fileOne)
	}

	return
}

func (this_ *LocalService) File(path string) (file *FileInfo, err error) {
	file, err = getFileInfo(path)
	return
}

func getFileInfo(path string) (fileInfo *FileInfo, err error) {
	path = util.FormatPath(path)
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}
		return
	}

	fileInfo = getFileInfoByStat(path, stat)
	return
}

func getFileInfoByStat(path string, stat os.FileInfo) (fileInfo *FileInfo) {

	fileInfo = &FileInfo{
		Name:     stat.Name(),
		Path:     path,
		IsDir:    stat.IsDir(),
		ModTime:  util.GetTimeTime(stat.ModTime()),
		FileMode: stat.Mode().String(),
		Size:     stat.Size(),
	}
	return
}

func loadFileCount(path string, onCount func(count int)) (fileCount int, err error) {
	var isDir bool

	var info os.FileInfo
	info, err = os.Stat(path)
	if err != nil {
		return
	}
	isDir = info.IsDir()

	fileCount++
	onCount(fileCount)
	if isDir {
		var files []os.FileInfo
		files, err = util.LocalLoadFiles(path)
		if err != nil {
			return
		}

		for _, f := range files {
			var fileCount_ int
			fileCount_, err = loadFileCount(path+"/"+f.Name(), func(count int) {
				onCount(fileCount + count)
			})
			if err != nil {
				return
			}
			fileCount += fileCount_
			onCount(fileCount)
		}
	}
	return
}
