package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DatabaseMysql Mysql操作
type DatabaseMysql struct {
	config DatabaseConfig
	db     *sqlx.DB
}

// NewDatabaseMysql 根据Mysql配置创建DatabaseMysql
func NewDatabaseMysql(config DatabaseConfig) (res *DatabaseMysql, err error) {
	res = &DatabaseMysql{
		config: config,
	}
	err = res.init()
	return
}

func (this_ *DatabaseMysql) init() (err error) {
	url := fmt.Sprint(this_.config.Username, ":", this_.config.Password, "@tcp(", this_.config.Host, ":", this_.config.Port, ")/", this_.config.Database, "?charset=utf8mb4&loc=Local&parseTime=true")
	var db *sqlx.DB
	db, err = sqlx.Open("mysql", url)
	if err != nil {
		return
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	this_.db = db
	return
}

func (this_ *DatabaseMysql) GetConfig() (config DatabaseConfig) {
	config = this_.config
	return
}

func (this_ *DatabaseMysql) Open() (err error) {
	err = this_.db.Ping()
	return
}

func (this_ *DatabaseMysql) Close() (err error) {
	err = this_.db.Close()
	return
}

func (this_ *DatabaseMysql) Exec(sql string, params []interface{}) (rowsAffected int64, err error) {
	result, err := this_.db.Exec(sql, params...)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (this_ *DatabaseMysql) Count(sql string, params []interface{}) (count int64, err error) {
	rows, err := this_.db.Query(sql, params...)
	if err != nil {
		return
	}
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return
		}
	}
	rows.Close()
	return
}

func (this_ *DatabaseMysql) Query(sql string, params []interface{}, columnTypes map[string]string) (list []map[string]interface{}, err error) {
	rows, err := this_.db.Query(sql, params...)
	if err != nil {
		return
	}
	list, err = formatRowsValue(rows, columnTypes)
	if err != nil {
		return
	}
	rows.Close()
	return
}
