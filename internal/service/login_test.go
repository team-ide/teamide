package service

import (
	"fmt"
	"teamide/internal/model"
	"testing"
)

func TestLoginMysql(t *testing.T) {
	service := NewLoginService(getMysqlDBWorker())

	testLogin(service)
}

func TestLoginSqlite(t *testing.T) {
	service := NewLoginService(getSqliteDBWorker())

	testLogin(service)
}

func testLogin(service *LoginService) {
	var err error

	login := &model.LoginModel{
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
