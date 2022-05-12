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
				sql += "" + param.packingCharacterColumn(column.Name) + " = " + param.packingCharacterColumnStringValue(column, value) + ", "
			} else {
				sql += "" + param.packingCharacterColumn(column.Name) + " = ?, "
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
					sql += "" + param.packingCharacterColumn(column.Name) + " = " + param.packingCharacterColumnStringValue(column, dataWhere[column.Name]) + " AND "
				} else {
					sql += "" + param.packingCharacterColumn(column.Name) + " = ? AND "
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
					sql += "" + param.packingCharacterColumn(column.Name) + " = " + param.packingCharacterColumnStringValue(column, dataWhere[column.Name]) + " AND "
				} else {
					sql += "" + param.packingCharacterColumn(column.Name) + " = ? AND "
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
