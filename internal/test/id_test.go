package test

import (
	"fmt"
	"teamide/internal/module/module_id"
	"testing"
)

func TestIDMysql(t *testing.T) {
	var err error
	service := module_id.NewIDService(getMysqlServerContext())
	id, err := service.GetNextID(module_id.IDTypeUser)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ID:%d\n", id)
}

func TestIDSqlite(t *testing.T) {
	var err error
	service := module_id.NewIDService(getSqliteServerContext())
	id, err := service.GetNextID(module_id.IDTypeUser)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ID:%d\n", id)
}
