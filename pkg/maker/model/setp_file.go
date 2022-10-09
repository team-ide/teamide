package model

type StepFileModel struct {
	*StepModel

	File string `json:"file,omitempty"` // 文件
	Path string `json:"path,omitempty"` // 路径，用于文件、ZK等操作

}
