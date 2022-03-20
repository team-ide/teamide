package util

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GetAbsolutePath(path string) (absolutePath string) {
	var abs string
	abs, _ = filepath.Abs(path)

	absolutePath = filepath.ToSlash(abs)
	return
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func LoadDirFiles(fileMap map[string][]byte, dir string) (err error) {
	var exist bool
	exist, err = PathExists(dir)
	if err != nil {
		return
	}
	if !exist {
		return
	}
	//获取当前目录下的所有文件或目录信息
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {

		} else {
			var abs string
			abs, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			fileAbsolutePath := filepath.ToSlash(abs)
			name := strings.TrimPrefix(fileAbsolutePath, dir)
			name = strings.TrimPrefix(name, "/")
			var f *os.File
			f, err = os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			var bs []byte
			bs, err = io.ReadAll(f)
			if err != nil {
				return err
			}
			fileMap[name] = bs
		}
		return nil
	})
	return
}

func LoadDirFilenames(filenames *[]string, dir string) (err error) {
	var exist bool
	exist, err = PathExists(dir)
	if err != nil {
		return
	}
	if !exist {
		return
	}
	//获取当前目录下的所有文件或目录信息
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {

		} else {
			var abs string
			abs, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			fileAbsolutePath := filepath.ToSlash(abs)
			name := strings.TrimPrefix(fileAbsolutePath, dir)
			name = strings.TrimPrefix(name, "/")
			*filenames = append(*filenames, name)
		}
		return nil
	})
	sort.Strings(*filenames)
	return
}

func ReadFile(filename string) (bs []byte, err error) {
	var f *os.File
	var exists bool
	exists, err = PathExists(filename)
	if err != nil {
		return
	}
	if !exists {
		return
	} else {
		f, err = os.Open(filename)
	}
	if err != nil {
		return
	}
	defer f.Close()
	bs, err = io.ReadAll(f)
	if err != nil {
		return
	}
	return
}

func WriteFile(filename string, bs []byte) (err error) {
	var f *os.File
	var exists bool
	exists, err = PathExists(filename)
	if err != nil {
		return
	}
	if !exists {
		f, err = os.Create(filename)
	} else {
		f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	}
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(bs)

	if err != nil {
		return
	}
	return
}
