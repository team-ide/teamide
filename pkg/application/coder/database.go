package coder

import (
	"fmt"
	"regexp"
	"strings"
	common2 "teamide/pkg/application/common"
	model2 "teamide/pkg/application/model"
)

const (
	CREATE_DATABASE = `CREATE DATABASE IF NOT EXISTS {database}[ CHARACTER SET '{characterSet}'][ COLLATE '{collate}']`
	CREATE_TABLE    = `CREATE TABLE IF NOT EXISTS {table} (
  [{columns}]
  [PRIMARY KEY {primaryKeys}]
  [{indexs}]
)[ ENGINE={engine}][ DEFAULT CHARSET={defaultCharset}][ COMMENT='{comment}']`
	CREATE_TABLE_COLUMN       = `{column}[ {type}][ CHARACTER SET {characterSet}][ NOT NULL][ DEFAULT '{default}'][ AUTO_INCREMENT][ COMMENT '{comment}']`
	CREATE_TABLE_INDEX        = `KEY {name} ({columns})[ COMMENT '{comment}']`
	CREATE_TABLE_INDEX_UNIQUE = `UNIQUE KEY {name} ({columns})[ COMMENT '{comment}']`

	ORACLE_CREATE_TABLE = `CREATE TABLE {table} (
		[{columns}]
		[PRIMARY KEY {primaryKeys}]
	  )`
	ORACLE_CREATE_TABLE_COLUMN = `{column}[ {type}][ DEFAULT {default}][ NOT NULL]`
)

func GetCreateTableSqls(app common2.IApplication, database *model2.DatasourceDatabase) (sqls []string, err error) {
	var sqls_ []string
	// sql_, err = GetDatabaseDDL(*database)
	// if err != nil {
	// 	return
	// }
	// if sql_ == "" {
	// 	return
	// }
	// sqls = append(sqls, sql_)
	// sqls = append(sqls, "USE "+database.Database+"")
	if len(app.GetContext().Structs) > 0 {
		for _, one := range app.GetContext().Structs {
			database_ := app.GetContext().GetDatasourceDatabase(one.Database)
			if database_ != database {
				continue
			}
			sqls_, err = GetTableDDL(database_, one)
			if err != nil {
				return
			}
			if len(sqls_) == 0 {
				continue
			}
			sqls = append(sqls, sqls_...)
		}
	}
	return
}

