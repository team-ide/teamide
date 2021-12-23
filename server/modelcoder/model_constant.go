package modelcoder

type ConstantModel struct {
	Name                string `json:"name,omitempty"`                // 名称，同一个应用中唯一
	Type                string `json:"type,omitempty"`                // 类型 string int long  byte
	Value               string `json:"value,omitempty"`               // 值
	EnvironmentVariable string `json:"environmentVariable,omitempty"` // 环境变量  优先取环境变量中的值
}
