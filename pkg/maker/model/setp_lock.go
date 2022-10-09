package model

type StepLockModel struct {
	*StepModel

	Lock     string `json:"lock,omitempty"`     // 锁
	LockType string `json:"lockType,omitempty"` // 锁类型
	Unlock   string `json:"unlock,omitempty"`   // 解锁
}
