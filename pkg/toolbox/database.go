package toolbox

import (
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"go.uber.org/zap"
	"os"
	"teamide/pkg/db"
	"teamide/pkg/util"
)

type DatabaseBaseRequest struct {
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

func DatabaseWork(work string, config *db.DatabaseConfig, data map[string]interface{}) (res map[string]interface{}, err error) {
	var service *db.Service

	service, err = getDatabaseService(config)
	if err != nil {
		return
	}

	dataBS, err := json.Marshal(data)
	if err != nil {
		return
	}
	request := &DatabaseBaseRequest{}
	err = json.Unmarshal(dataBS, request)
	if err != nil {
		return
	}

	var param = &db.Param{}
	err = json.Unmarshal(dataBS, param)
	if err != nil {
		return
	}
	if param.ParamModel == nil {
		param.ParamModel = &dialect.ParamModel{}
	}

	res = map[string]interface{}{}
	switch work {
	case "data":
		res["columnTypeInfoList"] = service.GetTargetDialect(param).GetColumnTypeInfos()
		res["indexTypeInfoList"] = service.GetTargetDialect(param).GetIndexTypeInfos()

		break
	case "owners":
		var owners []*dialect.OwnerModel
		owners, err = service.OwnersSelect(param)
		if err != nil {
			return
		}
		res["owners"] = owners

		break
	case "ownerCreate":
		var owner = &dialect.OwnerModel{}
		err = json.Unmarshal(dataBS, owner)
		if err != nil {
			return
		}
		var created bool
		created, err = service.OwnerCreate(param, owner)
		if err != nil {
			return
		}
		res["created"] = created

		break
	case "ownerDelete":
		var deleted bool
		deleted, err = service.OwnerDelete(param, request.OwnerName)
		if err != nil {
			return
		}
		res["deleted"] = deleted

		break
	case "ownerCreateSql":
		var owner = &dialect.OwnerModel{}
		err = json.Unmarshal(dataBS, owner)
		if err != nil {
			return
		}
		var sqlList []string
		sqlList, err = service.DatabaseWorker.OwnerCreateSql(param.ParamModel, owner)
		if err != nil {
			return
		}
		res["sqlList"] = sqlList

		break
	case "ddl":
		var sqlList []string
		sqlList, err = service.DDL(param, request.OwnerName, request.TableName)
		if err != nil {
			return
		}
		res["sqlList"] = sqlList
		break
	case "model":
		var content string
		content, err = service.Model(param, request.OwnerName, request.TableName)
		if err != nil {
			return
		}
		res["content"] = content
		break
	case "tables":
		var tables []*dialect.TableModel
		tables, err = service.TablesSelect(param, request.OwnerName)
		if err != nil {
			return
		}
		res["tables"] = tables

		break
	case "tableDetail":
		var table *dialect.TableModel
		table, err = service.TableDetail(param, request.OwnerName, request.TableName)
		if err != nil {
			return
		}
		res["table"] = table

		break
	case "tableCreate":
		var table = &dialect.TableModel{}
		err = json.Unmarshal(dataBS, table)
		if err != nil {
			return
		}
		err = service.TableCreate(param, request.OwnerName, table)
		if err != nil {
			return
		}

		break
	case "tableCreateSql":
		var table = &dialect.TableModel{}
		err = json.Unmarshal(dataBS, table)
		if err != nil {
			return
		}
		var sqlList []string
		sqlList, err = service.TableCreateSql(param, request.OwnerName, table)
		if err != nil {
			return
		}
		res["sqlList"] = sqlList

		break
	case "tableUpdate":
		var updateTableParam = &db.UpdateTableParam{}
		err = json.Unmarshal(dataBS, updateTableParam)
		if err != nil {
			return
		}
		err = service.TableUpdate(param, request.OwnerName, request.TableName, updateTableParam)
		if err != nil {
			return
		}

		break
	case "tableUpdateSql":
		var updateTableParam = &db.UpdateTableParam{}
		err = json.Unmarshal(dataBS, updateTableParam)
		if err != nil {
			return
		}
		var sqlList []string
		sqlList, err = service.TableUpdateSql(param, request.OwnerName, request.TableName, updateTableParam)
		if err != nil {
			return
		}
		res["sqlList"] = sqlList

		break
	case "tableDelete":
		err = service.TableDelete(param, request.OwnerName, request.TableName)
		if err != nil {
			return
		}

		break
	case "tableDataTrim":
		err = service.TableDataTrim(param, request.OwnerName, request.TableName)
		if err != nil {
			return
		}

		break
	case "tableData":
		var dataListRequest db.DataListResult
		dataListRequest, err = service.TableData(param, request.OwnerName, request.TableName, request.ColumnList, request.Wheres, request.Orders, request.PageSize, request.PageNo)
		if err != nil {
			return
		}
		res["sql"] = dataListRequest.Sql
		res["total"] = dataListRequest.Total
		res["dataList"] = dataListRequest.DataList

		break
	case "dataListSql":
		var sqlList []string
		sqlList, err = service.DataListSql(param, request.OwnerName, request.TableName, request.ColumnList,
			request.InsertList,
			request.UpdateList, request.UpdateWhereList,
			request.DeleteList,
		)
		if err != nil {
			return
		}
		res["sqlList"] = sqlList
		break
	case "dataListExec":
		err = service.DataListExec(param, request.OwnerName, request.TableName, request.ColumnList,
			request.InsertList,
			request.UpdateList, request.UpdateWhereList,
			request.DeleteList,
		)
		if err != nil {
			return
		}
		break
	case "executeSQL":
		res["executeList"], res["error"], err = service.ExecuteSQL(param, request.OwnerName, request.ExecuteSQL)
		if err != nil {
			return
		}
		break
	case "import":
		var importParam = &worker.TaskImportParam{}
		err = json.Unmarshal(dataBS, importParam)
		if err != nil {
			return
		}

		var task *worker.Task
		task, err = service.StartImport(param, importParam)
		if err != nil {
			return
		}
		res["task"] = task
		break
	case "export":
		var exportParam = &worker.TaskExportParam{}
		err = json.Unmarshal(dataBS, exportParam)
		if err != nil {
			return
		}

		var task *worker.Task
		task, err = service.StartExport(param, exportParam)
		if err != nil {
			return
		}
		res["task"] = task

		break
	case "sync":
		var syncParam = &worker.TaskSyncParam{}
		err = json.Unmarshal(dataBS, syncParam)
		if err != nil {
			return
		}

		var task *worker.Task
		task, err = service.StartSync(param, syncParam)
		if err != nil {
			return
		}
		res["task"] = task
		break
	case "taskStatus":
		task := worker.GetTask(request.TaskId)
		res["task"] = task
		break
	case "taskStop":
		worker.StopTask(request.TaskId)
		break
	case "taskClean":
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
		break
	}
	return
}

type SqlConditionalOperation struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	SqlConditionalOperations []*SqlConditionalOperation
)

