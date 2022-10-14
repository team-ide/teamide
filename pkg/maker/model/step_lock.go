package model

import "strings"

type StepLockModel struct {
	*StepModel `json:",inline"`

	Lock     string `json:"lock,omitempty"`     // 锁
	LockType string `json:"lockType,omitempty"` // 锁类型
}

func (this_ *StepLockModel) GetType() *StepLockType {
	for _, one := range StepLockTypes {
		if strings.EqualFold(one.Value, this_.Lock) {
			return one
		}
	}
	return nil
}

type StepLockType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	StepLockTypes []*StepLockType
	LockGet       = appendStepLockType("get", "")
)

func appendStepLockType(value string, text string) *StepLockType {
	res := &StepLockType{
		Value: value,
		Text:  text,
	}
	StepLockTypes = append(StepLockTypes, res)
	return res
}

var (
	docTemplateStepLockName = "step_lock"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepLockName,
		Fields: []*docTemplateField{
			{Name: "lock", Comment: "锁操作"},
			{Name: "lockType", Comment: "锁类型"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepLockModel{}
		},
	})
}
