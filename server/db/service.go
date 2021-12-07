package db

import (
	"server/base"
	"server/worker"
	"strings"
)

type Automatic struct {
	databaseConfig    DatabaseConfig
	automaticShutdown *worker.AutomaticShutdown
}

func (automatic *Automatic) CreateAutomaticShutdown(automaticShutdown *worker.AutomaticShutdown) error {
	service, err := CreateService(automatic.databaseConfig)
	if err != nil {
		return err
	}
	dbService := service.(interface{ DatabaseService })
	err = dbService.Open()
	if err != nil {
		dbService.Close()
		return err
	}
	// 默认10分钟自动关闭
	automaticShutdown.AutomaticShutdown = 10 * 60
	automaticShutdown.Service = service
	automaticShutdown.Stop = func() {
		dbService.Close()
	}
	automatic.automaticShutdown = automaticShutdown

	return err
}

func GetService(databaseConfig DatabaseConfig) (service DatabaseService, err error) {
	automatic := &Automatic{
		databaseConfig: databaseConfig,
	}
	key := "database-" + base.ToJSON(databaseConfig)
	var automaticShutdown *worker.AutomaticShutdown
	automaticShutdown, err = worker.GetAutomaticShutdown(key, automatic.CreateAutomaticShutdown)
	if err != nil {
		return
	}
	automaticShutdown.LastUseTimestamp = base.GetNowTime()
	service = automaticShutdown.Service.(interface{ DatabaseService })
	return
}

func CreateService(databaseConfig DatabaseConfig) (service interface{}, err error) {
	if strings.ToLower(databaseConfig.Type) == "mysql" {
		service, err = CreateMysqlService(databaseConfig)
	}
	return
}

type DatabaseService interface {
	Open() error
	Close() error
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
	Database string `json:"database"`
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
