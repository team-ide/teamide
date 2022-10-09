package model

type StepServiceModel struct {
	*StepModel

	Service string `json:"service,omitempty"` // 调用服务
}
