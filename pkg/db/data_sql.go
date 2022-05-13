package db

import (
	"errors"
	"strings"
)

func DataListInsertSql(param *GenerateParam, database string, table string, columnList []*TableColumnModel, dataList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error) {
	if len(dataList) == 0 {
		return
	}
	var keys []string
	for _, column := range columnList {
		if column.PrimaryKey {
			keys = append(keys, column.Name)
		}
	}

	for _, data := range dataList {

		var values []interface{}
		insertColumns := ""
		insertValues := ""
		for _, column := range columnList {
			value, valueOk := data[column.Name]
			if !valueOk {
				continue
			}
			insertColumns += param.packingCharacterColumn(column.Name) + ", "
			if param.AppendSqlValue {
				insertValues += param.packingCharacterColumnStringValue(column, value) + ", "
			} else {
				insertValues += "?, "
				values = append(values, param.formatColumnValue(column, value))
			}
		}
		insertColumns = strings.TrimSuffix(insertColumns, ", ")
		insertValues = strings.TrimSuffix(insertValues, ", ")

		sql := "INSERT INTO "

		if param.AppendDatabase && database != "" {
			sql += param.packingCharacterDatabase(database) + "."
		}
		sql += param.packingCharacterTable(table)
		if insertColumns != "" {
			sql += "(" + insertColumns + ")"
		}
		if insertValues != "" {
			sql += " VALUES (" + insertValues + ")"
		}

		sqlList = append(sqlList, sql)
		valuesList = append(valuesList, values)
	}
	return
}

func DataListUpdateSql(param *GenerateParam, database string, table string, columnList []*TableColumnModel, dataList []map[string]interface{}, dataWhereList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error) {
	if len(dataList) == 0 {
		return
	}
	if len(dataList) != len(dataWhereList) {
		err = errors.New("更新数据与更新条件数量不一致")
		return
	}
	var keyColumnList []*TableColumnModel
	for _, column := range columnList {
		if column.PrimaryKey {
			keyColumnList = append(keyColumnList, column)
		}
	}

	for index, data := range dataList {
		dataWhere := dataWhereList[index]
		if len(dataWhere) == 0 {
			err = errors.New("更新数据条件丢失")
			return
		}

		sql := "UPDATE "
		var values []interface{}

		if param.AppendDatabase && database != "" {
			sql += param.packingCharacterDatabase(database) + "."
		}
		sql += param.packingCharacterTable(table)
		sql += " SET "

		for _, column := range columnList {
			value, valueOK := data[column.Name]
			if !valueOK {
				continue
			}
			if param.AppendSqlValue {
				sql += "" + param.packingCharacterColumn(column.Name) + "=" + param.packingCharacterColumnStringValue(column, value) + ", "
			} else {
				sql += "" + param.packingCharacterColumn(column.Name) + "=?, "
				values = append(values, param.formatColumnValue(column, value))
			}
		}
		sql = strings.TrimSuffix(sql, ", ")

		sql += " WHERE "
		whereColumnList := keyColumnList
		if len(keyColumnList) == 0 {
			whereColumnList = columnList
		} else {
			for _, column := range whereColumnList {
				if param.AppendSqlValue {
					sql += "" + param.packingCharacterColumn(column.Name) + "=" + param.packingCharacterColumnStringValue(column, dataWhere[column.Name]) + " AND "
				} else {
					sql += "" + param.packingCharacterColumn(column.Name) + "=? AND "
					values = append(values, param.formatColumnValue(column, dataWhere[column.Name]))
				}
			}
		}
		sql = strings.TrimSuffix(sql, " AND ")

		sqlList = append(sqlList, sql)
		valuesList = append(valuesList, values)
	}
	return
}

