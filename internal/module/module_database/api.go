package module_database

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"github.com/team-ide/go-tool/db"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
	"teamide/pkg/ssh"
)

type api struct {
	toolboxService *module_toolbox.ToolboxService
}

func NewApi(toolboxService *module_toolbox.ToolboxService) *api {
	return &api{
		toolboxService: toolboxService,
	}
}

var (
	Power               = base.AppendPower(&base.PowerAction{Action: "database", Text: "数据库", ShouldLogin: true, StandAlone: true})
	infoPower           = base.AppendPower(&base.PowerAction{Action: "info", Text: "数据库信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	dataPower           = base.AppendPower(&base.PowerAction{Action: "data", Text: "数据库基础数据", ShouldLogin: true, StandAlone: true, Parent: Power})
	ownersPower         = base.AppendPower(&base.PowerAction{Action: "owners", Text: "数据库查询", ShouldLogin: true, StandAlone: true, Parent: Power})
	ownerCreatePower    = base.AppendPower(&base.PowerAction{Action: "ownerCreate", Text: "数据库库创建", ShouldLogin: true, StandAlone: true, Parent: Power})
	ownerDeletePower    = base.AppendPower(&base.PowerAction{Action: "ownerDelete", Text: "数据库库删除", ShouldLogin: true, StandAlone: true, Parent: Power})
	ownerCreateSqlPower = base.AppendPower(&base.PowerAction{Action: "ownerCreateSql", Text: "数据库库删除SQL", ShouldLogin: true, StandAlone: true, Parent: Power})
	ddlPower            = base.AppendPower(&base.PowerAction{Action: "ddl", Text: "数据库DDL", ShouldLogin: true, StandAlone: true, Parent: Power})
	modelPower          = base.AppendPower(&base.PowerAction{Action: "model", Text: "数据库模型", ShouldLogin: true, StandAlone: true, Parent: Power})
	tablesPower         = base.AppendPower(&base.PowerAction{Action: "tables", Text: "数据库库表查询", ShouldLogin: true, StandAlone: true, Parent: Power})
	tableDetailPower    = base.AppendPower(&base.PowerAction{Action: "tableDetail", Text: "数据库库表详细信息查询", ShouldLogin: true, StandAlone: true, Parent: Power})
	tableCreatePower    = base.AppendPower(&base.PowerAction{Action: "tableCreate", Text: "数据库创建表", ShouldLogin: true, StandAlone: true, Parent: Power})
	tableCreateSqlPower = base.AppendPower(&base.PowerAction{Action: "tableCreateSql", Text: "数据库创建表SQL", ShouldLogin: true, StandAlone: true, Parent: Power})
	tableUpdatePower    = base.AppendPower(&base.PowerAction{Action: "tableUpdate", Text: "数据库修改表", ShouldLogin: true, StandAlone: true, Parent: Power})
	tableUpdateSqlPower = base.AppendPower(&base.PowerAction{Action: "tableUpdateSql", Text: "数据库修改表SQL", ShouldLogin: true, StandAlone: true, Parent: Power})
	tableDeletePower    = base.AppendPower(&base.PowerAction{Action: "tableDelete", Text: "数据库删除表", ShouldLogin: true, StandAlone: true, Parent: Power})
	tableDataTrimPower  = base.AppendPower(&base.PowerAction{Action: "tableDataTrim", Text: "数据库表数据清空", ShouldLogin: true, StandAlone: true, Parent: Power})
	tableDataPower      = base.AppendPower(&base.PowerAction{Action: "tableData", Text: "数据库表数据查询", ShouldLogin: true, StandAlone: true, Parent: Power})
	dataListSqlPower    = base.AppendPower(&base.PowerAction{Action: "dataListSql", Text: "数据库数据转换SQL", ShouldLogin: true, StandAlone: true, Parent: Power})
	dataListExecPower   = base.AppendPower(&base.PowerAction{Action: "dataListExec", Text: "数据库数据执行", ShouldLogin: true, StandAlone: true, Parent: Power})
	executeSQLPower     = base.AppendPower(&base.PowerAction{Action: "executeSQL", Text: "数据库SQL执行", ShouldLogin: true, StandAlone: true, Parent: Power})
	importPower         = base.AppendPower(&base.PowerAction{Action: "import", Text: "数据库导入", ShouldLogin: true, StandAlone: true, Parent: Power})
	exportPower         = base.AppendPower(&base.PowerAction{Action: "export", Text: "数据库导出", ShouldLogin: true, StandAlone: true, Parent: Power})
	exportDownloadPower = base.AppendPower(&base.PowerAction{Action: "exportDownload", Text: "数据库导出下载", ShouldLogin: true, StandAlone: true, Parent: Power})
	syncPower           = base.AppendPower(&base.PowerAction{Action: "sync", Text: "数据库同步", ShouldLogin: true, StandAlone: true, Parent: Power})
	taskStatusPower     = base.AppendPower(&base.PowerAction{Action: "taskStatus", Text: "数据库任务状态查询", ShouldLogin: true, StandAlone: true, Parent: Power})
	taskStopPower       = base.AppendPower(&base.PowerAction{Action: "taskStop", Text: "数据库任务停止", ShouldLogin: true, StandAlone: true, Parent: Power})
	taskCleanPower      = base.AppendPower(&base.PowerAction{Action: "taskClean", Text: "数据库任务清理", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower          = base.AppendPower(&base.PowerAction{Action: "close", Text: "数据库关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: infoPower, Do: this_.info})
	apis = append(apis, &base.ApiWorker{Power: dataPower, Do: this_.data, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: ownersPower, Do: this_.owners})
	apis = append(apis, &base.ApiWorker{Power: ownerCreatePower, Do: this_.ownerCreate})
	apis = append(apis, &base.ApiWorker{Power: ownerDeletePower, Do: this_.ownerDelete})
	apis = append(apis, &base.ApiWorker{Power: ownerCreateSqlPower, Do: this_.ownerCreateSql})
	apis = append(apis, &base.ApiWorker{Power: ddlPower, Do: this_.ddl})
	apis = append(apis, &base.ApiWorker{Power: modelPower, Do: this_.model})
	apis = append(apis, &base.ApiWorker{Power: tablesPower, Do: this_.tables})
	apis = append(apis, &base.ApiWorker{Power: tableDetailPower, Do: this_.tableDetail})
	apis = append(apis, &base.ApiWorker{Power: tableCreatePower, Do: this_.tableCreate})
	apis = append(apis, &base.ApiWorker{Power: tableCreateSqlPower, Do: this_.tableCreateSql})
	apis = append(apis, &base.ApiWorker{Power: tableUpdatePower, Do: this_.tableUpdate})
	apis = append(apis, &base.ApiWorker{Power: tableUpdateSqlPower, Do: this_.tableUpdateSql})
	apis = append(apis, &base.ApiWorker{Power: tableDeletePower, Do: this_.tableDelete})
	apis = append(apis, &base.ApiWorker{Power: tableDataTrimPower, Do: this_.tableDataTrim})
	apis = append(apis, &base.ApiWorker{Power: tableDataPower, Do: this_.tableData})
	apis = append(apis, &base.ApiWorker{Power: dataListSqlPower, Do: this_.dataListSql})
	apis = append(apis, &base.ApiWorker{Power: dataListExecPower, Do: this_.dataListExec})
	apis = append(apis, &base.ApiWorker{Power: executeSQLPower, Do: this_.executeSQL})
	apis = append(apis, &base.ApiWorker{Power: importPower, Do: this_._import})
	apis = append(apis, &base.ApiWorker{Power: exportPower, Do: this_.export})
	apis = append(apis, &base.ApiWorker{Power: exportDownloadPower, Do: this_.exportDownload})
	apis = append(apis, &base.ApiWorker{Power: syncPower, Do: this_.sync})
	apis = append(apis, &base.ApiWorker{Power: taskStatusPower, Do: this_.taskStatus, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: taskStopPower, Do: this_.taskStop})
	apis = append(apis, &base.ApiWorker{Power: taskCleanPower, Do: this_.taskClean})
	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *db.Config, sshConfig *ssh.Config, err error) {
	config = &db.Config{}
	sshConfig, err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func getService(config *db.Config, sshConfig *ssh.Config) (res db.IService, err error) {
	key := fmt.Sprint("database-", config.Type, "-", config.Host, "-", config.Port)
	if config.DatabasePath != "" {
		key += "-" + config.DatabasePath
	}
	if config.Database != "" {
		key += "-" + config.Database
	}
	if config.DbName != "" {
		key += "-" + config.DbName
	}
	if config.OdbcDsn != "" {
		key += "-" + config.OdbcDsn
	}
	if config.OdbcDialectName != "" {
		key += "-" + config.OdbcDialectName
	}
	if config.Username != "" {
		key += "-" + base.GetMd5String(key+config.Username)
	}
	if config.Password != "" {
		key += "-" + base.GetMd5String(key+config.Password)
	}
	if sshConfig != nil {
		key += "-ssh-" + sshConfig.Address
		key += "-ssh-" + sshConfig.Username
	}

	var serviceInfo *base.ServiceInfo
	serviceInfo, err = base.GetService(key, func() (res *base.ServiceInfo, err error) {
		var s db.IService
		if sshConfig != nil {
			config.SSHClient, err = ssh.NewClient(*sshConfig)
			if err != nil {
				util.Logger.Error("getDatabaseService ssh NewClient error", zap.Any("key", key), zap.Error(err))
				return
			}
		}
		s, err = db.New(config)
		if err != nil {
			util.Logger.Error("getDatabaseService error", zap.Any("key", key), zap.Error(err))
			return
		}
		res = &base.ServiceInfo{
			WaitTime:    10 * 60 * 1000,
			LastUseTime: util.GetNowTime(),
			Service:     s,
			Stop:        s.Stop,
		}
		return
	})
	if err != nil {
		return
	}
	res = serviceInfo.Service.(db.IService)
	serviceInfo.SetLastUseTime()
	return
}

type BaseRequest struct {
	WorkerId     string                 `json:"workerId"`
	OwnerName    string                 `json:"ownerName"`
	TableName    string                 `json:"tableName"`
	TaskId       string                 `json:"taskId"`
	ExecuteSQL   string                 `json:"executeSQL"`
	ColumnList   []*dialect.ColumnModel `json:"columnList"`
	Wheres       []*dialect.Where       `json:"wheres"`
	Orders       []*dialect.Order       `json:"orders"`
	PageNo       int                    `json:"pageNo"`
	PageSize     int                    `json:"pageSize"`
	DatabaseType string                 `json:"databaseType"`

	InsertList      []map[string]interface{} `json:"insertList"`
	UpdateList      []map[string]interface{} `json:"updateList"`
	UpdateWhereList []map[string]interface{} `json:"updateWhereList"`
	DeleteList      []map[string]interface{} `json:"deleteList"`
}

func (this_ *api) info(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	res, err = service.Info()
	if err != nil {
		return
	}
	return
}

func (this_ *api) data(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	param := this_.getParam(requestBean, c)

	data := make(map[string]interface{})
	data["columnTypeInfoList"] = service.GetTargetDialect(param).GetColumnTypeInfos()
	data["indexTypeInfoList"] = service.GetTargetDialect(param).GetIndexTypeInfos()
	res = data
	return
}

func (this_ *api) owners(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	param := this_.getParam(requestBean, c)

	res, err = service.OwnersSelect(param)
	if err != nil {
		return
	}
	return
}

func (this_ *api) ownerCreate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	param := this_.getParam(requestBean, c)
	var owner = &dialect.OwnerModel{}
	if !base.RequestJSON(owner, c) {
		return
	}
	res, err = service.OwnerCreate(param, owner)
	if err != nil {
		return
	}
	return
}

func (this_ *api) ownerDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	param := this_.getParam(requestBean, c)
	res, err = service.OwnerDelete(param, request.OwnerName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) ownerCreateSql(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	param := this_.getParam(requestBean, c)

	var owner = &dialect.OwnerModel{}
	if !base.RequestJSON(owner, c) {
		return
	}
	res, err = service.OwnerCreateSql(param, owner)
	if err != nil {
		return
	}
	return
}

func (this_ *api) ddl(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	res, err = service.DDL(param, request.OwnerName, request.TableName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) model(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	res, err = service.Model(param, request.OwnerName, request.TableName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) tables(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	res, err = service.TablesSelect(param, request.OwnerName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) tableDetail(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	res, err = service.TableDetail(param, request.OwnerName, request.TableName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) tableCreate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	var table = &dialect.TableModel{}
	if !base.RequestJSON(table, c) {
		return
	}

	err = service.TableCreate(param, request.OwnerName, table)
	if err != nil {
		return
	}
	return
}

func (this_ *api) tableCreateSql(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	var table = &dialect.TableModel{}
	if !base.RequestJSON(table, c) {
		return
	}

	res, err = service.TableCreateSql(param, request.OwnerName, table)
	if err != nil {
		return
	}
	return
}

func (this_ *api) tableUpdate(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	var updateTableParam = &db.UpdateTableParam{}
	if !base.RequestJSON(updateTableParam, c) {
		return
	}

	err = service.TableUpdate(param, request.OwnerName, request.TableName, updateTableParam)
	if err != nil {
		return
	}
	return
}

func (this_ *api) tableUpdateSql(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	var updateTableParam = &db.UpdateTableParam{}
	if !base.RequestJSON(updateTableParam, c) {
		return
	}

	res, err = service.TableUpdateSql(param, request.OwnerName, request.TableName, updateTableParam)
	if err != nil {
		return
	}
	return
}

func (this_ *api) tableDelete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	err = service.TableDelete(param, request.OwnerName, request.TableName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) tableDataTrim(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	err = service.TableDataTrim(param, request.OwnerName, request.TableName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) tableData(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	res, err = service.TableData(param, request.OwnerName, request.TableName, request.ColumnList, request.Wheres, request.Orders, request.PageSize, request.PageNo)
	if err != nil {
		return
	}
	return
}

func (this_ *api) dataListSql(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	res, err = service.DataListSql(param, request.OwnerName, request.TableName, request.ColumnList,
		request.InsertList,
		request.UpdateList, request.UpdateWhereList,
		request.DeleteList,
	)
	if err != nil {
		return
	}
	return
}

func (this_ *api) getParam(requestBean *base.RequestBean, c *gin.Context) (param *db.Param) {

	param = &db.Param{}
	if !base.RequestJSON(param, c) {
		return
	}
	if param.ParamModel == nil {
		param.ParamModel = &dialect.ParamModel{}
	}

	return
}

func (this_ *api) dataListExec(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	err = service.DataListExec(param, request.OwnerName, request.TableName, request.ColumnList,
		request.InsertList,
		request.UpdateList, request.UpdateWhereList,
		request.DeleteList,
	)
	if err != nil {
		return
	}
	return
}

func (this_ *api) executeSQL(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	data := make(map[string]interface{})
	data["executeList"], data["error"], err = service.ExecuteSQL(param, request.OwnerName, request.ExecuteSQL)
	if err != nil {
		return
	}
	res = data
	return
}

func (this_ *api) _import(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	var importParam = &worker.TaskImportParam{}
	if !base.RequestJSON(importParam, c) {
		return
	}

	var task *worker.Task
	task, err = service.StartImport(param, importParam)
	if err != nil {
		return
	}
	res = task

	if task != nil {
		addWorkerTask(request.WorkerId, task.TaskId)
	}
	return
}

func (this_ *api) export(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	var exportParam = &worker.TaskExportParam{}
	if !base.RequestJSON(exportParam, c) {
		return
	}

	var task *worker.Task
	task, err = service.StartExport(param, exportParam)
	if err != nil {
		return
	}
	res = task

	if task != nil {
		addWorkerTask(request.WorkerId, task.TaskId)
	}
	return
}

func (this_ *api) exportDownload(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	data := map[string]string{}
	err = c.Bind(&data)
	if err != nil {
		return
	}

	taskId := data["taskId"]
	if taskId == "" {
		err = errors.New("taskId获取失败")
		return
	}

	task := worker.GetTask(taskId)
	if task == nil {
		err = errors.New("任务不存在")
		return
	}
	if task.Extend == nil || task.Extend["downloadPath"] == "" {
		err = errors.New("任务导出文件丢失")
		return
	}
	tempDir, err := util.GetTempDir()
	if err != nil {
		return
	}

	path := tempDir + task.Extend["downloadPath"].(string)
	exists, err := util.PathExists(path)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New("文件不存在")
		return
	}
	var fileName string
	var fileSize int64
	ff, err := os.Lstat(path)
	if err != nil {
		return
	}
	var fileInfo *os.File
	if ff.IsDir() {
		exists, err = util.PathExists(path + ".zip")
		if err != nil {
			return
		}
		if !exists {
			err = util.Zip(path, path+".zip")
			if err != nil {
				return
			}
		}
		ff, err = os.Lstat(path + ".zip")
		if err != nil {
			return
		}
		fileInfo, err = os.Open(path + ".zip")
		if err != nil {
			return
		}
	} else {
		fileInfo, err = os.Open(path)
		if err != nil {
			return
		}
	}
	fileName = ff.Name()
	fileSize = ff.Size()

	defer func() {
		_ = fileInfo.Close()
	}()

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", fmt.Sprint(fileSize))
	c.Header("download-file-name", fileName)

	_, err = io.Copy(c.Writer, fileInfo)
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	res = base.HttpNotResponse
	return
}

func (this_ *api) sync(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, sshConfig, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config, sshConfig)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	param := this_.getParam(requestBean, c)

	var syncParam = &worker.TaskSyncParam{}
	if !base.RequestJSON(syncParam, c) {
		return
	}

	var task *worker.Task
	task, err = service.StartSync(param, syncParam)
	if err != nil {
		return
	}
	res = task

	if task != nil {
		addWorkerTask(request.WorkerId, task.TaskId)
	}
	return
}

func (this_ *api) taskStatus(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res = worker.GetTask(request.TaskId)
	return
}

func (this_ *api) taskStop(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	worker.StopTask(request.TaskId)
	return
}

func (this_ *api) taskClean(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	task := worker.GetTask(request.TaskId)
	if task != nil {
		if task.Extend != nil {
			if task.Extend["dirPath"] != "" {
				_ = os.RemoveAll(task.Extend["dirPath"].(string))
			}
			if task.Extend["zipPath"] != "" {
				_ = os.Remove(task.Extend["zipPath"].(string))
			}
		}
	}
	worker.ClearTask(request.TaskId)
	return
}

func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	removeWorkerTasks(request.WorkerId)
	return
}

var workerTasksCache = map[string][]string{}
var workerTasksCacheLock = &sync.Mutex{}

func addWorkerTask(workerId string, taskId string) {
	workerTasksCacheLock.Lock()
	defer workerTasksCacheLock.Unlock()
	taskIds := workerTasksCache[workerId]
	if util.StringIndexOf(taskIds, taskId) < 0 {
		taskIds = append(taskIds, taskId)
		workerTasksCache[workerId] = taskIds
	}
	return
}
func removeWorkerTasks(workerId string) {
	workerTasksCacheLock.Lock()
	defer workerTasksCacheLock.Unlock()
	taskIds := workerTasksCache[workerId]
	for _, taskId := range taskIds {
		task := worker.GetTask(taskId)
		if task != nil {
			if task.Extend != nil {
				if task.Extend["dirPath"] != "" {
					_ = os.RemoveAll(task.Extend["dirPath"].(string))
				}
				if task.Extend["zipPath"] != "" {
					_ = os.Remove(task.Extend["zipPath"].(string))
				}
			}
			worker.ClearTask(taskId)
		}
	}
	delete(workerTasksCache, workerId)
	return
}
