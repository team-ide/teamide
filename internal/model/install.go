package model

const (
	// ModuleInstall install模块
	ModuleInstall = "install"
	// TableInstall Install信息表
	TableInstall = "TM_INSTALL"
)

var (
	InstallStages []*InstallStageModel
)

func appendInstallStage(installStages *InstallStageModel) {
	InstallStages = append(InstallStages, installStages)
}

type InstallStageModel struct {
	Version string                `json:"version,omitempty"`
	Module  string                `json:"module,omitempty"`
	Stage   string                `json:"stage,omitempty"`
	Sql     *InstallStageSqlModel `json:"sql,omitempty"`
}

type InstallStageSqlModel struct {
	Mysql  []string `json:"mysql,omitempty"`
	Sqlite []string `json:"sqlite,omitempty"`
}
