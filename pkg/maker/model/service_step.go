package model

type ServiceStepModel struct {
	Name    string              `json:"name,omitempty"`    // 名称，同一个服务中唯一
	Comment string              `json:"comment,omitempty"` // 说明
	Note    string              `json:"note,omitempty"`    // 注释
	If      string              `json:"if,omitempty"`      // 条件script，不填写或函数执行为true、1则为真，其它将跳过该阶段执行
	Steps   []*ServiceStepModel `json:"steps,omitempty"`   // 阶段
}
