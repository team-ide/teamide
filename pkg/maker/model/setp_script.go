package model

type StepScriptModel struct {
	*StepModel

	Script string `json:"script,omitempty"` // 脚本
}
