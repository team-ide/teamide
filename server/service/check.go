package service

import (
	"fmt"
	"server/base"
	"server/component"
	"server/factory"
	sqlModel "server/model/sql"
	"strings"
)

func CheckModel() {

	for _, one := range factory.Models {
		checkSelectSqlModels(one.GetSelectSqls())
		checkInsertSqlModels(one.GetInsertSqls())
		checkUpdateSqlModels(one.GetUpdateSqls())
		checkDeleteSqlModels(one.GetDeleteSqls())
	}
}

func checkSelectSqlModels(models []*sqlModel.Select) {

	for _, one := range models {
		checkSelectSqlModel(one)
	}
}

func checkInsertSqlModels(models []*sqlModel.Insert) {

	for _, one := range models {
		checkInsertSqlModel(one)
	}
}

func checkUpdateSqlModels(models []*sqlModel.Update) {

	for _, one := range models {
		checkUpdateSqlModel(one)
	}
}

func checkDeleteSqlModels(models []*sqlModel.Delete) {

	for _, one := range models {
		checkDeleteSqlModel(one)
	}
}

func checkSelectSqlModel(model *sqlModel.Select) {
	tableColumns := model.GetTableColumns()

	for table, columns := range tableColumns {
		err := checkTableColumns(table, columns)
		if err != nil {
			fmt.Println("SqlModel表字段验证失败:", base.ToJSON(model))
			panic(err)
		}
	}
}

func checkInsertSqlModel(model *sqlModel.Insert) {
	tableColumns := model.GetTableColumns()

	for table, columns := range tableColumns {
		err := checkTableColumns(table, columns)
		if err != nil {
			fmt.Println("SqlModel表字段验证失败:", base.ToJSON(model))
			panic(err)
		}
	}
}

func checkUpdateSqlModel(model *sqlModel.Update) {
	tableColumns := model.GetTableColumns()

	for table, columns := range tableColumns {
		err := checkTableColumns(table, columns)
		if err != nil {
			fmt.Println("SqlModel表字段验证失败:", base.ToJSON(model))
			panic(err)
		}
	}
}

func checkDeleteSqlModel(model *sqlModel.Delete) {
	tableColumns := model.GetTableColumns()

	for table, columns := range tableColumns {
		err := checkTableColumns(table, columns)
		if err != nil {
			fmt.Println("SqlModel表字段验证失败:", base.ToJSON(model))
			panic(err)
		}
	}
}

func checkTableColumns(table string, columns []string) (err error) {

	sql := "SELECT "
	for _, column := range columns {
		sql += column + ","
	}
	sql = strings.TrimSuffix(sql, ",")
	sql += " FROM " + table
	_, err = component.DB.Exec(base.SqlParam{Sql: sql})
	if err != nil {
		return
	}
	return
}
