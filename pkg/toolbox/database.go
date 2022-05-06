package toolbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
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
		var table TableDetailInfo
		table, err = service.TableDetail(request.Database, request.Table)
		if err != nil {
			return
		}
		res["table"] = table
	case "tableDDL":
		var table TableDetailInfo
		table, err = service.TableDetail(request.Database, request.Table)
		if err != nil {
			return
		}
		var databaseType string = request.DatabaseType
		if databaseType == "" {
			databaseType = databaseConfig.Type
		}
		var sqls []string
		sqls, err = ToTableDDL(databaseType, table)
		if err != nil {
			return
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

const (
	CREATE_DATABASE = `CREATE DATABASE IF NOT EXISTS {database}[ CHARACTER SET '{characterSet}'][ COLLATE '{collate}']`
	CREATE_TABLE    = `CREATE TABLE IF NOT EXISTS {table} (
  [{columns}]
  [PRIMARY KEY {primaryKeys}]
  [{indexs}]
)[ ENGINE={engine}][ DEFAULT CHARSET={defaultCharset}][ COMMENT='{comment}']`
	CREATE_TABLE_COLUMN       = `{column}[ {type}][ CHARACTER SET {characterSet}][ NOT NULL][ DEFAULT '{default}'][ AUTO_INCREMENT][ COMMENT '{comment}']`
	CREATE_TABLE_INDEX        = `KEY {name} ({columns})[ COMMENT '{comment}']`
	CREATE_TABLE_INDEX_UNIQUE = `UNIQUE KEY {name} ({columns})[ COMMENT '{comment}']`

	ORACLE_CREATE_TABLE = `CREATE TABLE {table} (
		[{columns}]
		[PRIMARY KEY {primaryKeys}]
	  )`
	ORACLE_CREATE_TABLE_COLUMN = `{column}[ {type}][ DEFAULT {default}][ NOT NULL]`
)

var (
	mysqlTypeMap  = map[string]func(length int, decimal int) string{}
	oracleTypeMap = map[string]func(length int, decimal int) string{}
)

func init() {
	oracleTypeMap["varchar"] = func(length int, decimal int) string {
		return fmt.Sprintf("varchar(%d)", length)
	}
	oracleTypeMap["text"] = func(length int, decimal int) string {
		return fmt.Sprintf("text")
	}
	oracleTypeMap["int"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	oracleTypeMap["bigint"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
}

func ToTableDDL(databaseType string, table TableDetailInfo) (sqls []string, err error) {
	sqls = []string{}
	var columns string
	var primaryKeys string
	var indexs string
	var data map[string]string

	if len(table.Columns) > 0 {
		var columnSql string
		for _, one := range table.Columns {
			data = map[string]string{}
			if one.Name == "" {
				continue
			}
			if one.PrimaryKey {
				primaryKeys += "" + one.Name + ","
			}
			data["column"] = fmt.Sprint("", one.Name, "")
			data["comment"] = fmt.Sprint("", one.Comment, "")
			data["default"] = fmt.Sprint("", one.Default, "")
			if one.NotNull {
				data["NOT NULL"] = "true"
			}
			var typeFunc func(length int, decimal int) string
			if DatabaseIsMySql(databaseType) {
				typeFunc = mysqlTypeMap[strings.ToLower(one.Type)]
			} else if DatabaseIsOracle(databaseType) {
				typeFunc = oracleTypeMap[strings.ToLower(one.Type)]
			}
			if typeFunc == nil {
				err = errors.New("字段类型[" + one.Type + "]未映射!")
				return
			}
			data["type"] = typeFunc(one.Length, one.Decimal)
			if DatabaseIsMySql(databaseType) {
				columnSql, err = foramtSql(CREATE_TABLE_COLUMN, data)
			} else if DatabaseIsOracle(databaseType) {
				columnSql, err = foramtSql(ORACLE_CREATE_TABLE_COLUMN, data)
			}
			if err != nil {
				return
			}
			if columnSql == "" {
				continue
			}
			if columns != "" {
				columns += "  "
			}
			columns += columnSql + ",\n"
		}
	}

	if DatabaseIsMySql(databaseType) && len(table.Indexs) > 0 {
		var indexSql string
		for _, one := range table.Indexs {
			data = map[string]string{}
			if one.Name == "" || one.Columns == "" {
				continue
			}
			data["name"] = fmt.Sprint("", one.Name, "")
			data["columns"] = fmt.Sprint("", one.Columns, "")
			data["comment"] = fmt.Sprint("", one.Comment, "")

			switch one.Type {
			case "UNIQUE", "unique":
				indexSql, err = foramtSql(CREATE_TABLE_INDEX_UNIQUE, data)
			default:
				indexSql, err = foramtSql(CREATE_TABLE_INDEX, data)
			}

			if err != nil {
				return
			}
			if indexSql == "" {
				continue
			}
			if indexs != "" {
				indexs += "  "
			}
			indexs += indexSql + ",\n"
		}
	}
	data = map[string]string{}
	data["table"] = fmt.Sprint("", table.Name, "")

	columns = strings.TrimSuffix(columns, "\n")
	if primaryKeys == "" && indexs == "" {
		columns = strings.TrimSuffix(columns, ",")
	}
	primaryKeys = strings.TrimSuffix(primaryKeys, ",")
	if primaryKeys != "" {
		primaryKeys = "(" + primaryKeys + ")"
	}
	if indexs != "" {
		primaryKeys += ","
	}
	indexs = strings.TrimSuffix(indexs, "\n")
	indexs = strings.TrimSuffix(indexs, ",")
	data["columns"] = columns
	data["primaryKeys"] = primaryKeys
	data["indexs"] = indexs
	data["comment"] = table.Comment
	var sql string
	if DatabaseIsMySql(databaseType) {
		sql, err = foramtSql(CREATE_TABLE, data)
		if err != nil {
			return
		}
		if sql != "" {
			sqls = append(sqls, sql)
		}
	} else if DatabaseIsOracle(databaseType) {
		sql, err = foramtSql(ORACLE_CREATE_TABLE, data)
		if err != nil {
			return
		}
		if sql != "" {
			sqls = append(sqls, sql)
		}
		// 添加注释
		if table.Comment != "" {
			sqls = append(sqls, `COMMENT ON TABLE `+table.Name+` IS '`+table.Comment+`'`)
		}
		if len(table.Columns) > 0 {
			for _, one := range table.Columns {
				if one.Name == "" || one.Comment == "" {
					continue
				}
				sqls = append(sqls, `COMMENT ON COLUMN `+table.Name+`.`+one.Name+` IS '`+one.Comment+`'`)
			}
		}
	}
	return
}

func foramtSql(sql string, data map[string]string) (foramtSql string, err error) {
	var re *regexp.Regexp
	re, err = regexp.Compile(`\[(.+?)\]`)
	if err != nil {
		return
	}
	indexsList := re.FindAllIndex([]byte(sql), -1)
	var lastIndex int = 0
	var sql_ string
	var formatValueSql string
	var find bool = true
	for _, indexs := range indexsList {
		sql_ = sql[lastIndex:indexs[0]]
		formatValueSql, find = foramtValueSql(sql_, data)
		if find {
			foramtSql += formatValueSql
		}

		lastIndex = indexs[1]

		sql_ = sql[indexs[0]+1 : indexs[1]-1]

		if !strings.Contains(sql_, `{`) {
			if data[strings.TrimSpace(sql_)] != "" {
				foramtSql += sql_
			}
		} else {
			formatValueSql, find = foramtValueSql(sql_, data)
			if find {
				foramtSql += formatValueSql
			}
		}
	}
	sql_ = sql[lastIndex:]
	formatValueSql, find = foramtValueSql(sql_, data)
	if find {
		foramtSql += formatValueSql
	}
	return
}

func foramtValueSql(sql string, data map[string]string) (res string, find bool) {
	var re *regexp.Regexp
	re, _ = regexp.Compile(`{(.+?)}`)
	find = true
	indexsList := re.FindAllIndex([]byte(sql), -1)
	var lastIndex int = 0
	for _, indexs := range indexsList {
		res += sql[lastIndex:indexs[0]]

		lastIndex = indexs[1]

		key := sql[indexs[0]+1 : indexs[1]-1]
		value := data[key]
		if value == "" {
			find = false
			return
		}
		res += value
	}
	res += sql[lastIndex:]
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
