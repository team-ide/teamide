package util

import (
	"errors"
	"io"
	"os"
	"sort"
	"strings"
)

type FileInfo struct {
	Name     string `json:"name,omitempty"`
	IsDir    bool   `json:"isDir,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Path     string `json:"path,omitempty"`
	ModTime  int64  `json:"modTime,omitempty"`
	FileMode string `json:"fileMode,omitempty"`
}

type FileRenameProgress struct {
	WaitCall     bool  `json:"-"`
	StartTime    int64 `json:"startTime"`
	EndTime      int64 `json:"endTime"`
	Timestamp    int64 `json:"timestamp"`
	Count        int64 `json:"count"`
	Size         int64 `json:"size"`
	SuccessCount int64 `json:"successCount"`
	Error        error
}

type FileRemoveProgress struct {
	WaitCall     bool  `json:"-"`
	StartTime    int64 `json:"startTime"`
	EndTime      int64 `json:"endTime"`
	Timestamp    int64 `json:"timestamp"`
	Count        int64 `json:"count"`
	Size         int64 `json:"size"`
	SuccessCount int64 `json:"successCount"`
	Error        error
}

type FileCopyProgress struct {
	WaitCall     bool         `json:"-"`
	StartTime    int64        `json:"startTime"`
	EndTime      int64        `json:"endTime"`
	Timestamp    int64        `json:"timestamp"`
	Size         int64        `json:"size"`
	SuccessSize  int64        `json:"successSize"`
	Count        int64        `json:"count"`
	SuccessCount int64        `json:"successCount"`
	Copying      *FileCopying `json:"copying,omitempty"`
	Error        error
}

type FileCopying struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	SuccessSize int64  `json:"successSize"`
	StartTime   int64  `json:"startTime"`
	EndTime     int64  `json:"endTime"`
	Timestamp   int64  `json:"timestamp"`
}

type FileUploadProgress struct {
	WaitCall     bool           `json:"-"`
	StartTime    int64          `json:"startTime"`
	EndTime      int64          `json:"endTime"`
	Timestamp    int64          `json:"timestamp"`
	Size         int64          `json:"size"`
	SuccessSize  int64          `json:"successSize"`
	Count        int64          `json:"count"`
	SuccessCount int64          `json:"successCount"`
	Uploading    *FileUploading `json:"uploading,omitempty"`
	Error        error
}

type FileDownloadProgress struct {
	WaitCall     bool  `json:"-"`
	StartTime    int64 `json:"startTime"`
	EndTime      int64 `json:"endTime"`
	Timestamp    int64 `json:"timestamp"`
	Size         int64 `json:"size"`
	SuccessSize  int64 `json:"successSize"`
	Count        int64 `json:"count"`
	SuccessCount int64 `json:"successCount"`
	Error        error
}

type FileUploading struct {
	StartTime   int64  `json:"startTime"`
	EndTime     int64  `json:"endTime"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	SuccessSize int64  `json:"successSize"`
}

type FileConfirmInfo struct {
	ConfirmId   string `json:"confirmId,omitempty"`
	IsConfirm   bool   `json:"isConfirm,omitempty"`
	Confirm     string `json:"confirm,omitempty"`
	Path        string `json:"path,omitempty"`
	Name        string `json:"name,omitempty"`
	IsFileExist bool   `json:"isFileExist,omitempty"`
	IsOk        bool   `json:"isOk,omitempty"`
	IsCancel    bool   `json:"isCancel,omitempty"`
	WorkerId    string `json:"workerId,omitempty"`
}

func LocalLoadFiles(path string) (files []os.FileInfo, err error) {
	ds, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, d := range ds {
		var file os.FileInfo
		file, err = os.Lstat(path + "/" + d.Name())
		if err != nil {
			return
		}
		files = append(files, file)
	}
	return
}

func LoadFiles(
	lstat func(string) (os.FileInfo, error),
	loadFiles func(string) ([]os.FileInfo, error),
	dir string,
) (files []*FileInfo, err error) {
	files = []*FileInfo{
		{
			Name:  "..",
			IsDir: true,
		},
	}

	fileInfo, err := lstat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}
		return
	}

	if !fileInfo.IsDir() {
		err = errors.New("路径[" + dir + "]不是目录")
		return
	}

	fs, err := loadFiles(dir)
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
		f := fMap[one]
		if err != nil {
			return
		}
		ModTime := f.ModTime()
		files = append(files, &FileInfo{
			Name:     one,
			IsDir:    true,
			ModTime:  GetTimeTime(ModTime),
			FileMode: f.Mode().String(),
		})
	}
	for _, one := range fileNames {
		f := fMap[one]
		if err != nil {
			return
		}
		ModTime := f.ModTime()
		files = append(files, &FileInfo{
			Name:     one,
			Size:     f.Size(),
			ModTime:  GetTimeTime(ModTime),
			FileMode: f.Mode().String(),
		})
	}

	return
}

