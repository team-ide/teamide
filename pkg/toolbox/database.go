package toolbox

import (
	"encoding/json"
	"fmt"
	"strings"
	"teamide/pkg/form"
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
	Database     string            `json:"database"`
	Table        string            `json:"table"`
	Columns      []TableColumnInfo `json:"columns"`
	Wheres       []Where           `json:"wheres"`
	PageIndex    int               `json:"pageIndex"`
	PageSize     int               `json:"pageSize"`
	DatabaseType string            `json:"databaseType"`
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
		var tables []TableDetailInfo
		tables, err = service.TableDetails(request.Database, request.Table)
		if err != nil {
			return
		}
		if len(tables) > 0 {
			res["table"] = tables[0]
		}
	case "ddl":
		var databaseType string = request.DatabaseType
		if databaseType == "" {
			databaseType = databaseConfig.Type
		}
		var sqls []string
		if request.Table == "" {
			var sqls_ []string
			sqls_, err = ToDatabaseDDL(request.Database, databaseType)
			if err != nil {
				return
			}
			sqls = append(sqls, sqls_...)
		}

		var tables []TableDetailInfo
		tables, err = service.TableDetails(request.Database, request.Table)
		if err != nil {
			return
		}
		for _, table := range tables {
			var sqls_ []string
			sqls_, err = ToTableDDL(databaseType, table)
			if err != nil {
				return
			}
			sqls = append(sqls, sqls_...)
		}

		res["sqls"] = sqls
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
func DatabaseIsMySql(databaseType string) bool {
	return strings.ToLower(databaseType) == "mysql"
}
func DatabaseIsOracle(databaseType string) bool {
	return strings.ToLower(databaseType) == "oracle"
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
	res.SetLastUseTime()
	return
}

func CreateDatabaseService(config DatabaseConfig) (service DatabaseService, err error) {
	service, err = CreateMysqlService(config)
	return
}

type DatabaseService interface {
	GetWaitTime() int64
	GetLastUseTime() int64
	SetLastUseTime()
	Stop()
	Open() error
	Databases() ([]DatabaseInfo, error)
	Tables(database string) ([]TableInfo, error)
	TableDetails(database string, table string) ([]TableDetailInfo, error)
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
	Name       string `json:"name"`
	Comment    string `json:"comment"`
	Type       string `json:"type"`
	Length     int    `json:"length"`
	Decimal    int    `json:"decimal"`
	PrimaryKey bool   `json:"primaryKey"`
	NotNull    bool   `json:"notNull"`
	Default    string `json:"default"`
}

type TableIndexInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Columns string `json:"columns"`
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
