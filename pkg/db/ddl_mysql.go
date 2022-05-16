package db

import (
	"strings"
)

type DatabaseMySqlDialect struct {
	*BaseDialect
}

func (this_ *DatabaseMySqlDialect) DatabaseDDL(param *GenerateParam, database *DatabaseModel) (sqlList []string, err error) {

	var sql string
	sql = `CREATE DATABASE ` + param.packingCharacterDatabase(database.Name)
	if param.CharacterSet != "" {
		sql += ` CHARACTER SET ` + param.CharacterSet
	}
	if param.Collate != "" {
		sql += ` COLLATE '` + param.Collate + "'"
	}

	sqlList = append(sqlList, sql)

	return
}

func (this_ *DatabaseMySqlDialect) DatabaseDeleteDDL(param *GenerateParam, database string) (sqlList []string, err error) {

	var sql string
	sql = `DROP DATABASE ` + param.packingCharacterDatabase(database)

	sqlList = append(sqlList, sql)

	return
}

func (this_ *DatabaseMySqlDialect) TableDDL(param *GenerateParam, database string, table *TableModel) (sqlList []string, err error) {
	sqlList = []string{}

	createTableSql := `CREATE TABLE `

	if param.AppendDatabase {
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

			if param.CharacterSet != "" {
				columnSql += ` CHARACTER SET ` + param.CharacterSet
			}
			if column.NotNull {
				columnSql += ` NOT NULL`
			}
			if column.Default != nil {
				columnSql += " DEFAULT " + formatStringValue("'", GetStringValue(column.Default))
			}
			if column.Comment != "" {
				columnSql += " COMMENT " + formatStringValue("'", column.Comment)
			}

			if column.PrimaryKey {
				primaryKeys += "" + column.Name + ","
			}
			createTableSql += "\t" + columnSql
			createTableSql += ",\n"
		}
	}
	if primaryKeys != "" {
		primaryKeys = strings.TrimSuffix(primaryKeys, ",")
		createTableSql += "\tPRIMARY KEY (" + param.packingCharacterColumns(primaryKeys) + ")"
	}

	createTableSql = strings.TrimSuffix(createTableSql, ",\n")
	createTableSql += "\n"

	createTableSql += `)`
	if param.CharacterSet != "" {
		createTableSql += ` DEFAULT CHARSET ` + param.CharacterSet
	}

	sqlList = append(sqlList, createTableSql)

	var sqlList_ []string
	// 添加注释
	if table.Comment != "" {
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
			sqlList_, err = this_.TableIndexAddDDL(param, database, table.Name, one)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)

		}
	}
	return
}