func FileRemove(
	lstat func(string) (os.FileInfo, error),
	loadFiles func(string) ([]os.FileInfo, error),
	remove func(string) error,
	path string,
	progress *FileRemoveProgress,
) {
	var err error
	defer func() {
		progress.Error = err
		progress.EndTime = GetNowTime()
	}()
	progress.StartTime = GetNowTime()

	progress.Count, progress.Size, err = LoadFileCount(lstat, loadFiles, path, func(count, size int64) {
		progress.Count = count
		progress.Size = size
	})
	if err != nil {
		return
	}

	err = fileRemoveAll(lstat, loadFiles, remove, path, progress)
	if err != nil {
		return
	}
	return
}

func fileRemoveAll(
	lstat func(string) (os.FileInfo, error),
	loadFiles func(string) ([]os.FileInfo, error),
	remove func(string) error,
	path string,
	progress *FileRemoveProgress,
) (err error) {
	var isDir bool

	var info os.FileInfo
	info, err = lstat(path)
	if err != nil {
		return
	}
	isDir = info.IsDir()

	if isDir {
		var files []os.FileInfo
		files, err = loadFiles(path)
		if err != nil {
			return
		}
		for _, f := range files {
			err = fileRemoveAll(lstat, loadFiles, remove, path+"/"+f.Name(), progress)
			if err != nil {
				return
			}
		}

	}
	err = remove(path)
	if err != nil {
		return
	}
	progress.SuccessCount++
	return
}

func LoadFileCount(
	lstat func(string) (os.FileInfo, error),
	loadFiles func(string) ([]os.FileInfo, error),
	path string,
	onLoad func(count, size int64),
) (fileCount int64, fileSize int64, err error) {
	var isDir bool

	var thisFileSize int64
	var info os.FileInfo
	info, err = lstat(path)
	if err != nil {
		return
	}
	isDir = info.IsDir()
	if !isDir {
		thisFileSize = info.Size()
	}

	fileCount++
	fileSize += thisFileSize
	onLoad(fileCount, fileSize)
	if isDir {
		var files []os.FileInfo
		files, err = loadFiles(path)
		if err != nil {
			return
		}

		for _, f := range files {
			var fileCount_ int64
			var fileSize_ int64
			fileCount_, fileSize_, err = LoadFileCount(lstat, loadFiles, path+"/"+f.Name(), onLoad)
			if err != nil {
				return
			}
			fileCount += fileCount_
			fileSize += fileSize_
			onLoad(fileCount, fileSize)
		}
	}
	return
}

func LocalFileOpen(s string) (io.Reader, error) {
	return os.Open(s)
}

func LocalFileWrite(s string) (io.Writer, error) {
	return os.Create(s)
}

func FileRead(
	lstat func(string) (os.FileInfo, error),
	fileReader func(string) (io.Reader, error),
	path string,
	maxRead int64,
) (size int64, text string, err error) {
	fileInfo, err := lstat(path)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		err = errors.New("路径[" + path + "]是目录")
		return
	}
	size = fileInfo.Size()
	if size > maxRead {
		err = errors.New("文件过大，暂不支持读取")
		return
	}
	f, err := fileReader(path)
	if err != nil {
		return
	}
	defer func() {
		closer, ok := f.(io.Closer)
		if ok {
			_ = closer.Close()
		}
	}()
	bs, err := io.ReadAll(f)
	if err != nil {
		return
	}
	if len(bs) > 0 {
		text = string(bs)
	}
	return
}

func FileWrite(
	lstat func(string) (os.FileInfo, error),
	fileWriter func(string) (io.Writer, error),
	path string,
	bytes []byte,
) (err error) {
	fileInfo, err := lstat(path)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		err = errors.New("路径[" + path + "]是目录")
		return
	}
	f, err := fileWriter(path)
	if err != nil {
		return
	}
	defer func() {
		closer, ok := f.(io.Closer)
		if ok {
			_ = closer.Close()
		}
	}()
	_, err = f.Write(bytes)
	if err != nil {
		return
	}
	return
}

