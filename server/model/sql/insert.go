package sqlModel

import (
	"fmt"
	"server/base"
	"strings"
)

type Insert struct {
	Database string          `json:"database"` // 库名
	Table    string          `json:"table"`    // 表名
	Columns  []*InsertColumn `json:"columns"`  // 新增字段
}

type InsertColumn struct {
	IfScript      string `json:"ifScript"`      // 条件  满足该条件 添加
	Name          string `json:"name"`          // 字段名称
	ValueScript   string `json:"valueScript"`   // 字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值
	Required      bool   `json:"required"`      // 必填
	AutoIncrement bool   `json:"autoIncrement"` // 自增列
	AllowEmpty    bool   `json:"allowEmpty"`    // 允许空值，如果是null或空字符串则也设置值
}

func (this_ *Insert) AppendColumn(column *InsertColumn) *Insert {
	this_.Columns = append(this_.Columns, column)
	return this_
}

func (this_ *Insert) GetSqlParam(data map[string]interface{}) (sqlParam base.SqlParam, err error) {
	wrapTable := WrapTableName(this_.Database, this_.Table)

	params := []interface{}{}

	insertColumn := ""
	insertValue := ""

	for _, column := range this_.Columns {
		if IsEmpty(column.Name) {
			continue
		}
		if !IfScriptValue(data, column.IfScript) {
			continue
		}
		if column.AutoIncrement {
			continue
		}
		value := GetColumnValue(data, column.Name, column.ValueScript)

		if !column.AllowEmpty && IsEmptyObj(value) {
			continue
		}

		wrapColumn := WrapColumnName("", column.Name)
		insertColumn += wrapColumn + ", "
		insertValue += "?, "
		params = append(params, value)

	}

	insertColumn = strings.TrimSuffix(insertColumn, ", ")
	insertValue = strings.TrimSuffix(insertValue, ", ")

	sql := fmt.Sprint("INSERT INTO ", wrapTable, "(", insertColumn, ")VALUES(", insertValue, ")")

	sqlParam.Sql = sql
	sqlParam.Params = params
	return
}

func (this_ *Insert) GetSqlParams(dataList ...map[string]interface{}) (sqlParams []base.SqlParam, err error) {
	if len(dataList) == 0 {
		return
	}
	for _, data := range dataList {
		var sqlParam base.SqlParam
		sqlParam, err = this_.GetSqlParam(data)
		if err != nil {
			return
		}
		sqlParams = append(sqlParams, sqlParam)
	}

	return
}
