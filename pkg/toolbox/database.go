package toolbox

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"go.uber.org/zap"
	"strings"
	"teamide/pkg/db"
	"teamide/pkg/form"
	"teamide/pkg/javascript"
	"teamide/pkg/util"
	"time"
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
	ColumnList      []*db.TableColumnModel   `json:"columnList"`
	Wheres          []Where                  `json:"wheres"`
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
		var dataListRequest DataListResult
		dataListRequest, err = service.DataList(DataListParam{
			Database:   request.Database,
			Table:      request.Table,
			ColumnList: request.ColumnList,
			Wheres:     request.Wheres,
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
		taskKey := util.GenerateUUID()
		importDataForStrategyTask := &importDataForStrategyTask{
			Key:            taskKey,
			Database:       request.Database,
			Table:          request.Table,
			ColumnList:     request.ColumnList,
			ImportDataList: request.ImportDataList,
			service:        service,
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

var (
	importDataForStrategyTaskCache = map[string]*importDataForStrategyTask{}
)

func addImportDataForStrategyTask(task *importDataForStrategyTask) {
	importDataForStrategyTaskCache[task.Key] = task
	go task.Start()
}

type importDataForStrategyTask struct {
	Key               string                   `json:"key,omitempty"`
	Database          string                   `json:"database,omitempty"`
	Table             string                   `json:"table,omitempty"`
	ColumnList        []*db.TableColumnModel   `json:"columnList,omitempty"`
	ImportDataList    []map[string]interface{} `json:"importDataList,omitempty"`
	ImportBatchNumber int                      `json:"importBatchNumber,omitempty"`
	DataCount         int                      `json:"dataCount"`
	ReadyDataCount    int                      `json:"readyDataCount"`
	ImportSuccess     int                      `json:"importSuccess"`
	ImportError       int                      `json:"importError"`
	IsEnd             bool                     `json:"isEnd,omitempty"`
	StartTime         time.Time                `json:"startTime,omitempty"`
	EndTime           time.Time                `json:"endTime,omitempty"`
	Error             string                   `json:"error,omitempty"`
	UseTime           int64                    `json:"useTime"`
	IsStop            bool                     `json:"isStop"`
	service           DatabaseService
}

func (this_ *importDataForStrategyTask) Stop() {
	this_.IsStop = true
}
func (this_ *importDataForStrategyTask) Start() {
	this_.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			Logger.Error("根据策略导入数据异常", zap.Any("error", err))
			this_.Error = fmt.Sprint(err)
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
	}()

	for _, importData := range this_.ImportDataList {
		importCount := 0
		if importData["_$importCount"] != nil {
			importCount = int(importData["_$importCount"].(float64))
		}
		if importCount <= 0 {
			importCount = 0
		}
		importData["_$importCount"] = importCount
		this_.DataCount += importCount
	}

	for _, importData := range this_.ImportDataList {
		if this_.IsStop {
			break
		}
		err := this_.importData(this_.Database, this_.Table, this_.ColumnList, importData)
		if err != nil {
			panic(err)
		}
	}
}
func (this_ *importDataForStrategyTask) importData(database, table string, columnList []*db.TableColumnModel, importData map[string]interface{}) (err error) {
	importCount := importData["_$importCount"].(int)
	if importCount <= 0 {
		return
	}
	if this_.IsStop {
		return
	}

	var dataList []map[string]interface{}
	importBatchNumber := this_.ImportBatchNumber
	if importBatchNumber <= 0 {
		importBatchNumber = 10
	}
	scriptContext := javascript.GetContext()

	vm := goja.New()

	for key, value := range scriptContext {
		err = vm.Set(key, value)
		if err != nil {
			return
		}
	}

	for i := 0; i < importCount; i++ {
		data := map[string]interface{}{}
		err = vm.Set("_$index", i)
		if err != nil {
			return
		}

		for _, column := range columnList {

			if this_.IsStop {
				return
			}

			value, valueOk := importData[column.Name]
			if !valueOk {
				continue
			}
			valueString, valueStringOk := value.(string)
			if valueStringOk && valueString != "" {
				var scriptValue goja.Value
				scriptValue, err = vm.RunString(valueString)
				if err != nil {
					Logger.Error("表达式执行异常", zap.Any("script", valueString), zap.Error(err))
					return
				}
				value = scriptValue.Export()
			}
			data[column.Name] = value
			err = vm.Set(column.Name, value)

			if err != nil {
				return
			}
		}
		this_.ReadyDataCount++
		dataList = append(dataList, data)
		if len(dataList) >= importBatchNumber {

			if this_.IsStop {
				return
			}
			err = this_.doImportData(database, table, columnList, dataList)
			if err != nil {
				this_.ImportError += len(dataList)
				return
			} else {
				this_.ImportSuccess += len(dataList)
			}
			dataList = []map[string]interface{}{}
		}
	}
	err = this_.doImportData(database, table, columnList, dataList)
	if err != nil {
		this_.ImportError += len(dataList)
		return
	} else {
		this_.ImportSuccess += len(dataList)
	}
	return
}

func (this_ *importDataForStrategyTask) doImportData(database, table string, columnList []*db.TableColumnModel, dataList []map[string]interface{}) (err error) {

	if len(dataList) == 0 {
		return
	}
	var sqlList []string
	var paramsList [][]interface{}
	for _, data := range dataList {
		var sqlParam *SqlParam
		sqlParam, err = this_.getImportDataFinder(database, table, columnList, data)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlParam.Sql)
		paramsList = append(paramsList, sqlParam.Params)
	}

	_, err = this_.service.Execs(sqlList, paramsList)
	if err != nil {
		return
	}
	return
}

func (this_ *importDataForStrategyTask) getImportDataFinder(database, table string, columnList []*db.TableColumnModel, data map[string]interface{}) (sqlParam *SqlParam, err error) {

	insertColumns := ""
	insertValues := ""
	var values []interface{}
	for _, column := range columnList {
		value, valueOk := data[column.Name]
		if !valueOk {
			continue
		}
		insertColumns += column.Name + ","
		insertValues += "?,"
		values = append(values, value)
	}
	insertColumns = strings.TrimSuffix(insertColumns, ",")
	insertValues = strings.TrimSuffix(insertValues, ",")
	sql := "INSERT INTO " + database + "." + table + ""
	sql += "(" + insertColumns + ")"
	sql += " VALUES "
	sql += "(" + insertValues + ")"

	sqlParam = &SqlParam{
		Sql:    sql,
		Params: values,
	}
	return
}

type saveDataListTask struct {
	Key             string                 `json:"key,omitempty"`
	Database        string                 `json:"database,omitempty"`
	Table           string                 `json:"table,omitempty"`
	ColumnList      []*db.TableColumnModel `json:"columnList,omitempty"`
	generateParam   *db.GenerateParam
	UpdateList      []map[string]interface{} `json:"-"`
	UpdateWhereList []map[string]interface{} `json:"-"`
	InsertList      []map[string]interface{} `json:"-"`
	DeleteList      []map[string]interface{} `json:"-"`
	SaveBatchNumber int                      `json:"saveBatchNumber,omitempty"`
	DataCount       int                      `json:"dataCount"`
	ReadyDataCount  int                      `json:"readyDataCount"`
	SaveSuccess     int                      `json:"saveSuccess"`
	SaveError       int                      `json:"saveError"`
	IsEnd           bool                     `json:"isEnd,omitempty"`
	StartTime       time.Time                `json:"startTime,omitempty"`
	EndTime         time.Time                `json:"endTime,omitempty"`
	Error           string                   `json:"error,omitempty"`
	UseTime         int64                    `json:"useTime"`
	IsStop          bool                     `json:"isStop"`
	service         DatabaseService
}

func (this_ *saveDataListTask) Stop() {
	this_.IsStop = true
}

func (this_ *saveDataListTask) Start() (err error) {
	this_.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			Logger.Error("根据保存数据异常", zap.Any("error", err))
			this_.Error = fmt.Sprint(err)
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
	}()
	this_.DataCount += len(this_.InsertList)
	this_.DataCount += len(this_.UpdateList)
	this_.DataCount += len(this_.DeleteList)

	saveBatchNumber := this_.SaveBatchNumber
	if saveBatchNumber <= 0 {
		saveBatchNumber = 10
	}

	var sqlList []string
	var valuesList [][]interface{}

	sqlList, valuesList, err = this_.SaveDataListSql()
	if err != nil {
		return
	}
	res, err := this_.service.Execs(sqlList, valuesList)
	if err != nil {
		this_.SaveError += len(sqlList)
		return
	}
	this_.SaveSuccess += int(res)
	return
}

