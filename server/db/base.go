package db

import (
	"base"
	"database/sql"
	"fmt"
	"reflect"
)

type SqlParam struct {
	Sql    string        `json:"sql,omitempty"`
	Params []interface{} `json:"params,omitempty"`
}

func NewSqlParam(sql_ string, params []interface{}) (sqlParam SqlParam) {
	if params == nil {
		params = []interface{}{}
	}
	sqlParam = SqlParam{
		Sql:    sql_,
		Params: params,
	}
	return
}

func InsertSqlByBean(table string, beans ...interface{}) (sqlParam SqlParam) {
	params := []interface{}{}
	if len(beans) == 0 {
		return
	}

	refType := base.GetRefType(beans[0])

	fieldCount := refType.NumField() // field count
	insertColumns := ""

	for i := 0; i < fieldCount; i++ {
		fieldType := refType.Field(i) // field type
		column := base.GetColumnNameByType(fieldType)
		if column != "" {
			insertColumns += column + ","
		}
	}
	insertValuesList := []string{}
	for _, bean := range beans {
		refValue := base.GetRefValue(bean) // value
		insertValues := ""
		for i := 0; i < fieldCount; i++ {
			fieldType := refType.Field(i)   // field type
			fieldValue := refValue.Field(i) // field vlaue
			column := base.GetColumnNameByType(fieldType)
			if column == "" {
				continue
			}
			value := base.GetFieldTypeValue(fieldType.Type, fieldValue)
			switch fieldType.Type.Name() {
			case "string":
				if value == "" {
					insertValues += "NULL,"
					continue
				}
				insertValues += "?,"
				params = append(params, value)
			default:
				if value == nil {
					insertValues += "NULL,"
					continue
				} else if base.IsZero(value) {
					insertValues += "NULL,"
					continue
				}
				insertValues += "?,"
				params = append(params, value)
			}
		}
		if len(insertValues) > 0 {
			insertValuesList = append(insertValuesList, insertValues)
		}
	}
	if len(insertColumns) > 0 {
		insertColumns = insertColumns[0 : len(insertColumns)-1]
	}
	sql := " INSERT INTO " + table + "(" + insertColumns + ") VALUES "
	for i, values := range insertValuesList {
		if i > 0 {
			sql += ","
		}
		values = values[0 : len(values)-1]
		sql += "("
		sql += values
		sql += ")"

	}
	sqlParam = SqlParam{sql, params}
	return sqlParam
}

func GetColumnSqlByBean(refType reflect.Type, alias string) (columnSql string) {
	columnSql = ""
	fieldCount := refType.NumField() // field count
	for i := 0; i < fieldCount; i++ {
		fieldType := refType.Field(i) // field type
		column := base.GetColumnNameByType(fieldType)
		if column == "" {
			continue
		}
		if alias != "" {
			columnSql += alias + "."
		}
		columnSql += column + ","
	}
	columnSql = columnSql[0 : len(columnSql)-1]
	return
}

func ResultToBeans(rows *sql.Rows, newBean func() interface{}) (list []interface{}, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	bean := newBean()
	columnTypes := base.GetColumnTypes(bean)

	list = []interface{}{}
	for rows.Next() {
		values := []interface{}{}
		beanColumnTypes := []base.ColumnType{}

		for _, column := range columns {
			columnType := columnTypes[column]

			if columnType.FieldType == nil {
				beanColumnTypes = append(beanColumnTypes, base.ColumnType{})
				continue
			}

			beanColumnTypes = append(beanColumnTypes, columnType)
			typeName := columnType.FieldType.Type.Name()
			if typeName == "int64" {
				var value sql.NullInt64
				values = append(values, &value)
			} else if typeName == "int32" {
				var value sql.NullInt32
				values = append(values, &value)
			} else if typeName == "int16" {
				var value sql.NullInt32
				values = append(values, &value)
			} else if typeName == "int8" {
				var value sql.NullInt32
				values = append(values, &value)
			} else if typeName == "int" {
				var value sql.NullInt32
				values = append(values, &value)
			} else if typeName == "float64" {
				var value sql.NullFloat64
				values = append(values, &value)
			} else if typeName == "float32" {
				var value sql.NullFloat64
				values = append(values, &value)
			} else if typeName == "bool" {
				var value sql.NullBool
				values = append(values, &value)
			} else if typeName == "time.Time" || typeName == "Time" {
				var value sql.NullTime
				values = append(values, &value)
			} else {
				var value sql.NullString
				values = append(values, &value)
			}
		}
		err = rows.Scan(values...)
		if err != nil {
			fmt.Println("ResultToBeans error:", err)
			return nil, err
		}
		refValue := base.GetRefValue(bean)
		for index := range columns {
			beanColumnType := beanColumnTypes[index]
			if beanColumnType.Column == "" {
				continue
			}

			value := values[index]
			typeName := beanColumnType.FieldType.Type.Name()
			if value == nil {
				continue
			}
			if typeName == "int64" {
				v := value.(*sql.NullInt64)
				value = v.Int64
			} else if typeName == "int32" {
				v := value.(*sql.NullInt32)
				value = v.Int32
			} else if typeName == "int16" {
				v := value.(*sql.NullInt32)
				value = int16(v.Int32)
			} else if typeName == "int8" {
				v := value.(*sql.NullInt32)
				value = int8(v.Int32)
			} else if typeName == "int" {
				v := value.(*sql.NullInt32)
				value = int(v.Int32)
			} else if typeName == "float64" {
				v := value.(*sql.NullFloat64)
				value = v.Float64
			} else if typeName == "float32" {
				v := value.(*sql.NullFloat64)
				value = float32(v.Float64)
			} else if typeName == "bool" {
				v := value.(*sql.NullBool)
				value = v.Bool
			} else if typeName == "time.Time" || typeName == "Time" {
				v := value.(*sql.NullTime)
				value = v.Time
			} else {
				v := value.(*sql.NullString)
				value = v.String
			}
			val := reflect.ValueOf(value)
			refValue.FieldByName(beanColumnType.FieldType.Name).Set(val)
		}

		list = append(list, bean)

		bean = newBean()
	}
	return list, err
}
