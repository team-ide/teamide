package model

type ServiceModel struct {
	Name    string              `json:"name,omitempty"`    // 名称，同一个应用中唯一
	Comment string              `json:"comment,omitempty"` // 说明
	Note    string              `json:"note,omitempty"`    // 注释
	Steps   []*ServiceStepModel `json:"steps,omitempty"`   // 阶段
}
