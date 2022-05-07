package sql_ddl

import (
	"errors"
	"fmt"
	"strings"
)

const (
	OracleCreateTable = `CREATE TABLE {table} (
[{columns}]
	[PRIMARY KEY {primaryKeys}]
)`
	OracleCreateTableColumn = `{column}[ {type}][ DEFAULT {default}][ NOT NULL]`
)

var (
	OracleTypeMap = map[string]func(length int, decimal int) string{}
)

func init() {
	OracleTypeMap["varchar"] = func(length int, decimal int) string {
		return fmt.Sprintf("varchar(%d)", length)
	}
	OracleTypeMap["text"] = func(length int, decimal int) string {
		return fmt.Sprintf("text")
	}
	OracleTypeMap["int"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	OracleTypeMap["bigint"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	OracleTypeMap["number"] = func(length int, decimal int) string {
		if decimal > 0 {
			return fmt.Sprintf("number(%d, %d)", length, decimal)
		}
		return fmt.Sprintf("number(%d)", length)
	}
	OracleTypeMap["datetime"] = func(length int, decimal int) string {
		return fmt.Sprintf("datetime")
	}
}

func ToDatabaseDDLForOracle(database string) (sqls []string, err error) {

	return
}

func ToTableDDLForOracle(table TableDetailInfo) (sqls []string, err error) {
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
			var typeFunc = OracleTypeMap[strings.ToLower(one.Type)]
			if typeFunc == nil {
				err = errors.New("字段类型[" + one.Type + "]未映射!")
				return
			}
			data["type"] = typeFunc(one.Length, one.Decimal)

			columnSql, err = formatSql(OracleCreateTableColumn, data)
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

	sql, err = formatSql(OracleCreateTable, data)
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
