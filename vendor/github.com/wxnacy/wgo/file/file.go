package file

import (
    "os"
    "strings"
    "io/ioutil"
    "errors"
)

// 判断文件路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return Exists(path) && !IsDir(path)
}

func SaveFile(path, content string) error{
    if strings.Contains(path, "/") {

        paths := strings.Split(path, "/")
        dir := strings.Join(paths[0:len(paths) - 1], "/")
        if !IsDir(dir) {
            err := os.MkdirAll(dir, os.ModePerm)
            if err != nil {
                return err
            }
        }
    }
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    f.WriteString(content)
    defer f.Close()

    return nil
}

func ReadFile(path string) ( string, error) {
    if IsFile(path) {
        d, err := ioutil.ReadFile(path)
        if err != nil {
            return "", err
        }
        return string(d), nil
    }
    return "", errors.New(path + "is not exists")
}

