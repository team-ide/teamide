package model

type StepCacheModel struct {
	*StepModel

	Cache string `json:"cache,omitempty"` // 本地缓存
}
