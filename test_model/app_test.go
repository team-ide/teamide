package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"teamide/application"
	"teamide/application/common"
	"testing"
)

func TestServer(t *testing.T) {
	initApplication()
	err := _app.StartServers()
	if err != nil {
		panic(err)
	}
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

func TestApplication(t *testing.T) {
	initApplication()
	createAllTable()

	for _, test := range _app.GetContext().Tests {
		invokeTest(test)
	}
}

func TestServiceJavascript(t *testing.T) {
	initApplication()
	var javascript string
	for _, one := range _app.GetContext().Services {
		javascript_, err := common.GetServiceJavascriptByService(_app, one)
		if err != nil {
			panic(err)
		}
		javascript += javascript_
		javascript += "\n"
		javascript += "\n"
	}
	fmt.Println("service javascript:")
	fmt.Println(javascript)
}

func TestDoc(t *testing.T) {
	initApplication()

	var err error

	loader := &application.ContextDoc{Dir: "app", App: _app}

	err = loader.Out()
	if err != nil {
		panic(err)
	}
}
func TestDDL(t *testing.T) {
	initApplication()

	createTable(_app.GetContext().GetDatasourceDatabase("oracle"))
}
