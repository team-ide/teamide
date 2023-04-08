package module_setting

const (
	// ModuleSetting 设置模块
	ModuleSetting = "setting"
	// TableSetting 设置表
	TableSetting        = "TM_SETTING"
	TableSettingComment = "设置"
)

// SettingModel 设置模型，和设置表对应
type SettingModel struct {
	Name       string `json:"name,omitempty"`
	Value      string `json:"value,omitempty"`
	CreateTime int64  `json:"createTime,omitempty"`
	UpdateTime int64  `json:"updateTime,omitempty"`
}
