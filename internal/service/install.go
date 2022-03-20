package service

import (
	"encoding/json"
	"strings"
	"teamide/internal/model"
	"teamide/pkg/db"
	"time"
)

// NewInstallService 根据库配置创建InstallService
func NewInstallService(dbWorker db.DatabaseWorker) (res *InstallService) {
	res = &InstallService{
		dbWorker: dbWorker,
	}
	return
}

// InstallService 安装程序服务
type InstallService struct {
	dbWorker db.DatabaseWorker
}

// Check 检测基础配置是否完整，如基础表
func (this_ *InstallService) Check() (err error) {
	var isExist bool
	switch strings.ToLower(this_.dbWorker.GetConfig().Type) {
	case "mysql":
		sql := "SELECT count(1) FROM information_schema.TABLES WHERE table_schema='" + this_.dbWorker.GetConfig().Database + "' AND table_name ='" + model.TableInstall + "'"
		var count int64
		count, err = this_.dbWorker.Count(sql, []interface{}{})
		if err != nil {
			return
		}
		isExist = count > 0
	case "sqlite":
		sql := `SELECT count(1) FROM sqlite_master WHERE type="table" AND name="` + model.TableInstall + `" `
		var count int64
		count, err = this_.dbWorker.Count(sql, []interface{}{})
		if err != nil {
			return
		}
		isExist = count > 0
	}
	if !isExist {
		installSql := `
CREATE TABLE ` + model.TableInstall + ` (
	version varchar(50) NOT NULL,
	module varchar(50) NOT NULL,
	stage varchar(200) NOT NULL,
	details varchar(5000) NOT NULL,
	createTime datetime NOT NULL,
	PRIMARY KEY (version, module, stage)
);
`
		_, err = this_.dbWorker.Exec(installSql, []interface{}{})
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
	installStages := model.InstallStages
	for _, stage := range installStages {
		err = this_.InstallStep(stage)
		if err != nil {
			return
		}
	}

	return
}

// InstallStep 根据阶段安装
func (this_ *InstallService) InstallStep(stage *model.InstallStageModel) (err error) {
	var installed bool
	installed, err = this_.checkInstalled(stage.Version, stage.Module, stage.Stage)
	if err != nil {
		return
	}
	if installed {
		return
	}
	err = this_.execStage(stage)
	if err != nil {
		return
	}
	return
}

// checkExist 检测是否安装
func (this_ *InstallService) checkInstalled(version, module, stage string) (installed bool, err error) {

	sql := `SELECT COUNT(1) FROM ` + model.TableInstall + ` WHERE version=? AND module=? AND stage=? `
	var count int64
	count, err = this_.dbWorker.Count(sql, []interface{}{version, module, stage})
	if err != nil {
		return
	}

	installed = count > 0
	return
}

// execStage 执行阶段
func (this_ *InstallService) execStage(stage *model.InstallStageModel) (err error) {
	value := map[string]interface{}{}
	if stage.Sql != nil {
		value["sql"] = stage.Sql
		err = this_.execStageSql(stage.Sql)
		if err != nil {
			return
		}
	}

	var valueBytes []byte
	valueBytes, err = json.Marshal(value)
	if err != nil {
		return
	}

	sql := `INSERT INTO ` + model.TableInstall + ` (version, module, stage, details, createTime) VALUES(?, ?, ?, ?, ?)`
	_, err = this_.dbWorker.Exec(sql, []interface{}{stage.Version, stage.Module, stage.Stage, string(valueBytes), time.Now()})
	if err != nil {
		return
	}

	return
}

// execStage 执行阶段
func (this_ *InstallService) execStageSql(stageSql *model.InstallStageSqlModel) (err error) {
	var execSqls []string
	switch strings.ToLower(this_.dbWorker.GetConfig().Type) {
	case "mysql":
		execSqls = stageSql.Mysql
	case "sqlite":
		execSqls = stageSql.Sqlite
	}
	if len(execSqls) > 0 {
		for _, execSql := range execSqls {
			_, err = this_.dbWorker.Exec(execSql, []interface{}{})
			if err != nil {
				return
			}
		}
	}
	return
}