func (this_ *saveDataListTask) SaveDataListSql() (sqlList []string, valuesList [][]interface{}, err error) {

	var sqlList_ []string
	var valuesList_ [][]interface{}
	if len(this_.InsertList) > 0 {
		sqlList_, valuesList_, err = db.DataListInsertSql(this_.generateParam, this_.Database, this_.Table, this_.ColumnList, this_.InsertList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	if len(this_.UpdateList) > 0 {
		sqlList_, valuesList_, err = db.DataListUpdateSql(this_.generateParam, this_.Database, this_.Table, this_.ColumnList, this_.UpdateList, this_.UpdateWhereList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
	}
	if len(this_.DeleteList) > 0 {
		sqlList_, valuesList_, err = db.DataListDeleteSql(this_.generateParam, this_.Database, this_.Table, this_.ColumnList, this_.DeleteList)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
		valuesList = append(valuesList, valuesList_...)
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
	Databases() ([]*db.DatabaseModel, error)
	Tables(database string) ([]*db.TableModel, error)
	TableDetails(database string, table string) ([]*db.TableModel, error)
	DataList(dataListParam DataListParam) (DataListResult, error)
	Execs(sqlList []string, paramsList [][]interface{}) (res int64, err error)
}

type DataListParam struct {
	Database   string                 `json:"database"`
	Table      string                 `json:"table"`
	ColumnList []*db.TableColumnModel `json:"columnList"`
	Wheres     []Where                `json:"wheres"`
	PageIndex  int                    `json:"pageIndex"`
	PageSize   int                    `json:"pageSize"`
	Orders     []Order                `json:"orders"`
}

type DataListResult struct {
	Sql      string                   `json:"sql"`
	Total    int                      `json:"total"`
	Params   []interface{}            `json:"params"`
	DataList []map[string]interface{} `json:"dataList"`
}

type Where struct {
	Name                    string `json:"name"`
	Value                   string `json:"value"`
	Before                  string `json:"before"`
	After                   string `json:"after"`
	CustomSql               string `json:"customSql"`
	SqlConditionalOperation string `json:"sqlConditionalOperation"`
	AndOr                   string `json:"andOr"`
}

type Order struct {
	Name    string `json:"name"`
	DescAsc string `json:"descAsc"`
}
