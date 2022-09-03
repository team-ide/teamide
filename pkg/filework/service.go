package filework

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

type FileInfo struct {
	Name     string `json:"name,omitempty"`
	IsDir    bool   `json:"isDir,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Path     string `json:"path,omitempty"`
	ModTime  int64  `json:"modTime,omitempty"`
	FileMode string `json:"fileMode,omitempty"`
}

type Service interface {
	Exist(path string) (exist bool, err error)
	Create(path string, isDir bool) (file *FileInfo, err error)
	Write(path string, bytes []byte) (err error)
	Read(path string) (bytes []byte, err error)
	Rename(oldPath string, newPath string) (file *FileInfo, err error)
	Move(fromPath string, fromService Service, toPath string, onDo func(fileCount int, fileSize int64)) (err error)
	Copy(fromPath string, fromService Service, toPath string, onDo func(fileCount int, fileSize int64)) (err error)
	Remove(path string, onDo func(fileCount int, removeCount int)) (err error)
	Count(path string, onDo func(fileCount int)) (fileCount int, err error)
	CountSize(path string, onDo func(fileCount int, fileSize int64)) (fileCount int, fileSize int64, err error)
	Files(dir string) (parentPath string, files []*FileInfo, err error)
	File(path string) (file *FileInfo, err error)
}
