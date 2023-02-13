//go:build darwin

package db

import (
	"database/sql"
	"errors"
)

func initShenTongDatabase() {

	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			err = errors.New("darwin can not support [ShenTong] database.")
			return
		},
		DialectName: "shentong",
		matches:     []string{"ShenTong", "st"},
	})
}

func initOracleDatabase() {

	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			err = errors.New("darwin can not support [oracle] database.")
			return
		},
		DialectName: "oracle",
		matches:     []string{"oracle"},
	})
}

func initGBaseDatabase() {

	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			err = errors.New("darwin can not support [GBase] database.")
			return
		},
		DialectName: "gbase",
		matches:     []string{"gbase"},
	})
}

func initOdbcDatabase() {

	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			err = errors.New("darwin can not support [odbc] database.")
			return
		},
		DialectName: "odbc",
		matches:     []string{"odbc"},
	})
}
