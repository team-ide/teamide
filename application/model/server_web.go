package model

type ServerWebModel struct {
	Name        string          `json:"name,omitempty" yaml:"name,omitempty"`               // 名称，同一个应用中唯一
	Comment     string          `json:"comment,omitempty" yaml:"comment,omitempty"`         // 注释说明
	Host        string          `json:"host,omitempty" yaml:"host,omitempty"`               // 绑定地址
	Port        int             `json:"port,omitempty" yaml:"port,omitempty"`               // 端口
	ContextPath string          `json:"contextPath,omitempty" yaml:"contextPath,omitempty"` // 访问路径
	Token       *ServerWebToken `json:"token,omitempty" yaml:"token,omitempty"`             // Token
}

type ServerWebToken struct {
	Include string `json:"include,omitempty" yaml:"include,omitempty"` //
	Exclude string `json:"exclude,omitempty" yaml:"exclude,omitempty"` //

	Variables        []*VariableModel `json:"variables,omitempty" yaml:"variables,omitempty"`               // 输入变量
	Validates        []*ValidateModel `json:"validates,omitempty" yaml:"validates,omitempty"`               // 验证
	CreateAction     string           `json:"createAction,omitempty" yaml:"createAction,omitempty"`         //
	ValidateAction   string           `json:"validateAction,omitempty" yaml:"validateAction,omitempty"`     //
	VariableName     string           `json:"variableName,omitempty" yaml:"variableName,omitempty"`         //
	VariableDataType string           `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` //
}
