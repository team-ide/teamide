package modelcoder

import (
	"fmt"
	"strings"
)

func getSqlInsertSqlParams(sqlInsert *DaoSqlInsertModel, application *Application, variable *invokeVariable) (sql string, sqlParams []interface{}, err error) {
	wrapTable := WrapTableName(sqlInsert.Database, sqlInsert.Table)

	sqlParams = []interface{}{}

	insertColumn := ""
	insertValue := ""

	for _, column := range sqlInsert.Columns {
		if IsEmpty(column.Name) {
			continue
		}
		var ifScript bool
		ifScript, err = ifScriptValue(application, variable, column.IfScript)
		if err != nil {
			return
		}
		if !ifScript {
			continue
		}
		if column.AutoIncrement {
			continue
		}
		var value interface{}
		value, err = getColumnValue(application, variable, column.Name, column.ValueScript)
		if err != nil {
			return
		}
		if !column.AllowEmpty && IsEmptyObj(value) {
			continue
		}

		wrapColumn := WrapColumnName("", column.Name)
		insertColumn += wrapColumn + ", "
		insertValue += "?, "
		sqlParams = append(sqlParams, value)

	}

	insertColumn = strings.TrimSuffix(insertColumn, ", ")
	insertValue = strings.TrimSuffix(insertValue, ", ")

	sql = fmt.Sprint("INSERT INTO ", wrapTable, "(", insertColumn, ")VALUES(", insertValue, ")")

	return
}

func getSqlInsertTableColumns(sqlInsert *DaoSqlInsertModel) (tableColumns map[string][]string) {
	tableColumns = map[string][]string{}

	wrapTable := WrapTableName(sqlInsert.Database, sqlInsert.Table)

	var columns []string
	for _, column := range sqlInsert.Columns {
		if IsEmpty(column.Name) {
			continue
		}
		wrapColumn := WrapColumnName("", column.Name)
		columns = append(columns, wrapColumn)

	}
	tableColumns[wrapTable] = columns

	return
}
