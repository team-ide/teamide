package modelers

import "strings"

type StepFileModel struct {
	*StepModel `json:",inline"`

	File string `json:"file,omitempty"` // 文件
	Path string `json:"path,omitempty"` // 路径，用于文件、ZK等操作

}

func (this_ *StepFileModel) GetType() *StepFileType {
	for _, one := range StepFileTypes {
		if strings.EqualFold(one.Value, this_.File) {
			return one
		}
	}
	return nil
}

type StepFileType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	StepFileTypes []*StepFileType
	FileGet       = appendStepFileType("get", "")
)

func appendStepFileType(value string, text string) *StepFileType {
	res := &StepFileType{
		Value: value,
		Text:  text,
	}
	StepFileTypes = append(StepFileTypes, res)
	return res
}

var (
	docTemplateStepFileName = "step_file"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepFileName,
		Fields: []*docTemplateField{
			{Name: "file", Comment: "文件操作"},
			{Name: "path", Comment: "文件路径"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepFileModel{}
		},
	})
}
