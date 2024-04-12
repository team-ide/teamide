package modelers

import "strings"

type StepEsModel struct {
	*StepModel `json:",inline"`

	Config string `json:"config,omitempty"` //
	Es     string `json:"es,omitempty"`     // ES操作

}

func (this_ *StepEsModel) GetType() *StepEsType {
	for _, one := range StepEsTypes {
		if strings.EqualFold(one.Value, this_.Es) {
			return one
		}
	}
	return nil
}

type StepEsType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	StepEsTypes []*StepEsType
	EsGet       = appendStepEsType("get", "")
)

func appendStepEsType(value string, text string) *StepEsType {
	res := &StepEsType{
		Value: value,
		Text:  text,
	}
	StepEsTypes = append(StepEsTypes, res)
	return res
}

var (
	docTemplateStepEsName = "step_es"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepEsName,
		Fields: []*docTemplateField{
			{Name: "es", Comment: "ES操作"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepEsModel{}
		},
	})
}
