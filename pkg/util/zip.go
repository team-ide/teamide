package util

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//GZipBytes 压缩
func GZipBytes(data []byte) ([]byte, error) {

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//UGZipBytes 解压
func UGZipBytes(data []byte) ([]byte, error) {
	var in bytes.Buffer
	in.Write(data)
	r, err := gzip.NewReader(&in)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	bs, err := io.ReadAll(r)
	return bs, err
}

//Zip zip压缩 srcFile 文件路径，destZip压缩包保存路径
func Zip(srcFile string, destZip string) error {
	srcFiles := srcFile
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	archive := zip.NewWriter(zipFile)
	defer archive.Close()
	err = filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		path = strings.Replace(path, "\\", "/", -1)
		srcFiles = strings.Replace(srcFiles, "\\", "/", -1)
		header.Name = strings.Replace(path, srcFiles, "", -1)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			var file *os.File
			file, err = os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})
	return err
}

//Unzip zip解压 zipFile 压缩包地址 destDir 解压保存文件夹
func Unzip(zipFile string, destDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	dirLen := 0
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(destDir+string(filepath.Separator)+file.Name[dirLen:], 0755)
			if err != nil {
				return err
			}
			continue
		} else {
			srcFile, err := file.Open()
			if err != nil {
				return err
			}
			defer srcFile.Close()
			destFile, err := os.Create(destDir + string(filepath.Separator) + file.Name[dirLen:])
			if err != nil {
				return err
			}
			defer destFile.Close()
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
