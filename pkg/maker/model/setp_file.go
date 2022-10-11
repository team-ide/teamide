package model

type StepFileModel struct {
	*StepModel `json:",inline"`

	File string `json:"file,omitempty"` // 文件
	Path string `json:"path,omitempty"` // 路径，用于文件、ZK等操作

}

var (
	docTemplateStepFileName = "step_file"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:   docTemplateStepFileName,
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		Fields: []*docTemplateField{
			{
				Name:    "file",
				Comment: "文件操作",
			},
			{
				Name:    "path",
				Comment: "文件路径",
			},
		},
	})
}
