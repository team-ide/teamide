package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"teamide/internal/context"
	"teamide/internal/install"
	"teamide/internal/module/module_id"
	"teamide/internal/module/module_login"
	"teamide/internal/module/module_power"
	"teamide/internal/module/module_register"
	"teamide/internal/module/module_toolbox"
	"teamide/internal/module/module_user"
	"time"
)

// NewInstallService 根据库配置创建InstallService
func NewInstallService(ServerContext *context.ServerContext) (res *InstallService) {
	res = &InstallService{
		ServerContext: ServerContext,
	}
	return
}

// InstallService 安装程序服务
type InstallService struct {
	*context.ServerContext
}

// Check 检测基础配置是否完整，如基础表
func (this_ *InstallService) Check() (err error) {
	var isExist bool
	switch strings.ToLower(this_.DatabaseWorker.GetConfig().Type) {
	case "mysql":
		sql := "SELECT count(1) FROM information_schema.TABLES WHERE table_schema='" + this_.DatabaseWorker.GetConfig().Database + "' AND table_name ='" + install.TableInstall + "'"
		var count int64
		count, err = this_.DatabaseWorker.Count(sql, []interface{}{})
		if err != nil {
			return
		}
		isExist = count > 0
	case "sqlite":
		sql := `SELECT count(1) FROM sqlite_master WHERE type="table" AND name="` + install.TableInstall + `" `
		var count int64
		count, err = this_.DatabaseWorker.Count(sql, []interface{}{})
		if err != nil {
			return
		}
		isExist = count > 0
	}
	if !isExist {
		installSql := `
CREATE TABLE ` + install.TableInstall + ` (
	version varchar(50) NOT NULL,
	module varchar(50) NOT NULL,
	stage varchar(200) NOT NULL,
	details varchar(5000) NOT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (version, module, stage)
);
`
		_, err = this_.DatabaseWorker.Exec(installSql, []interface{}{})
		if err != nil {
			return
		}
	}

	return
}

// Install 安装
func (this_ *InstallService) Install() (err error) {
	err = this_.Check()
	if err != nil {
		return
	}

	err = this_.InstallSteps(module_id.GetInstallStages())
	if err != nil {
		return
	}

	err = this_.InstallSteps(module_user.GetInstallStages())
	if err != nil {
		return
	}

	err = this_.InstallSteps(module_register.GetInstallStages())
	if err != nil {
		return
	}

	err = this_.InstallSteps(module_login.GetInstallStages())
	if err != nil {
		return
	}

	err = this_.InstallSteps(module_power.GetInstallStages())
	if err != nil {
		return
	}

	err = this_.InstallSteps(module_toolbox.GetInstallStages())
	if err != nil {
		return
	}

	return
}

// InstallSteps 根据阶段安装
func (this_ *InstallService) InstallSteps(stages []*install.StageModel) (err error) {
	for _, stage := range stages {
		err = this_.InstallStep(stage)
		if err != nil {
			return
		}
	}
	return
}

// InstallStep 根据阶段安装
func (this_ *InstallService) InstallStep(stage *install.StageModel) (err error) {
	var historyStageDetails *StageDetails
	historyStageDetails, err = this_.checkInstalled(stage.Version, stage.Module, stage.Stage)
	if err != nil {
		return
	}

	err = this_.execStage(historyStageDetails, stage)
	if err != nil {
		return
	}
	return
}

// checkExist 检测是否安装
func (this_ *InstallService) checkInstalled(version, module, stage string) (res *StageDetails, err error) {
	var list []map[string]interface{}

	sql := `SELECT details FROM ` + install.TableInstall + ` WHERE version=? AND module=? AND stage=? `
	list, err = this_.DatabaseWorker.QueryMap(sql, []interface{}{version, module, stage})
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = &StageDetails{}
		err = json.Unmarshal([]byte(list[0]["details"].(string)), &res)
		if err != nil {
			return
		}
	} else {
		res = nil
	}
	return
}

type StageDetails struct {
	Sql []string `json:"sql,omitempty"`
}

// execStage 执行阶段
func (this_ *InstallService) execStage(historyStageDetails *StageDetails, stage *install.StageModel) (err error) {

	var exeSQLs []string
	var exeErr error
	if stage.Sql != nil {
		exeSQLs, exeErr = this_.execStageSql(historyStageDetails, stage.Sql)
	}
	exeDetails := &StageDetails{
		Sql: exeSQLs,
	}
	var valueBytes []byte
	valueBytes, err = json.Marshal(exeDetails)
	if err != nil {
		return
	}

	if historyStageDetails != nil {
		sql := `UPDATE ` + install.TableInstall + ` SET details=?,updateTime=? WHERE version=? AND module=? AND stage=?`
		_, err = this_.DatabaseWorker.Exec(sql, []interface{}{string(valueBytes), time.Now(), stage.Version, stage.Module, stage.Stage})
	} else {
		sql := `INSERT INTO ` + install.TableInstall + ` (version, module, stage, details, createTime) VALUES(?, ?, ?, ?, ?)`
		_, err = this_.DatabaseWorker.Exec(sql, []interface{}{stage.Version, stage.Module, stage.Stage, string(valueBytes), time.Now()})
	}
	if err != nil {
		return
	}

	if exeErr != nil {
		return exeErr
	}
	return
}

// execStage 执行阶段
func (this_ *InstallService) execStageSql(historyStageDetails *StageDetails, stageSql *install.StageSqlModel) (exeSQLs []string, err error) {
	var sqs []string
	switch strings.ToLower(this_.DatabaseWorker.GetConfig().Type) {
	case "mysql":
		sqs = stageSql.Mysql
	case "sqlite":
		sqs = stageSql.Sqlite
	}
	var startIndex = -1
	if historyStageDetails != nil && len(historyStageDetails.Sql) > 0 {
		for index, sq := range historyStageDetails.Sql {
			exeSQLs = append(exeSQLs, sq)
			startIndex = index
			if index >= len(sqs) {
				this_.Logger.Error("安装脚本SQL数量小于已安装脚本SQL")
				err = errors.New(fmt.Sprint("安装脚本SQL数量小于已安装脚本SQL!"))
				return
			}
			if sq != sqs[index] {
				this_.Logger.Error("安装脚本SQL与已安装脚本SQL不一致")
				this_.Logger.Error("安装脚本SQL：" + sqs[index])
				this_.Logger.Error("已安装SQL：" + sq)
				err = errors.New(fmt.Sprint("安装脚本SQL与已安装脚本SQL不一致!"))
				return
			}
		}
	}
	if len(sqs) > 0 {
		for index, sq := range sqs {
			if index <= startIndex {
				continue
			}
			_, err = this_.DatabaseWorker.Exec(sq, []interface{}{})
			if err != nil {
				return
			}

			exeSQLs = append(exeSQLs, sq)
		}
	}
	return
}
