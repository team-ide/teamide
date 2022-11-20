//go:build !darwin

package db

import (
	"database/sql"
	"github.com/team-ide/go-driver/db_oracle"
)

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
