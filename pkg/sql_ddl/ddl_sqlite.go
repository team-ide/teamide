package sql_ddl

import (
	"errors"
	"fmt"
	"strings"
	"teamide/pkg/db"
)

const (
	SqliteCreateTable = `CREATE TABLE {table} (
[{columns}]
	[PRIMARY KEY {primaryKeys}]
)`
	SqliteCreateTableColumn = `{column}[ {type}][ DEFAULT {default}][ NOT NULL]`
)

var (
	SqliteTypeMap = map[string]func(length int, decimal int) string{}
)

func init() {
	SqliteTypeMap["varchar"] = func(length int, decimal int) string {
		return fmt.Sprintf("varchar(%d)", length)
	}
	SqliteTypeMap["text"] = func(length int, decimal int) string {
		return fmt.Sprintf("text")
	}
	SqliteTypeMap["int"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	SqliteTypeMap["bigint"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	SqliteTypeMap["number"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	SqliteTypeMap["datetime"] = func(length int, decimal int) string {
		return fmt.Sprintf("datetime")
	}
}

func ToDatabaseDDLForSqlite(database string) (sqls []string, err error) {

	return
}

func ToTableDDLForSqlite(table *TableDetailInfo) (sqls []string, err error) {
	sqls = []string{}
	var columns string
	var primaryKeys string
	var indexs string
	var data map[string]string

	if len(table.Columns) > 0 {
		var columnSql string
		for _, one := range table.Columns {
			data = map[string]string{}
			if one.Name == "" {
				continue
			}
			if one.PrimaryKey {
				primaryKeys += "" + one.Name + ","
			}
			data["column"] = fmt.Sprint("", one.Name, "")
			data["comment"] = fmt.Sprint("", one.Comment, "")
			data["default"] = fmt.Sprint("", one.Default, "")
			if one.NotNull {
				data["NOT NULL"] = "true"
			}
			var c = db.DatabaseTypeMySql.GetColumnTypeInfo(one.Type)
			if c == nil {
				err = errors.New("Sqlite字段类型[" + one.Type + "]未映射!")
				return
			}
			data["type"] = c.FormatColumnType(one.Length, one.Decimal)

			columnSql, err = formatSql(SqliteCreateTableColumn, data)
			if err != nil {
				return
			}
			if columnSql == "" {
				continue
			}
			columns += "\t" + columnSql + ",\n"
		}
	}

	data = map[string]string{}
	data["table"] = fmt.Sprint("", table.Name, "")

	columns = strings.TrimSuffix(columns, "\n")
	if primaryKeys == "" && indexs == "" {
		columns = strings.TrimSuffix(columns, ",")
	}
	primaryKeys = strings.TrimSuffix(primaryKeys, ",")
	if primaryKeys != "" {
		primaryKeys = "(" + primaryKeys + ")"
	}
	if indexs != "" {
		primaryKeys += ","
	}
	indexs = strings.TrimSuffix(indexs, "\n")
	indexs = strings.TrimSuffix(indexs, ",")
	data["columns"] = columns
	data["primaryKeys"] = primaryKeys
	data["indexs"] = indexs
	data["comment"] = table.Comment
	var sql string

	sql, err = formatSql(SqliteCreateTable, data)
	if err != nil {
		return
	}
	if sql != "" {
		sqls = append(sqls, sql)
	}
	// 添加注释
	if table.Comment != "" {
		sqls = append(sqls, `COMMENT ON TABLE `+table.Name+` IS '`+table.Comment+`'`)
	}
	if len(table.Columns) > 0 {
		for _, one := range table.Columns {
			if one.Name == "" || one.Comment == "" {
				continue
			}
			sqls = append(sqls, `COMMENT ON COLUMN `+table.Name+`.`+one.Name+` IS '`+one.Comment+`'`)
		}
	}

	if len(table.Indexs) > 0 {
		for _, one := range table.Indexs {
			data = map[string]string{}
			if one.Name == "" || one.Columns == "" {
				continue
			}

			name := table.Name + "_" + one.Name

			switch one.Type {
			case "UNIQUE", "unique":
				sqls = append(sqls, `CREATE UNIQUE INDEX `+name+` ON `+table.Name+`(`+one.Columns+`)`)
			default:
				sqls = append(sqls, `CREATE INDEX `+name+` ON `+table.Name+`(`+one.Columns+`)`)
			}

		}
	}
	return
}