func GetDatabaseDDL(database *model2.DatasourceDatabase) (sql string, err error) {
	if database.Database == "" {
		return
	}
	var data map[string]string = map[string]string{}
	data["database"] = fmt.Sprint("", database.Database, "")
	data["characterSet"] = database.CharacterSet
	data["collate"] = database.Collate
	sql, err = foramtSql(CREATE_DATABASE, data)
	return
}
func GetTableDDL(database *model2.DatasourceDatabase, structModel *model2.StructModel) (sqls []string, err error) {
	sqls = []string{}
	if structModel.Table == "" {
		return
	}
	var columns string
	var primaryKeys string
	var indexs string
	var data map[string]string

	if len(structModel.Fields) > 0 {
		var columnSql string
		for _, one := range structModel.Fields {
			data = map[string]string{}
			if one.Column == "" {
				continue
			}
			if one.PrimaryKey {
				primaryKeys += "" + one.Column + ","
			}
			data["column"] = fmt.Sprint("", one.Column, "")
			data["comment"] = fmt.Sprint("", one.Comment, "")
			data["default"] = fmt.Sprint("", one.Default, "")
			if one.NotNull {
				data["NOT NULL"] = "true"
			}
			dataType := one.DataType
			var typeStr string
			typeStr = one.ColumnType
			if common2.DatabaseIsMySql(database) {
				switch dataType {
				case "long", "int64":
					typeStr = "bigint"
				case "int", "int32":
					typeStr = "int"
				case "short", "int16":
					typeStr = "int"
				case "byte", "int8":
					typeStr = "int"
				case "date", "datetime", "time", "time.Time", "Time":
					typeStr = "datetime"
				case "boolean", "bool":
					typeStr = "int"
					one.ColumnLength = 1
				case "float", "float64":
					typeStr = "number"
				case "double", "float32":
					typeStr = "number"
				default:
					typeStr = "varchar"
				}
				if one.ColumnLength > 0 {
					if one.ColumnDecimal > 0 {
						typeStr = fmt.Sprint("", typeStr, "(", one.ColumnLength, ",", one.ColumnDecimal, ")")
					} else {
						typeStr = fmt.Sprint("", typeStr, "(", one.ColumnLength, ")")
					}
				} else {
					typeStr = fmt.Sprint("", typeStr, "")
				}
			} else if common2.DatabaseIsOracle(database) {
				switch dataType {
				case "long", "int64":
					typeStr = "number"
				case "int", "int32":
					typeStr = "number"
				case "short", "int16":
					typeStr = "number"
				case "byte", "int8":
					typeStr = "number"
				case "date", "datetime", "time", "time.Time", "Time":
					typeStr = "date"
				case "boolean", "bool":
					typeStr = "number"
					one.ColumnLength = 1
				case "float", "float64":
					typeStr = "number"
				case "double", "float32":
					typeStr = "number"
				default:
					typeStr = "varchar2"
				}
				if one.ColumnLength > 0 {
					if one.ColumnDecimal > 0 {
						typeStr = fmt.Sprint("", typeStr, "(", one.ColumnLength, ",", one.ColumnDecimal, ")")
					} else {
						typeStr = fmt.Sprint("", typeStr, "(", one.ColumnLength, ")")
					}
				} else {
					typeStr = fmt.Sprint("", typeStr, "")
				}
			}
			data["type"] = typeStr
			if common2.DatabaseIsMySql(database) {
				columnSql, err = foramtSql(CREATE_TABLE_COLUMN, data)
			} else if common2.DatabaseIsOracle(database) {
				columnSql, err = foramtSql(ORACLE_CREATE_TABLE_COLUMN, data)
			}
			if err != nil {
				return
			}
			if columnSql == "" {
				continue
			}
			if columns != "" {
				columns += "  "
			}
			columns += columnSql + ",\n"
		}
	}

	if common2.DatabaseIsMySql(database) && len(structModel.Indexs) > 0 {
		var indexSql string
		for _, one := range structModel.Indexs {
			data = map[string]string{}
			if one.Name == "" || one.Columns == "" {
				continue
			}
			data["name"] = fmt.Sprint("", one.Name, "")
			data["columns"] = fmt.Sprint("", one.Columns, "")
			data["comment"] = fmt.Sprint("", one.Comment, "")

			switch one.Type {
			case "UNIQUE", "unique":
				indexSql, err = foramtSql(CREATE_TABLE_INDEX_UNIQUE, data)
			default:
				indexSql, err = foramtSql(CREATE_TABLE_INDEX, data)
			}

			if err != nil {
				return
			}
			if indexSql == "" {
				continue
			}
			if indexs != "" {
				indexs += "  "
			}
			indexs += indexSql + ",\n"
		}
	}
	data = map[string]string{}
	data["table"] = fmt.Sprint("", structModel.Table, "")

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
	data["comment"] = structModel.Comment
	var sql string
	if common2.DatabaseIsMySql(database) {
		sql, err = foramtSql(CREATE_TABLE, data)
		if err != nil {
			return
		}
		if sql != "" {
			sqls = append(sqls, sql)
		}
	} else if common2.DatabaseIsOracle(database) {
		sql, err = foramtSql(ORACLE_CREATE_TABLE, data)
		if err != nil {
			return
		}
		if sql != "" {
			sqls = append(sqls, sql)
		}
		// 添加注释
		if structModel.Comment != "" {
			sqls = append(sqls, `COMMENT ON TABLE "`+structModel.Table+`" IS '`+structModel.Comment+`'`)
		}
		if len(structModel.Fields) > 0 {
			for _, one := range structModel.Fields {
				if one.Column == "" || one.Comment == "" {
					continue
				}
				sqls = append(sqls, `COMMENT ON COLUMN `+structModel.Table+`.`+one.Column+` IS '`+one.Comment+`'`)
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
