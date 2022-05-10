package sql_ddl

import (
	"errors"
	"fmt"
	"strings"
	"teamide/pkg/db"
)

const (
	MySqlCreateDatabase = `CREATE DATABASE {database}[ CHARACTER SET '{characterSet}'][ COLLATE '{collate}']`
	MySqlCreateTable    = `CREATE TABLE {table} (
[{columns}]
	[PRIMARY KEY {primaryKeys}]
[{indexs}]
)[ ENGINE={engine}][ DEFAULT CHARSET={defaultCharset}][ COMMENT='{comment}']`
	MySqlCreateTableColumn      = `{column}[ {type}][ CHARACTER SET {characterSet}][ NOT NULL][ DEFAULT '{default}'][ AUTO_INCREMENT][ COMMENT '{comment}']`
	MySqlCreateTableIndex       = `KEY {name} ({columns})[ COMMENT '{comment}']`
	MySqlCreateTableIndexUnique = `UNIQUE KEY {name} ({columns})[ COMMENT '{comment}']`
)

var (
	MysqlTypeMap = map[string]func(length int, decimal int) string{}
)

func init() {

	MysqlTypeMap["varchar"] = func(length int, decimal int) string {
		return fmt.Sprintf("varchar(%d)", length)
	}
	MysqlTypeMap["text"] = func(length int, decimal int) string {
		return fmt.Sprintf("text")
	}
	MysqlTypeMap["int"] = func(length int, decimal int) string {
		return fmt.Sprintf("number(%d)", length)
	}
	MysqlTypeMap["bigint"] = func(length int, decimal int) string {
		return fmt.Sprintf("number(%d)", length)
	}
	MysqlTypeMap["number"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	MysqlTypeMap["datetime"] = func(length int, decimal int) string {
		return fmt.Sprintf("datetime")
	}

}

func ToDatabaseDDLForMySql(database string) (sqls []string, err error) {

	var sql string
	var data map[string]string = map[string]string{}
	data["database"] = fmt.Sprint("", database, "")
	//data["characterSet"] = database.CharacterSet
	//data["collate"] = database.Collate
	sql, err = formatSql(MySqlCreateDatabase, data)

	sqls = append(sqls, sql)

	return
}

func ToTableDDLForMySql(table *TableDetailInfo) (sqls []string, err error) {
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
				err = errors.New("MySql字段类型[" + one.Type + "]未映射!")
				return
			}
			data["type"] = c.FormatColumnType(one.Length, one.Decimal)
			columnSql, err = formatSql(MySqlCreateTableColumn, data)
			if err != nil {
				return
			}
			if columnSql == "" {
				continue
			}
			columns += "\t" + columnSql + ",\n"
		}
	}

	if len(table.Indexs) > 0 {
		var indexSql string
		for _, one := range table.Indexs {
			data = map[string]string{}
			if one.Name == "" || one.Columns == "" {
				continue
			}
			data["name"] = fmt.Sprint("", one.Name, "")
			data["columns"] = fmt.Sprint("", one.Columns, "")
			data["comment"] = fmt.Sprint("", one.Comment, "")

			switch one.Type {
			case "UNIQUE", "unique":
				indexSql, err = formatSql(MySqlCreateTableIndexUnique, data)
			default:
				indexSql, err = formatSql(MySqlCreateTableIndex, data)
			}

			if err != nil {
				return
			}
			if indexSql == "" {
				continue
			}
			indexs += "\t" + indexSql + ",\n"
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
	sql, err = formatSql(MySqlCreateTable, data)
	if err != nil {
		return
	}
	if sql != "" {
		sqls = append(sqls, sql)
	}
	return
}
