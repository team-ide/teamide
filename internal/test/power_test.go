package test

import (
	"fmt"
	"teamide/internal/module/module_power"
	"testing"
)

func TestPowerMysql(t *testing.T) {
	service := module_power.NewPowerRoleService(getMysqlServerContext())

	testPower(service)
}

func TestPowerSqlite(t *testing.T) {
	service := module_power.NewPowerRoleService(getSqliteServerContext())

	testPower(service)
}

func testPower(service *module_power.PowerRoleService) {
	var err error

	powerRole := &module_power.PowerRoleModel{
		Name: "超管",
	}

	_, err = service.Insert(powerRole)
	if err != nil {
		panic(err)
	}
	fmt.Printf(fmt.Sprint(powerRole))
}