func init() {
	SqlConditionalOperations = []*SqlConditionalOperation{
		{Text: "等于", Value: "="},
		{Text: "不等于", Value: "<>"},
		{Text: "大于", Value: ">"},
		{Text: "大于或等于", Value: ">="},
		{Text: "小于", Value: "<"},
		{Text: "小于或等于", Value: "<="},
		{Text: "包含", Value: "like"},
		{Text: "不包含", Value: "not like"},
		{Text: "开始以", Value: "like start"},
		{Text: "开始不是以", Value: "not like start"},
		{Text: "结束以", Value: "like end"},
		{Text: "结束不是以", Value: "not like end"},
		{Text: "是null", Value: "is null"},
		{Text: "不是null", Value: "is not null"},
		{Text: "是空", Value: "is empty"},
		{Text: "不是空", Value: "is not empty"},
		{Text: "介于", Value: "between"},
		{Text: "不介于", Value: "not between"},
		{Text: "在列表", Value: "in"},
		{Text: "不在列表", Value: "not in"},
		{Text: "自定义", Value: "custom"},
	}
}

func getDatabaseService(config *db.DatabaseConfig) (res *db.Service, err error) {
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
	if config.Username != "" {
		key += "-" + util.GetMd5String(key+config.Username)
	}
	if config.Password != "" {
		key += "-" + util.GetMd5String(key+config.Password)
	}
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s *db.Service
		s, err = db.CreateService(config)
		if err != nil {
			util.Logger.Error("getDatabaseService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
			return
		}
		res = s
		return
	})
	if err != nil {
		return
	}
	res = service.(*db.Service)
	res.SetLastUseTime()
	return
}
