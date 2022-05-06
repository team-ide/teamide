package toolbox

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	CreateDatabase = `CREATE DATABASE {database}[ CHARACTER SET '{characterSet}'][ COLLATE '{collate}']`
	CreateTable    = `CREATE TABLE {table} (
[{columns}]
	[PRIMARY KEY {primaryKeys}]
[{indexs}]
)[ ENGINE={engine}][ DEFAULT CHARSET={defaultCharset}][ COMMENT='{comment}']`
	CreateTableColumn      = `{column}[ {type}][ CHARACTER SET {characterSet}][ NOT NULL][ DEFAULT '{default}'][ AUTO_INCREMENT][ COMMENT '{comment}']`
	CreateTableIndex       = `KEY {name} ({columns})[ COMMENT '{comment}']`
	CreateTableIndexUnique = `UNIQUE KEY {name} ({columns})[ COMMENT '{comment}']`

	OracleCreateTable = `CREATE TABLE {table} (
[{columns}]
	[PRIMARY KEY {primaryKeys}]
)`
	OracleCreateTableColumn = `{column}[ {type}][ DEFAULT {default}][ NOT NULL]`
)

var (
	mysqlTypeMap  = map[string]func(length int, decimal int) string{}
	oracleTypeMap = map[string]func(length int, decimal int) string{}
)

func init() {

	mysqlTypeMap["varchar"] = func(length int, decimal int) string {
		return fmt.Sprintf("varchar(%d)", length)
	}
	mysqlTypeMap["text"] = func(length int, decimal int) string {
		return fmt.Sprintf("text")
	}
	mysqlTypeMap["int"] = func(length int, decimal int) string {
		return fmt.Sprintf("number(%d)", length)
	}
	mysqlTypeMap["bigint"] = func(length int, decimal int) string {
		return fmt.Sprintf("number(%d)", length)
	}
	mysqlTypeMap["number"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	mysqlTypeMap["datetime"] = func(length int, decimal int) string {
		return fmt.Sprintf("datetime")
	}

	oracleTypeMap["varchar"] = func(length int, decimal int) string {
		return fmt.Sprintf("varchar(%d)", length)
	}
	oracleTypeMap["text"] = func(length int, decimal int) string {
		return fmt.Sprintf("text")
	}
	oracleTypeMap["int"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	oracleTypeMap["bigint"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	oracleTypeMap["number"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	oracleTypeMap["datetime"] = func(length int, decimal int) string {
		return fmt.Sprintf("datetime")
	}
}

func ToDatabaseDDL(database string, databaseType string) (sqls []string, err error) {

	if DatabaseIsMySql(databaseType) {
		var sql string
		var data map[string]string = map[string]string{}
		data["database"] = fmt.Sprint("", database, "")
		//data["characterSet"] = database.CharacterSet
		//data["collate"] = database.Collate
		sql, err = foramtSql(CreateDatabase, data)

		sqls = append(sqls, sql)
	}

	return
}

func ToTableDDL(databaseType string, table TableDetailInfo) (sqls []string, err error) {
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
			var typeFunc func(length int, decimal int) string
			if DatabaseIsMySql(databaseType) {
				typeFunc = mysqlTypeMap[strings.ToLower(one.Type)]
			} else if DatabaseIsOracle(databaseType) {
				typeFunc = oracleTypeMap[strings.ToLower(one.Type)]
			}
			if typeFunc == nil {
				err = errors.New("字段类型[" + one.Type + "]未映射!")
				return
			}
			data["type"] = typeFunc(one.Length, one.Decimal)
			if DatabaseIsMySql(databaseType) {
				columnSql, err = foramtSql(CreateTableColumn, data)
			} else if DatabaseIsOracle(databaseType) {
				columnSql, err = foramtSql(OracleCreateTableColumn, data)
			}
			if err != nil {
				return
			}
			if columnSql == "" {
				continue
			}
			columns += "\t" + columnSql + ",\n"
		}
	}

	if DatabaseIsMySql(databaseType) && len(table.Indexs) > 0 {
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
				indexSql, err = foramtSql(CreateTableIndexUnique, data)
			default:
				indexSql, err = foramtSql(CreateTableIndex, data)
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
	if DatabaseIsMySql(databaseType) {
		sql, err = foramtSql(CreateTable, data)
		if err != nil {
			return
		}
		if sql != "" {
			sqls = append(sqls, sql)
		}
	} else if DatabaseIsOracle(databaseType) {
		sql, err = foramtSql(OracleCreateTable, data)
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
	}
	return
}

func foramtSql(sql string, data map[string]string) (foramtSql string, err error) {
	var re *regexp.Regexp
	re, err = regexp.Compile(`\[(.+?)\]`)
	if err != nil {
		return
	}
	indexsList := re.FindAllIndex([]byte(sql), -1)
	var lastIndex int = 0
	var sql_ string
	var formatValueSql string
	var find bool = true
	for _, indexs := range indexsList {
		sql_ = sql[lastIndex:indexs[0]]
		formatValueSql, find = foramtValueSql(sql_, data)
		if find {
			foramtSql += formatValueSql
		}

		lastIndex = indexs[1]

		sql_ = sql[indexs[0]+1 : indexs[1]-1]

		if !strings.Contains(sql_, `{`) {
			if data[strings.TrimSpace(sql_)] != "" {
				foramtSql += sql_
			}
		} else {
			formatValueSql, find = foramtValueSql(sql_, data)
			if find {
				foramtSql += formatValueSql
			}
		}
	}
	sql_ = sql[lastIndex:]
	formatValueSql, find = foramtValueSql(sql_, data)
	if find {
		foramtSql += formatValueSql
	}
	return
}

func foramtValueSql(sql string, data map[string]string) (res string, find bool) {
	var re *regexp.Regexp
	re, _ = regexp.Compile(`{(.+?)}`)
	find = true
	indexsList := re.FindAllIndex([]byte(sql), -1)
	var lastIndex int = 0
	for _, indexs := range indexsList {
		res += sql[lastIndex:indexs[0]]

		lastIndex = indexs[1]

		key := sql[indexs[0]+1 : indexs[1]-1]
		value := data[key]
		if value == "" {
			find = false
			return
		}
		res += value
	}
	res += sql[lastIndex:]
	return
}
