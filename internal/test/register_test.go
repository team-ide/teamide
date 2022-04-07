package test

import (
	"fmt"
	"teamide/internal/module/module_register"
	"testing"
)

func TestRegisterMysql(t *testing.T) {
	service := module_register.NewRegisterService(getMysqlServerContext())

	testRegister(service)
}

func TestRegisterSqlite(t *testing.T) {
	service := module_register.NewRegisterService(getSqliteServerContext())

	testRegister(service)
}

func testRegister(service *module_register.RegisterService) {
	var err error

	register := &module_register.RegisterModel{
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