func FileUpload(
	lstat func(string) (os.FileInfo, error),
	mkdirAll func(string) error,
	fileReader func() (io.Reader, error),
	fileWriter func(string) (io.Writer, error),
	fileSize int64,
	fileName string,
	filePath string,
	callConfirm func(*FileConfirmInfo) (*FileConfirmInfo, error),
	progress *FileUploadProgress,
) {
	var err error
	defer func() {
		progress.Error = err
		progress.EndTime = GetNowTime()
	}()
	progress.StartTime = GetNowTime()
	progress.Count = 1
	progress.Size = fileSize

	pathDir := filePath[0:strings.LastIndex(filePath, "/")]

	_, err = lstat(pathDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = mkdirAll(pathDir)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	_, err = lstat(filePath)
	if err == nil {
		progress.WaitCall = true
		defer func() {
			progress.WaitCall = false
		}()
		confirmInfo := &FileConfirmInfo{
			IsFileExist: true,
			Path:        filePath,
			Name:        fileName,
		}
		var res *FileConfirmInfo
		res, err = callConfirm(confirmInfo)
		if err != nil {
			return
		}
		if res.IsCancel {
			progress.SuccessCount++
			progress.Size -= fileSize
			return
		}
		progress.WaitCall = false
	}

	write, err := fileWriter(filePath)
	if err != nil {
		return
	}
	defer func() {
		closer, ok := write.(io.Closer)
		if ok {
			_ = closer.Close()
		}
	}()

	read, err := fileReader()
	if err != nil {
		return
	}
	defer func() {
		closer, ok := read.(io.Closer)
		if ok {
			_ = closer.Close()
		}
	}()

	err = CopyBytes(write, read, func(readSize int64, writeSize int64) {
		progress.SuccessSize += writeSize
	})
	if err != nil {
		return
	}

	progress.SuccessCount++

	return
}

func FileCopy(
	fromLstat func(string) (os.FileInfo, error),
	fromLoadFiles func(string) ([]os.FileInfo, error),
	fromReader func(string) (io.Reader, error),
	fromPath string,
	toLstat func(string) (os.FileInfo, error),
	toMkdirAll func(string) error,
	toWriter func(string) (io.Writer, error),
	toPath string,
	callConfirm func(*FileConfirmInfo) (*FileConfirmInfo, error),
	progress *FileCopyProgress,
) {
	var err error
	defer func() {
		progress.Error = err
		progress.EndTime = GetNowTime()
	}()
	progress.StartTime = GetNowTime()

	progress.Count, progress.Size, err = LoadFileCount(fromLstat, fromLoadFiles, fromPath, func(count, size int64) {
		progress.Count = count
		progress.Size = size
	})
	if err != nil {
		return
	}
	err = fileCopyAll(fromLstat, fromLoadFiles, fromReader, fromPath, toLstat, toMkdirAll, toWriter, toPath, callConfirm, progress)
	if err != nil {
		return
	}
	return
}

func fileCopyAll(
	fromLstat func(string) (os.FileInfo, error),
	fromLoadFiles func(string) ([]os.FileInfo, error),
	fromReader func(string) (io.Reader, error),
	fromPath string,
	toLstat func(string) (os.FileInfo, error),
	toMkdirAll func(string) error,
	toWriter func(string) (io.Writer, error),
	toPath string,
	callConfirm func(*FileConfirmInfo) (*FileConfirmInfo, error),
	progress *FileCopyProgress,
) (err error) {

	var isDir bool
	var fileName string
	var fileSize int64

	var info os.FileInfo
	info, err = fromLstat(fromPath)
	if err != nil {
		return
	}
	isDir = info.IsDir()
	fileName = info.Name()
	if !isDir {
		fileSize = info.Size()
	}

	if isDir {
		progress.SuccessCount++
		var files []os.FileInfo
		files, err = fromLoadFiles(fromPath)
		if err != nil {
			return
		}

		for _, f := range files {
			err = fileCopyAll(fromLstat, fromLoadFiles, fromReader, fromPath+"/"+f.Name(), toLstat, toMkdirAll, toWriter, toPath+"/"+f.Name(), callConfirm, progress)
			if err != nil {
				return
			}
		}

	} else {

		var isExist bool

		_, err = toLstat(toPath)
		if err == nil {
			isExist = true
		}
		if isExist {
			progress.WaitCall = true
			defer func() {
				progress.WaitCall = false
			}()
			confirmInfo := &FileConfirmInfo{
				IsFileExist: true,
				Path:        toPath,
			}
			var res *FileConfirmInfo
			res, err = callConfirm(confirmInfo)
			if err != nil {
				return
			}
			if res.IsCancel {
				progress.Size -= fileSize
				progress.SuccessCount++
				return
			}
			progress.WaitCall = false
		}

		var fromIoReader io.Reader
		fromIoReader, err = fromReader(fromPath)

		if err != nil {
			return
		}

		defer func() {
			closer, ok := fromIoReader.(io.Closer)
			if ok {
				_ = closer.Close()
			}
		}()
		var toIoWriter io.Writer

		pathDir := toPath[0:strings.LastIndex(toPath, "/")]
		_, err = toLstat(pathDir)
		if err != nil {
			if os.IsNotExist(err) {
				err = toMkdirAll(pathDir)
				if err != nil {
					return
				}
			} else {
				return
			}
		}

		toIoWriter, err = toWriter(toPath)
		if err != nil {
			return
		}

		defer func() {
			closer, ok := toIoWriter.(io.Closer)
			if ok {
				_ = closer.Close()
			}
		}()

		Copying := &FileCopying{}
		Copying.Name = fileName
		Copying.StartTime = GetNowTime()
		Copying.Size = fileSize
		progress.Copying = Copying
		err = CopyBytes(toIoWriter, fromIoReader, func(readSize int64, writeSize int64) {
			progress.SuccessSize += writeSize
			Copying.SuccessSize += writeSize
		})
		if err != nil {
			return
		}

		progress.SuccessCount++
	}

	return
}
