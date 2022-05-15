package db

import "strings"

type DatabaseSqliteDialect struct {
	*BaseDialect
}

func (this_ *DatabaseSqliteDialect) TableDDL(param *GenerateParam, database string, table *TableModel) (sqlList []string, err error) {
	sqlList = []string{}

	createTableSql := `CREATE TABLE `

	if param.AppendDatabase && database != "" {
		createTableSql += param.packingCharacterDatabase(database) + "."
	}
	createTableSql += param.packingCharacterTable(table.Name)

	createTableSql += `(`
	createTableSql += "\n"
	primaryKeys := ""
	if len(table.ColumnList) > 0 {
		for _, column := range table.ColumnList {
			var columnSql = param.packingCharacterColumn(column.Name)

			var columnType string
			columnType, err = this_.GetColumnType(column.Type, column.Length, column.Decimal)
			if err != nil {
				return
			}
			columnSql += " " + columnType

			if column.NotNull {
				columnSql += ` NOT NULL`
			}
			if column.Default != nil {
				columnSql += ` DEFAULT ` + formatStringValue("'", GetStringValue(column.Default))
			}

			if column.PrimaryKey {
				primaryKeys += "" + column.Name + ","
			}
			createTableSql += "\t" + columnSql + ",\n"
		}
	}
	if primaryKeys != "" {
		primaryKeys = strings.TrimSuffix(primaryKeys, ",")
		createTableSql += "\tPRIMARY KEY (" + param.packingCharacterColumns(primaryKeys) + ")"
	}

	createTableSql = strings.TrimSuffix(createTableSql, ",\n")
	createTableSql += "\n"

	createTableSql += `)`

	sqlList = append(sqlList, createTableSql)

	// 添加注释
	if table.Comment != "" {
		var sqlList_ []string
		sqlList_, err = this_.TableCommentDDL(param, database, table.Name, table.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}
	if len(table.IndexList) > 0 {
		for _, one := range table.IndexList {
			if one.Name == "" || len(one.Columns) == 0 {
				continue
			}
			//name := table.Name + "_" + one.Name
			name := one.Name
			if !strings.HasPrefix(name, table.Name+"_INDEX_") {
				name = table.Name + "_INDEX_" + one.Name
			}
			one.Name = name
			var sqlList_ []string
			sqlList_, err = this_.TableIndexAddDDL(param, database, table.Name, one)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
	}
	return
}
