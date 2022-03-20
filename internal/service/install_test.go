package service

import (
	"testing"
)

func TestCheckMysql(t *testing.T) {
	var err error
	installService := NewInstallService(getMysqlDBWorker())
	err = installService.Check()
	if err != nil {
		panic(err)
	}
}

func TestCheckSqlite(t *testing.T) {
	var err error
	installService := NewInstallService(getSqliteDBWorker())
	err = installService.Check()
	if err != nil {
		panic(err)
	}
}

func TestInstallMysql(t *testing.T) {
	var err error
	installService := NewInstallService(getMysqlDBWorker())
	err = installService.Install()
	if err != nil {
		panic(err)
	}
}

func TestInstallSqlite(t *testing.T) {
	var err error
	installService := NewInstallService(getSqliteDBWorker())
	err = installService.Install()
	if err != nil {
		panic(err)
	}
}
