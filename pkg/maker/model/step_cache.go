package model

type StepCacheModel struct {
	*StepModel `json:",inline"`

	Cache      string `json:"cache,omitempty"` // 本地缓存
	Key        string `json:"key,omitempty"`
	Value      string `json:"value,omitempty"`
	SetVar     string `json:"setVar,omitempty"`
	SetVarType string `json:"setVarType,omitempty"`
}

var (
	docTemplateStepCacheName = "step_cache"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepCacheName,
		Fields: []*docTemplateField{
			{
				Name:    "cache",
				Comment: "本地缓存操作",
			},
			{
				Name:    "key",
				Comment: "操作的Key",
			},
			{
				Name:    "value",
				Comment: "操作的Value",
			},
			{
				Name:    "setVar",
				Comment: "设置变量",
			},
			{
				Name:    "setVarType",
				Comment: "设置变量类型",
			},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepCacheModel{}
		},
	})
}
