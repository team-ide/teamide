package model

type ConstantModel struct {
	Name                string `json:"name,omitempty" yaml:"name,omitempty"`                               // 名称，同一个应用中唯一
	Comment             string `json:"comment,omitempty" yaml:"comment,omitempty"`                         // 注释说明
	DataType            string `json:"dataType,omitempty" yaml:"dataType,omitempty"`                       // 类型 string int long  byte
	Value               string `json:"value,omitempty" yaml:"value,omitempty"`                             // 值
	EnvironmentVariable string `json:"environmentVariable,omitempty" yaml:"environmentVariable,omitempty"` // 环境变量  优先取环境变量中的值
}

func TextToConstantModel(namePath string, text string) (model *ConstantModel, err error) {
	var name string
	model = &ConstantModel{}
	name, err = TextToModel(namePath, text, model)
	if err != nil {
		return
	}
	model.Name = name
	return
}

func ConstantModelToText(model *ConstantModel) (text string, err error) {
	text, err = ModelToText(model)
	if err != nil {
		return
	}
	return
}
