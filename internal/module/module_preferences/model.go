package module_preferences

import "time"

const (
	// ModulePreferences 偏好模块
	ModulePreferences = "preferences"
	// TablePreferences 偏好表
	TablePreferences = "TM_PREFERENCES"
)

// PreferencesModel 偏好模型
type PreferencesModel struct {
	UserId     int64     `json:"userId,omitempty"`
	Name       string    `json:"name,omitempty"`
	Comment    string    `json:"comment,omitempty"`
	Option     string    `json:"option,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
}

func (entity *PreferencesModel) GetTableName() string {
	return TablePreferences
}

func (entity *PreferencesModel) GetPKColumnName() string {
	return ""
}
