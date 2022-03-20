package service

import (
	"fmt"
	"teamide/internal/model"
	"testing"
)

func TestIDMysql(t *testing.T) {
	var err error
	service := NewIDService(getMysqlDBWorker())
	id, err := service.GetNextID(model.IDTypeUser)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ID:%d\n", id)
}

func TestIDSqlite(t *testing.T) {
	var err error
	service := NewIDService(getSqliteDBWorker())
	id, err := service.GetNextID(model.IDTypeUser)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ID:%d\n", id)
}
