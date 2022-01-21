package main

import (
	"testing"
)

func TestServiceLogin(t *testing.T) {
	initApplication()
	createAllTable()

	invokeTest(_app.GetContext().GetTest("login/login"))
}

func TestServiceTokenCreate(t *testing.T) {
	initApplication()
	createAllTable()

	invokeTest(_app.GetContext().GetTest("token/create"))
}

func TestServiceTokenValidate(t *testing.T) {
	initApplication()
	createAllTable()

	invokeTest(_app.GetContext().GetTest("token/validate"))
}

func TestServiceUserInsert(t *testing.T) {
	initApplication()
	createAllTable()

	invokeTest(_app.GetContext().GetTest("user/insert"))
}

func TestServiceUserSelectPage(t *testing.T) {
	initApplication()
	createAllTable()

	invokeTest(_app.GetContext().GetTest("user/selectPage"))
}

func TestServiceUserBatchInsert(t *testing.T) {
	initApplication()
	createAllTable()

	invokeTest(_app.GetContext().GetTest("user/batchInsert"))
}
