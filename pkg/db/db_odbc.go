//go:build !darwin

package db

import (
	"database/sql"
	"github.com/team-ide/go-driver/db_odbc"
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
