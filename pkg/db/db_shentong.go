//go:build !darwin

package db

import (
	"database/sql"
	"github.com/team-ide/go-driver/db_shentong"
)

func initShenTongDatabase() {
	addDatabaseType(&DatabaseType{
		newDb: func(config *DatabaseConfig) (db *sql.DB, err error) {
			dsn := db_shentong.GetDSN(config.Username, config.Password, config.Host, config.Port, config.Database)
			db, err = db_shentong.Open(dsn)
			return
		},
		DialectName: db_shentong.GetDialect(),
		matches:     []string{"ShenTong", "st"},
	})
}
