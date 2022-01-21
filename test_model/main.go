package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"teamide/application"
	"teamide/application/base"
	"teamide/application/coder"
	"teamide/application/common"
	"teamide/application/model"
)

var (
	done = make(chan os.Signal, 1)
)
var (
	appDir = "app"
)

func main() {
	var baseRead bool
	var isServer bool
	var isInit bool
	var isDatabase bool
	var isTest bool
	var tests []string
	var isDoc bool
	var isDocWeb bool

	for _, v := range os.Args {
		if !baseRead {
			if v == "server" {
				isServer = true
				baseRead = true
				continue
			} else if v == "test" {
				isTest = true
				baseRead = true
				continue
			} else if v == "init" {
				isInit = true
				baseRead = true
				continue
			} else if v == "doc" {
				isDoc = true
				baseRead = true
				continue
			}
		}
		if isInit && v == "database" {
			isDatabase = true
			break
		} else if isDoc && v == "web" {
			isDocWeb = true
			break
		} else {
			if strings.HasPrefix(v, "dir=") {
				dir := strings.TrimPrefix(v, "dir=")
				dir = strings.TrimSpace(dir)
				if base.IsNotEmpty(dir) {
					appDir = dir
				}
				continue
			} else {
				if isTest {
					tests = append(tests, v)
				}
			}
		}
	}
	initApplication()

	if isServer {
		err := _app.StartServers()
		if err != nil {
			panic(err)
		}
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-done
	} else if isDatabase {
		createAllTable()
	} else if isDocWeb {
		outDoc()
	} else if isTest {
		for _, one := range tests {
			testModel := _app.GetContext().GetTest(one)
			if testModel == nil {
				if _app.GetLogger() != nil {
					_app.GetLogger().Error("test model [", one, "] not defind")
				}
			} else {
				invokeTest(testModel)
			}
		}
		outDoc()
	} else {
		if _app.GetLogger() != nil {
			_app.GetLogger().Error("未识别命令")
		}
	}
}

func outDoc() {
	initApplication()

	var err error

	loader := &application.ContextDoc{Dir: appDir, App: _app}

	err = loader.Out()
	if err != nil {
		panic(err)
	}
}

var (
	_app common.IApplication
)

func initApplication() {
	initModelContext()
	_app = application.NewApplication(_modelContext, &common.LoggerDefault{OutDebug_: true})
}

func invokeTest(test *model.TestModel) {

	res, err := _app.InvokeTest(test)

	if err != nil {
		panic(err)
	}
	if res.Error != nil {
		_app.GetLogger().Info("test [", test.Name, "] error:", res.Error)
	}
	_app.GetLogger().Info("test [", test.Name, "] result count:", res.Count, ",success:", res.SuccessCount, ",error:", res.ErrorCount)
	if len(res.Infos) > 0 {
		for _, info := range res.Infos {
			if info.Error != nil {
				_app.GetLogger().Error("test [", test.Name, "] [", info.ThreadName, "] [", info.ForName, "] error:", info.Error)
			} else {
				_app.GetLogger().Info("test [", test.Name, "] [", info.ThreadName, "] [", info.ForName, "] result:", base.ToJSON(info.Result))
			}
		}
	}

}

var (
	_modelContext = &model.ModelContext{}
)

func initModelContext() {
	var err error

	loader := &application.ContextLoader{Dir: appDir}

	_modelContext, err = loader.Load()
	if err != nil {
		panic(err)
	}
}

func createAllTable() {
	for _, database := range _app.GetContext().DatasourceDatabases {
		createTable(database)
	}
}

func createTable(database *model.DatasourceDatabase) {
	var err error
	executor, err := _app.GetSqlExecutor(database.Name)
	if err != nil {
		panic(err)
	}
	err = executor.Ping()
	if err != nil {
		panic(err)
	}
	var sql string
	if database.Type == "mysql" {
		sql, err = coder.GetDatabaseDDL(database)
		if err != nil {
			panic(err)
		}
		_app.GetLogger().Info("CREATE DATABASE DDL:")
		_app.GetLogger().Info(sql)
		err = executor.ExecSqls([]string{sql})
		if err != nil {
			panic(err)
		}
	}
	var sqls []string
	sqls, err = coder.GetCreateTableSqls(_app, database)
	if err != nil {
		panic(err)
	}
	_app.GetLogger().Info("DDL:")
	for _, one := range sqls {
		_app.GetLogger().Info(one)
	}
	err = executor.ExecSqls(sqls)
	if err != nil {
		panic(err)
	}
}
