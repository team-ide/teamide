package service

import (
	"fmt"
	"teamide/internal/model"
	"testing"
)

func TestRegisterMysql(t *testing.T) {
	service := NewRegisterService(getMysqlDBWorker())

	testRegister(service)
}

func TestRegisterSqlite(t *testing.T) {
	service := NewRegisterService(getSqliteDBWorker())

	testRegister(service)
}

func testRegister(service *RegisterService) {
	var err error

	register := &model.RegisterModel{
		Name:       "张三",
		Account:    "zhangsan",
		Email:      "zhangsan@teamide.com",
		Password:   "123456",
		SourceType: 1,
	}

	_, err = service.Register(register)
	if err != nil {
		panic(err)
	}
	fmt.Printf(fmt.Sprint(register))
}
