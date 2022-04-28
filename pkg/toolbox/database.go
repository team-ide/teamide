package toolbox

import (
	"encoding/json"
	"fmt"
)

func init() {
	worker_ := &Worker{
		Name: "database",
		Text: "Database",
		Work: databaseWork,
	}

	AddWorker(worker_)
}

type DatabaseBaseRequest struct {
	Database  string            `json:"database"`
	Table     string            `json:"table"`
	Columns   []TableColumnInfo `json:"columns"`
	Wheres    []Where           `json:"wheres"`
	PageIndex int               `json:"pageIndex"`
	PageSize  int               `json:"pageSize"`
}

func databaseWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {
	var service DatabaseService

	var databaseConfig DatabaseConfig
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

	bs, err = json.Marshal(data)
	if err != nil {
		return
	}
	request := &DatabaseBaseRequest{}
	err = json.Unmarshal(bs, request)
	if err != nil {
		return
	}

	res = map[string]interface{}{}
	switch work {
	case "checkConnect":
		err = service.Open()
		if err != nil {
			return
		}
	case "databases":
		var databases []DatabaseInfo
		databases, err = service.Databases()
		if err != nil {
			return
		}
		res["databases"] = databases
	case "tables":
		var tables []TableInfo
		tables, err = service.Tables(request.Database)
		if err != nil {
			return
		}
		res["tables"] = tables
	case "showCreateDatabase":
		var create string
		create, err = service.ShowCreateDatabase(request.Database)
		if err != nil {
			return
		}
		res["create"] = create
	case "showCreateTable":
		var create string
		create, err = service.ShowCreateTable(request.Database, request.Table)
		if err != nil {
			return
		}
		res["create"] = create
	case "tableDetail":
		var table TableDetailInfo
		table, err = service.TableDetail(request.Database, request.Table)
		if err != nil {
			return
		}
		res["table"] = table
	case "datas":
		var datasRequest DatasResult
		datasRequest, err = service.Datas(DatasParam{
			Database:  request.Database,
			Table:     request.Table,
			Columns:   request.Columns,
			Wheres:    request.Wheres,
			PageIndex: request.PageIndex,
			PageSize:  request.PageSize,
		})
		if err != nil {
			return
		}
		res["sql"] = datasRequest.Sql
		res["params"] = datasRequest.Params
		res["total"] = datasRequest.Total
		res["datas"] = datasRequest.Datas
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

func getDatabaseService(config DatabaseConfig) (res DatabaseService, err error) {
	key := fmt.Sprint("database-", config.Type, "-", config.Host, "-", config.Port, "-", config.Username, "-", config.Password)
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s DatabaseService
		s, err = CreateDatabaseService(config)
		if err != nil {
			return
		}
		_, err = s.Databases()
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
	return
}

func CreateDatabaseService(config DatabaseConfig) (service DatabaseService, err error) {
	service, err = CreateMysqlService(config)
	return
}

type DatabaseService interface {
	GetWaitTime() int64
	GetLastUseTime() int64
	Stop()
	Open() error
	Databases() ([]DatabaseInfo, error)
	Tables(database string) ([]TableInfo, error)
	TableDetail(database string, table string) (TableDetailInfo, error)
	ShowCreateDatabase(database string) (string, error)
	ShowCreateTable(database string, table string) (string, error)
	Datas(datasParam DatasParam) (DatasResult, error)
}

type DatabaseConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type DatabaseInfo struct {
	Name string `json:"name"`
}

type TableInfo struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type TableDetailInfo struct {
	Name    string            `json:"name"`
	Comment string            `json:"comment"`
	Columns []TableColumnInfo `json:"columns"`
	Indexs  []TableIndexInfo  `json:"indexs"`
}

type TableColumnInfo struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Type    string `json:"type"`
	Length  int    `json:"length"`
	Decimal int    `json:"decimal"`
}

type TableIndexInfo struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type DatasParam struct {
	Database  string            `json:"database"`
	Table     string            `json:"table"`
	Columns   []TableColumnInfo `json:"columns"`
	Wheres    []Where           `json:"wheres"`
	PageIndex int               `json:"pageIndex"`
	PageSize  int               `json:"pageSize"`
}

type DatasResult struct {
	Sql    string                   `json:"sql"`
	Total  string                   `json:"total"`
	Params []interface{}            `json:"params"`
	Datas  []map[string]interface{} `json:"datas"`
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
