package modelers

type ConstantModel struct {
	Name    string `json:"name,omitempty"`    // 字段名称，同一个结构体中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
	Type    string `json:"type,omitempty"`    // 类型
	Value   string `json:"value,omitempty"`   // 值
}
