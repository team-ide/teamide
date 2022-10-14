package model

import "strings"

type StepHttpModel struct {
	*StepModel `json:",inline"`

	Http string `json:"http,omitempty"` // HTTP操作
	Url  string `json:"url,omitempty"`  // HTTP地址
}

func (this_ *StepHttpModel) GetType() *StepHttpType {
	for _, one := range StepHttpTypes {
		if strings.EqualFold(one.Value, this_.Http) {
			return one
		}
	}
	return nil
}

type StepHttpType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	StepHttpTypes []*StepHttpType
	HttpGet       = appendStepHttpType("get", "")
)

func appendStepHttpType(value string, text string) *StepHttpType {
	res := &StepHttpType{
		Value: value,
		Text:  text,
	}
	StepHttpTypes = append(StepHttpTypes, res)
	return res
}

var (
	docTemplateStepHttpName = "step_http"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepHttpName,
		Fields: []*docTemplateField{
			{Name: "http", Comment: "HTTP操作"},
			{Name: "url", Comment: "地址"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepHttpModel{}
		},
	})
}
