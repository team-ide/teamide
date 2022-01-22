package model

type ErrorModel struct {
	Name    string `json:"name,omitempty"`                             // 名称，同一个应用中唯一
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"` // 注释说明
	Code    string `json:"code,omitempty"`                             // 错误码
	Msg     string `json:"msg,omitempty"`                              // 错误信息
}

func TextToErrorModel(namePath string, text string) (model *ErrorModel, err error) {
	var name string
	model = &ErrorModel{}
	name, err = TextToModel(namePath, text, model)
	if err != nil {
		return
	}
	model.Name = name
	return
}