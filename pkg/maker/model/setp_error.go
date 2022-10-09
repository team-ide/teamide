package model

type StepErrorModel struct {
	*StepModel

	Error string `json:"error,omitempty"` // 异常
}