func DataListDeleteSql(param *GenerateParam, database string, table string, columnList []*TableColumnModel, dataWhereList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error) {
	if len(dataWhereList) == 0 {
		return
	}
	var keyColumnList []*TableColumnModel
	for _, column := range columnList {
		if column.PrimaryKey {
			keyColumnList = append(keyColumnList, column)
		}
	}

	for _, dataWhere := range dataWhereList {
		if len(dataWhere) == 0 {
			err = errors.New("更新数据条件丢失")
			return
		}

		sql := "DELETE FROM "
		var values []interface{}

		if param.AppendDatabase && database != "" {
			sql += param.packingCharacterDatabase(database) + "."
		}
		sql += param.packingCharacterTable(table)

		sql += " WHERE "
		whereColumnList := keyColumnList
		if len(keyColumnList) == 0 {
			whereColumnList = columnList
		} else {
			for _, column := range whereColumnList {
				if param.AppendSqlValue {
					sql += "" + param.packingCharacterColumn(column.Name) + "=" + param.packingCharacterColumnStringValue(column, dataWhere[column.Name]) + " AND "
				} else {
					sql += "" + param.packingCharacterColumn(column.Name) + "=? AND "
					values = append(values, param.formatColumnValue(column, dataWhere[column.Name]))
				}
			}
		}
		sql = strings.TrimSuffix(sql, " AND ")

		sqlList = append(sqlList, sql)
		valuesList = append(valuesList, values)
	}
	return
}

type Where struct {
	Name                    string `json:"name"`
	Value                   string `json:"value"`
	Before                  string `json:"before"`
	After                   string `json:"after"`
	CustomSql               string `json:"customSql"`
	SqlConditionalOperation string `json:"sqlConditionalOperation"`
	AndOr                   string `json:"andOr"`
}

type Order struct {
	Name    string `json:"name"`
	DescAsc string `json:"descAsc"`
}

func DataListSelectSql(param *GenerateParam, database string, table string, columnList []*TableColumnModel, whereList []*Where, orderList []*Order) (sql string, values []interface{}, err error) {
	selectColumns := ""
	for _, column := range columnList {
		selectColumns += param.packingCharacterColumn(column.Name) + ","
	}
	selectColumns = strings.TrimSuffix(selectColumns, ",")
	if selectColumns == "" {
		selectColumns = "*"
	}
	sql = "SELECT " + selectColumns + " FROM "

	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += param.packingCharacterTable(table)

	//构造查询用的finder
	if len(whereList) > 0 {
		sql += " WHERE"
		for index, where := range whereList {
			sql += " " + param.packingCharacterColumn(where.Name)
			value := where.Value
			switch where.SqlConditionalOperation {
			case "like":
				sql += " LIKE ?"
				values = append(values, "%"+value+"%")
			case "not like":
				sql += " NOT LIKE ?"
				values = append(values, "%"+value+"%")
			case "like start":
				sql += " LIKE ?"
				values = append(values, ""+value+"%")
			case "not like start":
				sql += " NOT LIKE ?"
				values = append(values, ""+value+"%")
			case "like end":
				sql += " LIKE ?"
				values = append(values, "%"+value+"")
			case "not like end":
				sql += " NOT LIKE ?"
				values = append(values, "%"+value+"")
			case "is null":
				sql += " IS NULL"
			case "is not null":
				sql += " IS NOT NULL"
			case "is empty":
				sql += " = ?"
				values = append(values, "")
			case "is not empty":
				sql += " <> ?"
				values = append(values, "")
			case "between":
				sql += " BETWEEN ? AND ?"
				values = append(values, where.Before, where.After)
			case "not between":
				sql += " NOT BETWEEN ? AND ?"
				values = append(values, where.Before, where.After)
			case "in":
				sql += " IN (?)"
				values = append(values, value)
			case "not in":
				sql += " NOT IN (?)"
				values = append(values, value)
			default:
				sql += " " + where.SqlConditionalOperation + " ?"
				values = append(values, value)
			}
			// params_ = append(params_, where.Value)
			if index < len(whereList)-1 {
				sql += " " + where.AndOr + " "
			}
		}
	}
	if len(orderList) > 0 {
		sql += " ORDER BY"
		for index, order := range orderList {
			sql += " " + param.packingCharacterColumn(order.Name)
			if order.DescAsc != "" {
				sql += " " + order.DescAsc
			}
			// params_ = append(params_, where.Value)
			if index < len(orderList)-1 {
				sql += ","
			}
		}

	}
	return
}
