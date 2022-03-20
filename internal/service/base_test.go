package service

import "teamide/pkg/db"

var (
	mysqlConfig = db.DatabaseConfig{
		Type:     "mysql",
		Database: "TEAM_IDE_TEST",
		Host:     "mysql",
		Port:     3306,
		Username: "root",
		Password: "123456",
	}

	sqliteConfig = db.DatabaseConfig{
		Type:     "sqlite",
		Database: "TEAM_IDE_TEST",
	}
)

func getMysqlDBWorker() db.DatabaseWorker {

	dbWorker, err := db.NewDatabaseWorker(mysqlConfig)
	if err != nil {
		panic(err)
	}
	return dbWorker
}

func getSqliteDBWorker() db.DatabaseWorker {

	dbWorker, err := db.NewDatabaseWorker(sqliteConfig)
	if err != nil {
		panic(err)
	}
	return dbWorker
}
