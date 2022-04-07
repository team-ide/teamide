package test

import (
	"teamide/internal/config"
	"teamide/internal/context"
	"teamide/pkg/db"
)

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

func getServerContext(configPath string) (serverContext *context.ServerContext) {
	serverContext = &context.ServerContext{
		IsStandAlone: true,
		IsHtmlDev:    true,
		RootDir:      "./",
		UserHomeDir:  "./",
	}
	serverContext.HttpAesKey = "Q56hFAauWk18Gy2i"
	var serverConfig *config.ServerConfig
	serverConfig, err := config.CreateServerConfig(configPath)
	if err != nil {
		panic(err)
		return
	}
	//context.ServerConf = serverConf
	serverContext.ServerConfig = serverConfig
	err = serverContext.Init(serverConfig)
	if err != nil {
		panic(err)
		return
	}
	serverContext.Decryption, err = context.NewDefaultDecryption(serverContext.Logger)
	if err != nil {
		panic(err)
		return
	}
	return
}

func getMysqlServerContext() *context.ServerContext {
	serverContext := getServerContext("./conf/mysql.yaml")
	return serverContext
}

func getSqliteServerContext() *context.ServerContext {
	serverContext := getServerContext("./conf/sqlite.yaml")
	return serverContext
}
