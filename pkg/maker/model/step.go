package model

type StepModel struct {
	Name    string `json:"name,omitempty"`    // 名称，同一个服务中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
	If      string `json:"if,omitempty"`      // 条件script，不填写或函数执行为true、1则为真，其它将跳过该阶段执行

	Steps []interface{} `json:"steps,omitempty"` // 阶段

	Return string `json:"return,omitempty"` // 返回值变量
}

var (
	docTemplateStepName = "step"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:    docTemplateStepName,
		Comment: "服务文件，该文件用于动作处理，如数据库、redis、文件等地方",
		Fields: []*docTemplateField{
			{
				Name:    "name",
				Comment: "结构体名称",
			},
			{
				Name:    "comment",
				Comment: "结构体说明",
			},
			{
				Name:    "note",
				Comment: "结构体源码注释",
			},
			{
				Name:    "if",
				Comment: "条件script，不填写或函数执行为true、1则为真，其它将跳过该阶段执行",
			},

			// 此处添加各个阶段
			{
				Sons: []*docTemplateSon{
					{
						MatchKey:   "cache",
						StructName: docTemplateStepCacheName,
						newModel: func() interface{} {
							return &StepCacheModel{}
						},
					},
					{
						MatchKey:   "command",
						StructName: docTemplateStepCommandName,
						newModel: func() interface{} {
							return &StepCommandModel{}
						},
					},
					{
						MatchKey:   "db",
						StructName: docTemplateStepDbName,
						newModel: func() interface{} {
							return &StepDbModel{}
						},
					},
					{
						MatchKey:   "error",
						StructName: docTemplateStepErrorName,
						newModel: func() interface{} {
							return &StepErrorModel{}
						},
					},
					{
						MatchKey:   "es",
						StructName: docTemplateStepEsName,
						newModel: func() interface{} {
							return &StepEsModel{}
						},
					},
					{
						MatchKey:   "file",
						StructName: docTemplateStepFileName,
						newModel: func() interface{} {
							return &StepFileModel{}
						},
					},
					{
						MatchKey:   "http",
						StructName: docTemplateStepHttpName,
						newModel: func() interface{} {
							return &StepHttpModel{}
						},
					},
					{
						MatchKey:   "lock",
						StructName: docTemplateStepLockName,
						newModel: func() interface{} {
							return &StepLockModel{}
						},
					},
					{
						MatchKey:   "mq",
						StructName: docTemplateStepMqName,
						newModel: func() interface{} {
							return &StepMqModel{}
						},
					},
					{
						MatchKey:   "redis",
						StructName: docTemplateStepRedisName,
						newModel: func() interface{} {
							return &StepRedisModel{}
						},
					},
					{
						MatchKey:   "script",
						StructName: docTemplateStepScriptName,
						newModel: func() interface{} {
							return &StepScriptModel{}
						},
					},
					{
						MatchKey:   "service",
						StructName: docTemplateStepServiceName,
						newModel: func() interface{} {
							return &StepServiceModel{}
						},
					},
					{
						MatchKey:   "var",
						StructName: docTemplateStepVarName,
						newModel: func() interface{} {
							return &StepVarModel{}
						},
					},
					{
						MatchKey:   "zk",
						StructName: docTemplateStepZkName,
						newModel: func() interface{} {
							return &StepZkModel{}
						},
					},
				},
				defaultNewModel: func() interface{} {
					return &StepModel{}
				},
			},

			{
				Name:       "steps",
				Comment:    "阶段",
				IsList:     true,
				StructName: docTemplateStepName,
			},
			{
				Name:    "return",
				Comment: "返回值",
			},
		},
		defaultNewModel: func() interface{} {
			return &StepModel{}
		},
	})
}
