package model

type StepCacheModel struct {
	*StepModel `json:",inline"`

	Cache string `json:"cache,omitempty"` // 本地缓存
}

var (
	docTemplateStepCacheName = "step_cache"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:   docTemplateStepCacheName,
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
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
		},
	})
}
