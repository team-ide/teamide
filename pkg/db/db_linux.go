//go:build linux

package db

import (
	"database/sql"
	"github.com/team-ide/go-driver/db_odbc"
	"github.com/team-ide/go-driver/db_oracle"
	"github.com/team-ide/go-driver/db_shentong"
)

func initOdbcDatabase() {
	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			dsn := db_odbc.GetDSN(config.OdbcName, config.Username, config.Password)
			db, err = db_odbc.Open(dsn)
			return
		},
		DialectName: db_odbc.GetDialect(),
		matches:     []string{"odbc"},
	})
}

func initOracleDatabase() {
	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			dsn := db_oracle.GetDSN(config.Username, config.Password, config.Host, config.Port, config.Sid)
			db, err = db_oracle.Open(dsn)
			return
		},
		DialectName: db_oracle.GetDialect(),
		matches:     []string{"oracle"},
	})
}

func initShenTongDatabase() {
	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			dsn := db_shentong.GetDSN(config.Username, config.Password, config.Host, config.Port, config.DbName)
			db, err = db_shentong.Open(dsn)
			return
		},
		DialectName: db_shentong.GetDialect(),
		matches:     []string{"ShenTong", "st"},
	})
}
