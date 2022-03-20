package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DatabaseSqlite Sqlite操作
type DatabaseSqlite struct {
	config DatabaseConfig
	db     *sqlx.DB
}

// NewDatabaseSqlite 根据Sqlite配置创建DatabaseSqlite
func NewDatabaseSqlite(config DatabaseConfig) (res *DatabaseSqlite, err error) {
	res = &DatabaseSqlite{
		config: config,
	}
	err = res.init()
	return
}

func (this_ *DatabaseSqlite) init() (err error) {
	var db *sqlx.DB
	db, err = sqlx.Open("sqlite3", this_.config.Database)
	if err != nil {
		return
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	this_.db = db
	return
}

func (this_ *DatabaseSqlite) GetConfig() (config DatabaseConfig) {
	config = this_.config
	return
}

func (this_ *DatabaseSqlite) Open() (err error) {
	err = this_.db.Ping()
	return
}

func (this_ *DatabaseSqlite) Close() (err error) {
	err = this_.db.Close()
	return
}

func (this_ *DatabaseSqlite) Exec(sql string, params []interface{}) (rowsAffected int64, err error) {
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

func (this_ *DatabaseSqlite) Count(sql string, params []interface{}) (count int64, err error) {
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

func (this_ *DatabaseSqlite) Query(sql string, params []interface{}, columnTypes map[string]string) (list []map[string]interface{}, err error) {
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
