package filework

import "io"

type ConfirmInfo struct {
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

type Service interface {
	InitClient() (err error)
	Write(path string, reader io.Reader, onDo func(fileCount *int, fileSize *int64), confirmInfo *ConfirmInfo) (err error)
	Read(path string, writer io.Writer, onDo func(fileCount *int, fileSize *int64)) (err error)
	Rename(oldName string, newName string) (err error)
	Move(fromPath string, fromService Service, toPath string, onDo func(fileCount *int, fileSize *int64), confirmInfo *ConfirmInfo) (err error)
	Copy(fromPath string, fromService Service, toPath string, onDo func(fileCount *int, fileSize *int64), confirmInfo *ConfirmInfo) (err error)
	Remove(path string, onDo func(fileCount *int)) (err error)
	Count(path string, onDo func(fileCount *int)) (err error)
	CountSize(path string, onDo func(fileCount *int, fileSize *int64)) (err error)
}
