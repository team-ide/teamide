package model

type StepRedisModel struct {
	*StepModel `json:",inline"`

	Redis      string `json:"redis,omitempty"` // redis操作
	Key        string `json:"key,omitempty"`
	Field      string `json:"field,omitempty"`
	Value      string `json:"value,omitempty"`
	Expire     string `json:"expire,omitempty"`
	SetVar     string `json:"setVar,omitempty"`
	SetVarType string `json:"setVarType,omitempty"`
}

var (
	docTemplateStepRedisName = "step_redis"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepRedisName,
		Fields: []*docTemplateField{
			{Name: "redis", Comment: "redis操作"},
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
