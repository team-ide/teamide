package modelers

import "strings"

type StepRedisModel struct {
	*StepModel `json:",inline"`

	Redis      string `json:"redis,omitempty"`  // redis操作
	Config     string `json:"config,omitempty"` //
	Key        string `json:"key,omitempty"`
	Field      string `json:"field,omitempty"`
	Value      string `json:"value,omitempty"`
	Expire     string `json:"expire,omitempty"`
	SetVar     string `json:"setVar,omitempty"`
	SetVarType string `json:"setVarType,omitempty"`
}

func (this_ *StepRedisModel) GetType() *StepRedisType {
	for _, one := range StepRedisTypes {
		if strings.EqualFold(one.Value, this_.Redis) {
			return one
		}
	}
	return nil
}

type StepRedisType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	StepRedisTypes []*StepRedisType
	RedisGet       = appendStepRedisType("get", "")
	RedisSet       = appendStepRedisType("set", "")
)

func appendStepRedisType(value string, text string) *StepRedisType {
	res := &StepRedisType{
		Value: value,
		Text:  text,
	}
	StepRedisTypes = append(StepRedisTypes, res)
	return res
}

var (
	docTemplateStepRedisName = "step_redis"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepRedisName,
		Fields: []*docTemplateField{
			{Name: "redis", Comment: "redis操作"},
			{Name: "config", Comment: ""},
			{Name: "key", Comment: "操作的Key"}, {Name: "field", Comment: "Hash的Key"},
			{Name: "value", Comment: "操作的Value"},
			{Name: "setVar", Comment: "设置变量"},
			{Name: "setVarType", Comment: "设置变量类型"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepRedisModel{}
		},
	})
}
