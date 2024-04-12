package modelers

import "strings"

type StepMqModel struct {
	*StepModel `json:",inline"`

	Config string `json:"config,omitempty"` //
	Mq     string `json:"mq,omitempty"`     // MQ操作
	Topic  string `json:"topic,omitempty"`  // MQ topic

}

func (this_ *StepMqModel) GetType() *StepMqType {
	for _, one := range StepMqTypes {
		if strings.EqualFold(one.Value, this_.Mq) {
			return one
		}
	}
	return nil
}

type StepMqType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	StepMqTypes []*StepMqType
	MqGet       = appendStepMqType("get", "")
)

func appendStepMqType(value string, text string) *StepMqType {
	res := &StepMqType{
		Value: value,
		Text:  text,
	}
	StepMqTypes = append(StepMqTypes, res)
	return res
}

var (
	docTemplateStepMqName = "step_mq"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepMqName,
		Fields: []*docTemplateField{
			{Name: "mq", Comment: "MQ操作"},
			{Name: "topic", Comment: "MQ主题"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepMqModel{}
		},
	})
}