func (this_ *DatabaseMySqlDialect) TableCommentDDL(param *GenerateParam, database string, table string, comment string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)
	sql += " COMMENT " + formatStringValue("'", comment)
	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseMySqlDialect) TableDeleteDDL(param *GenerateParam, database string, table string) (sqlList []string, err error) {

	var sql string
	sql = `DROP TABLE `

	if param.AppendDatabase {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += param.packingCharacterTable(table)
	sqlList = append(sqlList, sql)

	return
}
func (this_ *DatabaseMySqlDialect) TableColumnDeleteDDL(param *GenerateParam, database string, table string, column string) (sqlList []string, err error) {

	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += param.packingCharacterTable(table)

	sql += ` DROP COLUMN  `
	sql += param.packingCharacterColumn(column)

	sqlList = append(sqlList, sql)
	return
}
func (this_ *DatabaseMySqlDialect) TableColumnUpdateDDL(param *GenerateParam, database string, table string, column *TableColumnModel) (sqlList []string, err error) {
	var columnType string
	columnType, err = this_.GetColumnType(column.Type, column.Length, column.Decimal)
	if err != nil {
		return
	}

	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)
	if column.OldName != "" && column.Name != column.OldName {
		sql += " CHANGE COLUMN " + param.packingCharacterColumn(column.OldName)
	} else {
		sql += " MODIFY COLUMN"
	}
	sql += " " + param.packingCharacterColumn(column.Name)
	sql += " " + columnType
	if column.NotNull {
		sql += " NOT NULL"
	}
	if column.Default == nil {
		sql += " DEFAULT NULL"
	} else {
		sql += " DEFAULT " + formatStringValue("'", GetStringValue(column.Default))
	}
	sql += " COMMENT " + formatStringValue("'", column.Comment)
	if column.BeforeColumn != "" {
		sql += " AFTER " + param.packingCharacterColumn(column.BeforeColumn)
	}
	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseMySqlDialect) TableColumnAddDDL(param *GenerateParam, database string, table string, column *TableColumnModel) (sqlList []string, err error) {
	var columnType string
	columnType, err = this_.GetColumnType(column.Type, column.Length, column.Decimal)
	if err != nil {
		return
	}

	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)
	sql += " ADD COLUMN " + param.packingCharacterColumn(column.Name)
	sql += " " + columnType
	if column.NotNull {
		sql += " NOT NULL"
	}
	if column.Default == nil {
		sql += " DEFAULT NULL"
	} else {
		sql += " DEFAULT " + formatStringValue("'", GetStringValue(column.Default))
	}
	sql += " COMMENT " + formatStringValue("'", column.Comment)
	if column.BeforeColumn != "" {
		sql += " AFTER " + param.packingCharacterColumn(column.BeforeColumn)
	}
	sqlList = append(sqlList, sql)
	return
}
func (this_ *DatabaseMySqlDialect) TableIndexAddDDL(param *GenerateParam, database string, table string, index *TableIndexModel) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)

	switch strings.ToUpper(index.Type) {
	case "PRIMARY":
		sql += " ADD PRIMARY KEY "
	case "UNIQUE":
		sql += " ADD UNIQUE "
	case "FULLTEXT":
		sql += " ADD FULLTEXT "
	default:
		sql += " ADD INDEX "
	}
	if index.Name != "" {
		sql += "" + index.Name + " "
	}
	sql += "(" + param.packingCharacterColumns(strings.Join(index.Columns, ",")) + ")"

	if index.Comment != "" {
		sql += " COMMENT " + formatStringValue("'", index.Comment)
	}

	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseMySqlDialect) TableIndexUpdateDDL(param *GenerateParam, database string, table string, index *TableIndexModel) (sqlList []string, err error) {

	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)

	if index.OldName != "" {
		sql += " DROP INDEX " + param.packingCharacterColumn(index.OldName) + ","
	} else {
		sql += " DROP INDEX " + param.packingCharacterColumn(index.Name) + ","
	}
	switch strings.ToUpper(index.Type) {
	case "UNIQUE":
		sql += " ADD UNIQUE INDEX "
	default:
		sql += "ADD INDEX"
	}
	sql += " " + param.packingCharacterColumn(index.Name) + "(" + param.packingCharacterColumns(strings.Join(index.Columns, ",")) + ")"

	if index.Comment != "" {
		sql += " COMMENT " + formatStringValue("'", index.Comment)
	}
	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseMySqlDialect) TableIndexRenameDDL(param *GenerateParam, database string, table string, index string, rename string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)

	sql += ` RENAME INDEX `
	sql += "" + param.packingCharacterColumn(index)
	sql += ` TO `
	sql += "" + param.packingCharacterColumn(rename)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseMySqlDialect) TableIndexDeleteDDL(param *GenerateParam, database string, table string, index string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)

	sql += ` DROP INDEX `
	sql += "" + param.packingCharacterColumn(index)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseMySqlDialect) TablePrimaryKeyDeleteDDL(param *GenerateParam, database string, table string, primaryKeys []string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)

	sql += ` DROP PRIMARY KEY `

	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseMySqlDialect) TablePrimaryKeyAddDDL(param *GenerateParam, database string, table string, primaryKeys []string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)

	sql += ` ADD PRIMARY KEY `

	sql += "(" + param.packingCharacterColumns(strings.Join(primaryKeys, ",")) + ")"

	sqlList = append(sqlList, sql)
	return
}
