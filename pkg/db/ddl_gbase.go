package db

import (
	"errors"
	"strings"
)

type DatabaseGBaseDialect struct {
	DatabaseDialect
}

func (this_ *DatabaseGBaseDialect) DatabaseDDL(param *GenerateParam, database *DatabaseModel) (sqlList []string, err error) {

	return
}

func (this_ *DatabaseGBaseDialect) TableDDL(param *GenerateParam, database string, table *TableModel) (sqlList []string, err error) {
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
			var c = DatabaseTypeGBase.GetColumnTypeInfo(one.Type)
			if c == nil {
				err = errors.New("GBase字段类型[" + one.Type + "]未映射!")
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
			sqlList_ := this_.ColumnComment(param, database, table.Name, one.Name, one.Comment)
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
			if !strings.HasPrefix(name, table.Name+"_INDEX_") {
				name = table.Name + "_INDEX_" + one.Name
			}
			sqlList_ := this_.Index(param, database, table.Name, name, one.Type, one.Columns)
			sqlList = append(sqlList, sqlList_...)

		}
	}
	return
}

func (this_ *DatabaseGBaseDialect) TableComment(param *GenerateParam, database string, table string, comment string) (sqlList []string) {
	sql := "COMMENT ON TABLE  "
	if param.AppendDatabase && database != "" {
		sql += this_.packingCharacterDatabase(param, database) + "."
	}
	sql += "" + this_.packingCharacterTable(param, table)
	sql += " IS " + formatStringValue("'", comment)
	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseGBaseDialect) ColumnComment(param *GenerateParam, database string, table string, column string, comment string) (sqlList []string) {
	sql := "COMMENT ON COLUMN "
	if param.AppendDatabase && database != "" {
		sql += this_.packingCharacterDatabase(param, database) + "."
	}
	sql += "" + this_.packingCharacterTable(param, table)
	sql += "." + this_.packingCharacterColumn(param, column)
	sql += " IS " + formatStringValue("'", comment)
	sqlList = append(sqlList, sql)
	return
}

func (this_ *DatabaseGBaseDialect) Index(param *GenerateParam, database string, table string, indexName string, indexType string, columns string) (sqlList []string) {

	sql := ""
	switch strings.ToUpper(indexType) {
	case "UNIQUE":
		sql += "CREATE UNIQUE INDEX "
	default:
		sql += "CREATE INDEX "
	}
	if indexName != "" {
		sql += "" + indexName + " "
	}
	sql += " ON "
	if param.AppendDatabase && database != "" {
		sql += this_.packingCharacterDatabase(param, database) + "."
	}
	sql += "" + this_.packingCharacterTable(param, table)

	sql += "(" + this_.packingCharacterColumns(param, columns) + ")"

	sqlList = append(sqlList, sql)
	return
}
