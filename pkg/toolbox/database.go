package toolbox

import (
	"encoding/json"
	"fmt"
	"teamide/pkg/db"
	"teamide/pkg/form"
	"teamide/pkg/util"
)

func init() {

	worker_ := &Worker{
		Name: "database",
		Text: "Database",
		Work: databaseWork,
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "类型", Name: "type", Type: "select", DefaultValue: "mysql",
					Options: []*form.Option{
						{Text: "MySql", Value: "mysql"},
					},
					Rules: []*form.Rule{
						{Required: true, Message: "数据库类型不能为空"},
					},
				},
				{
					Label: "Host（127.0.0.1）", Name: "host", DefaultValue: "127.0.0.1",
					Rules: []*form.Rule{
						{Required: true, Message: "数据库连接地址不能为空"},
					},
				},
				{
					Label: "Port（3306）", Name: "port", IsNumber: true, DefaultValue: "3306",
					Rules: []*form.Rule{
						{Required: true, Message: "数据库连接端口不能为空"},
					},
				},
				{Label: "Username", Name: "username"},
				{Label: "Password", Name: "password", Type: "password"},
			},
		},
	}

	AddWorker(worker_)
}

type DatabaseBaseRequest struct {
	Database        string                   `json:"database"`
	Table           string                   `json:"table"`
	TaskKey         string                   `json:"taskKey"`
	ExecuteSQL      string                   `json:"executeSQL"`
	ColumnList      []*db.TableColumnModel   `json:"columnList"`
	Wheres          []*db.Where              `json:"wheres"`
	Orders          []*db.Order              `json:"orders"`
	PageIndex       int                      `json:"pageIndex"`
	PageSize        int                      `json:"pageSize"`
	DatabaseType    string                   `json:"databaseType"`
	ImportDataList  []map[string]interface{} `json:"importDataList"`
	InsertList      []map[string]interface{} `json:"insertList"`
	UpdateList      []map[string]interface{} `json:"updateList"`
	UpdateWhereList []map[string]interface{} `json:"updateWhereList"`
	DeleteList      []map[string]interface{} `json:"deleteList"`
}

func databaseWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {
	var service DatabaseService

	var databaseConfig db.DatabaseConfig
	var bs []byte
	bs, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &databaseConfig)
	if err != nil {
		return
	}

	service, err = getDatabaseService(databaseConfig)
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

	res = map[string]interface{}{}
	switch work {
	case "databases":
		var databases []*db.DatabaseModel
		databases, err = service.Databases()
		if err != nil {
			return
		}
		res["databases"] = databases
	case "tables":
		var tables []*db.TableModel
		tables, err = service.Tables(request.Database)
		if err != nil {
			return
		}
		res["tables"] = tables
	case "tableDetail":
		var tables []*db.TableModel
		tables, err = service.TableDetails(request.Database, request.Table)
		if err != nil {
			return
		}
		if len(tables) > 0 {
			res["table"] = tables[0]
		}
	case "ddl":
		var generateParam = &db.GenerateParam{}
		err = json.Unmarshal(dataBS, generateParam)
		if err != nil {
			return
		}

		if generateParam.DatabaseType == "" {
			generateParam.DatabaseType = databaseConfig.Type
		}

		var sqlList []string
		if generateParam.GenerateDatabase {
			var sqlList_ []string
			sqlList_, err = db.ToDatabaseDDL(generateParam, &db.DatabaseModel{Name: request.Database})
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}

		var tables []*db.TableModel
		tables, err = service.TableDetails(request.Database, request.Table)
		if err != nil {
			return
		}
		for _, table := range tables {
			var sqlList_ []string
			sqlList_, err = db.ToTableDDL(generateParam, request.Database, table)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}

		res["sqlList"] = sqlList
	case "dataList":

		var generateParam = &db.GenerateParam{}
		err = json.Unmarshal(dataBS, generateParam)
		if err != nil {
			return
		}

		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		var dataListRequest DataListResult
		dataListRequest, err = service.DataList(generateParam, DataListParam{
			Database:   request.Database,
			Table:      request.Table,
			ColumnList: request.ColumnList,
			Wheres:     request.Wheres,
			Orders:     request.Orders,
			PageIndex:  request.PageIndex,
			PageSize:   request.PageSize,
		})
		if err != nil {
			return
		}
		res["sql"] = dataListRequest.Sql
		res["params"] = dataListRequest.Params
		res["total"] = dataListRequest.Total
		res["dataList"] = dataListRequest.DataList
	case "executeSQL":
		var generateParam = &db.GenerateParam{}
		err = json.Unmarshal(dataBS, generateParam)
		if err != nil {
			return
		}

		executeSQLTask := &executeSQLTask{
			Database:      request.Database,
			ExecuteSQL:    request.ExecuteSQL,
			service:       service,
			generateParam: generateParam,
		}
		executeSQLTask.Start()
		res["task"] = executeSQLTask

	case "dataListSql":
		var generateParam = &db.GenerateParam{}
		err = json.Unmarshal(dataBS, generateParam)
		if err != nil {
			return
		}
		saveDataListTask := &saveDataListTask{
			Database:        request.Database,
			Table:           request.Table,
			ColumnList:      request.ColumnList,
			InsertList:      request.InsertList,
			UpdateList:      request.UpdateList,
			UpdateWhereList: request.UpdateWhereList,
			DeleteList:      request.DeleteList,
			service:         service,
			generateParam:   generateParam,
		}

		var sqlList []string
		var valuesList [][]interface{}

		sqlList, valuesList, err = saveDataListTask.SaveDataListSql()
		if err != nil {
			return
		}
		res["sqlList"] = sqlList
		res["valuesList"] = valuesList
	case "saveDataList":
		var generateParam = &db.GenerateParam{}
		err = json.Unmarshal(dataBS, generateParam)
		if err != nil {
			return
		}

		generateParam.OpenTransaction = true
		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		saveDataListTask := &saveDataListTask{
			Database:        request.Database,
			Table:           request.Table,
			ColumnList:      request.ColumnList,
			InsertList:      request.InsertList,
			UpdateList:      request.UpdateList,
			UpdateWhereList: request.UpdateWhereList,
			DeleteList:      request.DeleteList,
			service:         service,
			generateParam:   generateParam,
		}
		err = saveDataListTask.Start()
		if err != nil {
			return
		}
		res["task"] = saveDataListTask

	case "importDataForStrategy":

		var generateParam = &db.GenerateParam{}
		err = json.Unmarshal(dataBS, generateParam)
		if err != nil {
			return
		}

		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		taskKey := util.GenerateUUID()
		importDataForStrategyTask := &importDataForStrategyTask{
			Key:            taskKey,
			Database:       request.Database,
			Table:          request.Table,
			ColumnList:     request.ColumnList,
			ImportDataList: request.ImportDataList,
			service:        service,
			generateParam:  generateParam,
		}
		addImportDataForStrategyTask(importDataForStrategyTask)

		res["taskKey"] = taskKey
	case "importDataForStrategyStatus":
		importDataForStrategyTask := importDataForStrategyTaskCache[request.TaskKey]
		res["task"] = importDataForStrategyTask
	case "importDataForStrategyStop":
		importDataForStrategyTask := importDataForStrategyTaskCache[request.TaskKey]
		if importDataForStrategyTask != nil {
			importDataForStrategyTask.Stop()
		}
	case "importDataForStrategyClean":
		delete(importDataForStrategyTaskCache, request.TaskKey)
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

func getDatabaseService(config db.DatabaseConfig) (res DatabaseService, err error) {
	key := fmt.Sprint("database-", config.Type, "-", config.Host, "-", config.Port)
	if config.Username != "" {
		key += "-" + util.GetMd5String(key+config.Username)
	}
	if config.Password != "" {
		key += "-" + util.GetMd5String(key+config.Password)
	}
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s DatabaseService
		s, err = CreateDatabaseService(config)
		if err != nil {
			return
		}
		res = s
		return
	})
	if err != nil {
		return
	}
	res = service.(DatabaseService)
	res.SetLastUseTime()
	return
}

func CreateDatabaseService(config db.DatabaseConfig) (service DatabaseService, err error) {
	service, err = CreateMysqlService(config)
	return
}

type DatabaseService interface {
	GetWaitTime() int64
	GetLastUseTime() int64
	SetLastUseTime()
	Stop()
	GetDatabaseWorker() *db.DatabaseWorker
	Databases() ([]*db.DatabaseModel, error)
	Tables(database string) ([]*db.TableModel, error)
	TableDetails(database string, table string) ([]*db.TableModel, error)
	DataList(param *db.GenerateParam, dataListParam DataListParam) (DataListResult, error)
	Execs(sqlList []string, paramsList [][]interface{}) (res int64, err error)
}

type DataListParam struct {
	Database   string                 `json:"database"`
	Table      string                 `json:"table"`
	ColumnList []*db.TableColumnModel `json:"columnList"`
	Wheres     []*db.Where            `json:"wheres"`
	PageIndex  int                    `json:"pageIndex"`
	PageSize   int                    `json:"pageSize"`
	Orders     []*db.Order            `json:"orders"`
}

type DataListResult struct {
	Sql      string                   `json:"sql"`
	Total    int                      `json:"total"`
	Params   []interface{}            `json:"params"`
	DataList []map[string]interface{} `json:"dataList"`
}
