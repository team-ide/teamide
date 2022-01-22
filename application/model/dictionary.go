package model

type DictionaryModel struct {
	Name    string `json:"name,omitempty"`                             // 名称，同一个应用中唯一
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"` // 注释说明
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
