package toolbox

import (
	"encoding/json"
	"fmt"
	"teamide/pkg/db"
	"teamide/pkg/db/task"
	"teamide/pkg/form"
	"teamide/pkg/util"
)

func databaseWorker() *Worker {

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

	return worker_
}

type DatabaseBaseRequest struct {
	Database     string                 `json:"database"`
	Table        string                 `json:"table"`
	TaskKey      string                 `json:"taskKey"`
	ExecuteSQL   string                 `json:"executeSQL"`
	ColumnList   []*db.TableColumnModel `json:"columnList"`
	Wheres       []*db.Where            `json:"wheres"`
	Orders       []*db.Order            `json:"orders"`
	PageIndex    int                    `json:"pageIndex"`
	PageSize     int                    `json:"pageSize"`
	DatabaseType string                 `json:"databaseType"`
}

func databaseWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {
	var service *db.Service

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

	var generateParam = &db.GenerateParam{}
	err = json.Unmarshal(dataBS, generateParam)
	if err != nil {
		return
	}

	if generateParam.DatabaseType == "" {
		generateParam.DatabaseType = databaseConfig.Type
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
	case "createDatabase":
		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		var database = &db.DatabaseModel{}
		err = json.Unmarshal(dataBS, database)
		if err != nil {
			return
		}

		var sqlList []string
		sqlList, err = db.ToDatabaseDDL(generateParam, database)
		if err != nil {
			return
		}
		if len(sqlList) > 0 {
			for _, sql := range sqlList {
				_, err = service.GetDatabaseWorker().Exec(sql, nil)
				if err != nil {
					return
				}
			}
		}
	case "createDatabaseSql":
		var database = &db.DatabaseModel{}
		err = json.Unmarshal(dataBS, database)
		if err != nil {
			return
		}

		var sqlList []string
		sqlList, err = db.ToDatabaseDDL(generateParam, database)
		if err != nil {
			return
		}
		res["sqlList"] = sqlList
	case "createTable":
		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		var table = &db.TableModel{}
		err = json.Unmarshal(dataBS, table)
		if err != nil {
			return
		}
		var sqlList []string
		sqlList, err = db.ToTableDDL(generateParam, request.Database, table)
		if err != nil {
			return
		}
		if len(sqlList) > 0 {
			for _, sql := range sqlList {
				_, err = service.GetDatabaseWorker().Exec(sql, nil)
				if err != nil {
					// 建表出现异常，删除SQL
					sqlList, _ = db.ToTableDeleteDDL(generateParam, request.Database, table.Name)
					if err != nil {
						return
					}
					if len(sqlList) > 0 {
						for _, sql = range sqlList {
							_, _ = service.GetDatabaseWorker().Exec(sql, nil)
						}
					}
					return
				}
			}
		}
	case "createTableSql":
		var table = &db.TableModel{}
		err = json.Unmarshal(dataBS, table)
		if err != nil {
			return
		}
		var sqlList []string
		sqlList, err = db.ToTableDDL(generateParam, request.Database, table)
		if err != nil {
			return
		}
		res["sqlList"] = sqlList

	case "updateTable":
		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		var table = &db.TableModel{}
		err = json.Unmarshal(dataBS, table)
		if err != nil {
			return
		}
		var sqlList []string
		sqlList, err = db.ToTableUpdateDDL(generateParam, request.Database, table)
		if err != nil {
			return
		}
		if len(sqlList) > 0 {
			for _, sql := range sqlList {
				_, err = service.GetDatabaseWorker().Exec(sql, nil)
				if err != nil {
					return
				}
			}
		}
	case "updateTableSql":
		var table = &db.TableModel{}
		err = json.Unmarshal(dataBS, table)
		if err != nil {
			return
		}
		var sqlList []string
		sqlList, err = db.ToTableUpdateDDL(generateParam, request.Database, table)
		if err != nil {
			return
		}
		res["sqlList"] = sqlList
	case "deleteDatabase":
		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		var sqlList []string
		sqlList, err = db.ToDatabaseDeleteDDL(generateParam, request.Database)
		if err != nil {
			return
		}
		if len(sqlList) > 0 {
			for _, sql := range sqlList {
				_, err = service.GetDatabaseWorker().Exec(sql, nil)
				if err != nil {
					return
				}
			}
		}
	case "deleteTable":
		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		var sqlList []string
		sqlList, err = db.ToTableDeleteDDL(generateParam, request.Database, request.Table)
		if err != nil {
			return
		}
		if len(sqlList) > 0 {
			for _, sql := range sqlList {
				_, err = service.GetDatabaseWorker().Exec(sql, nil)
				if err != nil {
					return
				}
			}
		}
	case "ddl":

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

		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		var dataListRequest db.DataListResult
		dataListRequest, err = service.DataList(generateParam, request.Database, request.Table, request.ColumnList, request.Wheres, request.Orders, request.PageSize, request.PageIndex)
		if err != nil {
			return
		}
		res["sql"] = dataListRequest.Sql
		res["params"] = dataListRequest.Params
		res["total"] = dataListRequest.Total
		res["dataList"] = dataListRequest.DataList
	case "executeSQL":

		executeSQLTask := &task.ExecuteSQLTask{
			Database:      request.Database,
			ExecuteSQL:    request.ExecuteSQL,
			Service:       service,
			GenerateParam: generateParam,
		}
		executeSQLTask.Start()
		res["task"] = executeSQLTask

	case "dataListSql":

		var saveDataTask = &task.SaveDataTask{}
		err = json.Unmarshal(dataBS, saveDataTask)
		if err != nil {
			return
		}

		saveDataTask.Service = service
		saveDataTask.GenerateParam = generateParam

		var sqlList []string
		var valuesList [][]interface{}

		sqlList, valuesList, err = saveDataTask.SaveDataListSql()
		if err != nil {
			return
		}
		res["sqlList"] = sqlList
		res["valuesList"] = valuesList
	case "saveDataList":

		generateParam.OpenTransaction = true
		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		var saveDataTask = &task.SaveDataTask{}
		err = json.Unmarshal(dataBS, saveDataTask)
		if err != nil {
			return
		}

		saveDataTask.Service = service
		saveDataTask.GenerateParam = generateParam

		err = saveDataTask.Start()
		if err != nil {
			return
		}
		res["task"] = saveDataTask

	case "import":
		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		taskKey := util.GenerateUUID()

		var importTask = &task.ImportTask{}
		err = json.Unmarshal(dataBS, importTask)
		if err != nil {
			return
		}

		importTask.Key = taskKey
		importTask.Service = service
		importTask.GenerateParam = generateParam

		task.StartImportTask(importTask)

		res["taskKey"] = taskKey
	case "importStatus":
		importTask := task.GetImportTask(request.TaskKey)
		res["task"] = importTask
	case "importStop":
		task.StopImportTask(request.TaskKey)
	case "importClean":
		task.CleanImportTask(request.TaskKey)

	case "export":
		generateParam.AppendDatabase = true
		generateParam.DatabasePackingCharacter = "`"
		generateParam.TablePackingCharacter = "`"
		generateParam.ColumnPackingCharacter = "`"

		taskKey := util.GenerateUUID()

		var exportTask = &task.ExportTask{}
		err = json.Unmarshal(dataBS, exportTask)
		if err != nil {
			return
		}

		exportTask.Key = taskKey
		exportTask.Service = service
		exportTask.GenerateParam = generateParam

		task.StartExportTask(exportTask)

		res["taskKey"] = taskKey
	case "exportStatus":
		exportTask := task.GetExportTask(request.TaskKey)
		res["task"] = exportTask
	case "exportStop":
		task.StopExportTask(request.TaskKey)
	case "exportClean":
		task.CleanExportTask(request.TaskKey)
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

func getDatabaseService(config db.DatabaseConfig) (res *db.Service, err error) {
	key := fmt.Sprint("database-", config.Type, "-", config.Host, "-", config.Port)
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
