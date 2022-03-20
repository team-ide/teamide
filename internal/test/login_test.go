package test

import (
	"fmt"
	"teamide/internal/module/module_login"
	"testing"
)

func TestLoginMysql(t *testing.T) {
	service := module_login.NewLoginService(getMysqlDBWorker())

	testLogin(service)
}

func TestLoginSqlite(t *testing.T) {
	service := module_login.NewLoginService(getSqliteDBWorker())

	testLogin(service)
}

func testLogin(service *module_login.LoginService) {
	var err error

	login := &module_login.LoginModel{
		Account:    "zhangsan",
		Password:   "123456",
		SourceType: 1,
	}

	_, err = service.Login(login)
	if err != nil {
		panic(err)
	}
	fmt.Printf(fmt.Sprint(login))

	_, err = service.Logout(login.LoginId)
	if err != nil {
		panic(err)
	}
	fmt.Printf(fmt.Sprint(login))
}
