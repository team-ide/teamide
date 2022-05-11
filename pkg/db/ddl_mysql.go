package db

import (
	"errors"
	"strings"
)

type DatabaseMySqlDialect struct {
	DatabaseDialect
}

func (this_ *DatabaseMySqlDialect) DatabaseDDL(param *GenerateParam, database *DatabaseModel) (sqlList []string, err error) {

	var sql string
	sql = `CREATE DATABASE ` + this_.packingCharacterDatabase(param, database.Name)
	if param.CharacterSet != "" {
		sql += ` CHARACTER SET ` + param.CharacterSet
	}
	if param.Collate != "" {
		sql += ` COLLATE '` + param.Collate + "'"
	}

	sqlList = append(sqlList, sql)

	return
}

func (this_ *DatabaseMySqlDialect) TableDDL(param *GenerateParam, database string, table *TableModel) (sqlList []string, err error) {
	sqlList = []string{}

	createTableSql := `CREATE TABLE `

	if param.AppendDatabase {
		createTableSql += this_.packingCharacterDatabase(param, database) + "."
	}
	createTableSql += this_.packingCharacterTable(param, table.Name)

	createTableSql += `(`
	createTableSql += "\n"
	primaryKeys := ""
	if len(table.ColumnList) > 0 {
		for _, one := range table.ColumnList {
			var columnSql = this_.packingCharacterColumn(param, one.Name)
			var c = DatabaseTypeMySql.GetColumnTypeInfo(one.Type)
			if c == nil {
				err = errors.New("MySql字段类型[" + one.Type + "]未映射!")
				return
			}
			columnType := c.FormatColumnType(one.Length, one.Decimal)

			columnSql += " " + columnType

			if param.CharacterSet != "" {
				columnSql += ` CHARACTER SET ` + param.CharacterSet
			}
			if one.Default != "" {
				columnSql += ` DEFAULT ` + formatStringValue("'", one.Default)
			}

			if one.NotNull {
				columnSql += ` NOT NULL`
			}

			if one.PrimaryKey {
				primaryKeys += "" + one.Name + ","
			}
			createTableSql += "\t" + columnSql
			createTableSql += ",\n"
		}
	}
	if primaryKeys != "" {
		primaryKeys = strings.TrimSuffix(primaryKeys, ",")
		createTableSql += "\tPRIMARY KEY (" + this_.packingCharacterColumns(param, primaryKeys) + ")"
	}

	createTableSql = strings.TrimSuffix(createTableSql, ",\n")
	createTableSql += "\n"

	createTableSql += `)`
	if param.CharacterSet != "" {
		createTableSql += ` DEFAULT CHARSET ` + param.CharacterSet
	}

	sqlList = append(sqlList, createTableSql)

	// 添加注释
	if table.Comment != "" {
		sqlList_ := this_.TableComment(param, database, table.Name, table.Comment)
		sqlList = append(sqlList, sqlList_...)
	}
	if len(table.ColumnList) > 0 {
		for _, one := range table.ColumnList {
			if one.Comment == "" {
				continue
			}
			var c = DatabaseTypeMySql.GetColumnTypeInfo(one.Type)
			if c == nil {
				err = errors.New("MySql字段类型[" + one.Type + "]未映射!")
				return
			}
			columnType := c.FormatColumnType(one.Length, one.Decimal)
			sqlList_ := this_.ColumnComment(param, database, table.Name, one.Name, columnType, one.Comment)
			sqlList = append(sqlList, sqlList_...)
		}
	}

	if len(table.IndexList) > 0 {
		for _, one := range table.IndexList {
			if one.Name == "" || one.Columns == "" {
				continue
			}
			//name := table.Name + "_" + one.Name
			name := one.Name
			sqlList_ := this_.Index(param, database, table.Name, name, one.Type, one.Columns, one.Comment)
			sqlList = append(sqlList, sqlList_...)

		}
	}
	return
}

func (this_ *DatabaseMySqlDialect) TableComment(param *GenerateParam, database string, table string, comment string) (sqlList []string) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += this_.packingCharacterDatabase(param, database) + "."
	}
	sql += "" + this_.packingCharacterTable(param, table)
	sql += " COMMENT " + formatStringValue("'", comment)
	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseMySqlDialect) ColumnComment(param *GenerateParam, database string, table string, column string, columnType string, comment string) (sqlList []string) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += this_.packingCharacterDatabase(param, database) + "."
	}
	sql += "" + this_.packingCharacterTable(param, table)
	sql += " MODIFY COLUMN " + this_.packingCharacterColumn(param, column)
	sql += " " + columnType
	sql += " COMMENT " + formatStringValue("'", comment)
	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseMySqlDialect) Index(param *GenerateParam, database string, table string, indexName string, indexType string, columns string, indexComment string) (sqlList []string) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += this_.packingCharacterDatabase(param, database) + "."
	}
	sql += "" + this_.packingCharacterTable(param, table)

	switch strings.ToUpper(indexType) {
	case "PRIMARY":
		sql += " ADD PRIMARY KEY "
	case "UNIQUE":
		sql += " ADD UNIQUE "
	case "FULLTEXT":
		sql += " ADD FULLTEXT "
	default:
		sql += " ADD INDEX "
	}
	if indexName != "" {
		sql += "" + indexName + " "
	}
	sql += "(" + this_.packingCharacterColumns(param, columns) + ")"

	if indexComment != "" {
		sql += " COMMENT " + formatStringValue("'", indexComment)
	}

	sqlList = append(sqlList, sql)
	return
}
