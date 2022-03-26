package model

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"teamide/pkg/application/base"
)

type FileInfo struct {
	Name         string                `json:"name,omitempty"`         // 文件名称
	Type         string                `json:"type,omitempty"`         // 文件类型
	Dir          string                `json:"dir,omitempty"`          // 目录
	Path         string                `json:"path,omitempty"`         // 路径，保存或传入的相对路径
	Size         int64                 `json:"size,omitempty"`         // 文件大小
	AbsolutePath string                `json:"absolutePath,omitempty"` // 绝对路径
	FileHeader   *multipart.FileHeader `json:"-"`                      // 上传的文件
	File         *os.File              `json:"-"`                      // 读取的文件
}

func (this_ *FileInfo) OpenFileHeader() (multipart.File, error) {
	return this_.FileHeader.Open()
}

func (this_ *FileInfo) GetReader() (reader io.Reader, err error) {
	if this_.FileHeader != nil {
		var file multipart.File
		file, err = this_.FileHeader.Open()
		if err != nil {
			return
		}
		reader = file
	} else if this_.File != nil {
		reader = this_.File
	} else if this_.AbsolutePath != "" {
		reader, err = os.Open(this_.AbsolutePath)
		if err != nil {
			return
		}
	} else if this_.Path != "" {
		reader, err = os.Open(this_.Path)
		if err != nil {
			return
		}
	}
	return
}

func getActionStepFileSaveByMap(value interface{}) (step ActionStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["fileSave"] == nil {
			return
		}
		switch data := v["fileSave"].(type) {
		case map[string]interface{}:
			v["fileSave"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to file save error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to action step file save error"))
	}
	if err != nil {
		return
	}
	step = &ActionStepFileSave{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}

func getActionStepFileGetByMap(value interface{}) (step ActionStep, err error) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if v["fileGet"] == nil {
			return
		}
		switch data := v["fileGet"].(type) {
		case map[string]interface{}:
			v["fileGet"] = data
		default:
			err = errors.New(fmt.Sprint("[", v, "] to file get error"))
		}
	default:
		err = errors.New(fmt.Sprint("[", v, "] to action step file get error"))
	}
	if err != nil {
		return
	}
	step = &ActionStepFileGet{}
	err = base.ToBean([]byte(base.ToJSON(value)), step)
	return
}
