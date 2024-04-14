package modelers

type ErrorModel struct {
	Name    string              `json:"name,omitempty"`    // 字段名称，同一个结构体中唯一
	Comment string              `json:"comment,omitempty"` // 说明
	Note    string              `json:"note,omitempty"`    // 注释
	Options []*ErrorOptionModel `json:"options,omitempty"`
}

type ErrorOptionModel struct {
	Name    string `json:"name,omitempty"`    // 字段名称，同一个结构体中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
	Code    string `json:"code,omitempty"`    // 错误码
	Msg     string `json:"msg,omitempty"`     // 错误信息
}
