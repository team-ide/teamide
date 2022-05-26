package model

type DictionaryModel struct {
	Name    string              `json:"name,omitempty" yaml:"name,omitempty"`       // 名称，同一个应用中唯一
	Comment string              `json:"comment,omitempty" yaml:"comment,omitempty"` // 注释说明
	Options []*DictionaryOption `json:"options,omitempty" yaml:"options,omitempty"` // 结构体字段
}

type DictionaryOption struct {
	Text     string `json:"text,omitempty" yaml:"text,omitempty"`         // 字段名称，同一个结构体中唯一
	Value    string `json:"value,omitempty" yaml:"value,omitempty"`       // 字段名称，同一个结构体中唯一
	DataType string `json:"dataType,omitempty" yaml:"dataType,omitempty"` // 数据类型
	Comment  string `json:"comment,omitempty" yaml:"comment,omitempty"`   // 注释
}

func TextToDictionaryModel(namePath string, text string) (model *DictionaryModel, err error) {
	var name string
	model = &DictionaryModel{}
	name, err = TextToModel(namePath, text, model)
	if err != nil {
		return
	}
	model.Name = name
	return
}

func DictionaryModelToText(model *DictionaryModel) (text string, err error) {
	text, err = ToText(model)
	if err != nil {
		return
	}
	return
}
