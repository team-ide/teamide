package modelers

import "strings"

type StepZkModel struct {
	*StepModel `json:",inline"`

	Zk                    string `json:"zk,omitempty"`     // ZK操作
	Config                string `json:"config,omitempty"` //
	Path                  string `json:"path,omitempty"`
	CreatePathIfNotExists bool   `json:"createPathIfNotExists,omitempty"`
	Value                 string `json:"value,omitempty"`
	Ephemeral             bool   `json:"ephemeral,omitempty"`
	SetVar                string `json:"setVar,omitempty"`
	SetVarType            string `json:"setVarType,omitempty"`
}

func (this_ *StepZkModel) GetType() *StepZKType {
	for _, one := range StepZKTypes {
		if strings.EqualFold(one.Value, this_.Zk) {
			return one
		}
	}
	return nil
}

type StepZKType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	StepZKTypes     []*StepZKType
	ZkGet           = appendStepZKType("get", "")
	ZkCreate        = appendStepZKType("create", "")
	ZkSet           = appendStepZKType("set", "")
	ZkStat          = appendStepZKType("stat", "")
	ZkChildren      = appendStepZKType("children", "")
	ZkDelete        = appendStepZKType("delete", "")
	ZkExists        = appendStepZKType("exists", "")
	ZkGetW          = appendStepZKType("getW", "")
	ZkWatchChildren = appendStepZKType("watchChildren", "")
)

func appendStepZKType(value string, text string) *StepZKType {
	res := &StepZKType{
		Value: value,
		Text:  text,
	}
	StepZKTypes = append(StepZKTypes, res)
	return res
}

var (
	docTemplateStepZkName = "step_zk"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepZkName,
		Fields: []*docTemplateField{
			{Name: "zk", Comment: "ZK操作"},
			{Name: "config", Comment: ""},
			{Name: "path", Comment: "路径"},
			{Name: "createPathIfNotExists", Comment: "如果路径不存在，则创建"},
			{Name: "value", Comment: "操作的Value"},
			{Name: "ephemeral", Comment: "是否是临时"},
			{Name: "setVar", Comment: "设置变量"},
			{Name: "setVarType", Comment: "设置变量类型"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepZkModel{}
		},
	})
}
