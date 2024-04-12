package modelers

import "strings"

type StepCacheModel struct {
	*StepModel `json:",inline"`

	Cache      string `json:"cache,omitempty"` // 本地缓存
	Key        string `json:"key,omitempty"`
	Value      string `json:"value,omitempty"`
	SetVar     string `json:"setVar,omitempty"`
	SetVarType string `json:"setVarType,omitempty"`
}

func (this_ *StepCacheModel) GetType() *StepCacheType {
	for _, one := range StepCacheTypes {
		if strings.EqualFold(one.Value, this_.Cache) {
			return one
		}
	}
	return nil
}

type StepCacheType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	StepCacheTypes []*StepCacheType
	CacheGet       = appendStepCacheType("get", "")
)

func appendStepCacheType(value string, text string) *StepCacheType {
	res := &StepCacheType{
		Value: value,
		Text:  text,
	}
	StepCacheTypes = append(StepCacheTypes, res)
	return res
}

var (
	docTemplateStepCacheName = "step_cache"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepCacheName,
		Fields: []*docTemplateField{
			{Name: "cache", Comment: "本地缓存操作"},
			{Name: "key", Comment: "操作的Key"},
			{Name: "value", Comment: "操作的Value"},
			{Name: "setVar", Comment: "设置变量"},
			{Name: "setVarType", Comment: "设置变量类型"},
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
