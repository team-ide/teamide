package model

type StepLockModel struct {
	*StepModel `json:",inline"`

	Lock     string `json:"lock,omitempty"`     // 锁
	LockType string `json:"lockType,omitempty"` // 锁类型
}

var (
	docTemplateStepLockName = "step_lock"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:   docTemplateStepLockName,
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		Fields: []*docTemplateField{
			{
				Name:    "lock",
				Comment: "锁操作",
			},
			{
				Name:    "lockType",
				Comment: "锁类型",
			},
		},
	})
}
