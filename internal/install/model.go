package install

const (
	// ModuleInstall install模块
	ModuleInstall = "install"
	// TableInstall Install信息表
	TableInstall = "TM_INSTALL"
)

var (
	Stages []*StageModel
)

func AppendInstallStage(installStages *StageModel) {
	Stages = append(Stages, installStages)
}

type StageModel struct {
	Version string         `json:"version,omitempty"`
	Module  string         `json:"module,omitempty"`
	Stage   string         `json:"stage,omitempty"`
	Sql     *StageSqlModel `json:"sql,omitempty"`
}

type StageSqlModel struct {
	Mysql  []string `json:"mysql,omitempty"`
	Sqlite []string `json:"sqlite,omitempty"`
}
