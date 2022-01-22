package util

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
