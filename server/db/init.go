package db

import (
	"fmt"
	"server/config"
)

var (
	DBService MysqlService
)

func init() {
	var service interface{}
	var err error
	databaseConfig := DatabaseConfig{
		Host:     config.Config.Mysql.Host,
		Port:     config.Config.Mysql.Port,
		Database: config.Config.Mysql.Database,
		Username: config.Config.Mysql.Username,
		Password: config.Config.Mysql.Password,
	}
	fmt.Println("数据库初始化：host:", databaseConfig.Host, ",port:", databaseConfig.Port, ",database:", databaseConfig.Database)
	service, err = CreateMysqlService(databaseConfig)
	if err != nil {
		panic(err)
	}
	DBService = *service.(*MysqlService)

	_, err = DBService.Exec(SqlParam{
		Sql:    "SELECT 1 FROM " + TABLE_INSTALL,
		Params: []interface{}{},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("数据库连接成功！")
}

func Init() {
}
