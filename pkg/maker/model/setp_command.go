package model

type StepCommandModel struct {
	*StepModel

	Command string `json:"command,omitempty"` // 执行命令
}
