package sqlModel

import (
	"fmt"
	"server/base"
	"strings"
)

type Update struct {
	Database string          `json:"database"` // 库名
	Table    string          `json:"table"`    // 表名
	Columns  []*UpdateColumn `json:"columns"`  // 字段
	Wheres   []*Where        `json:"wheres"`   // 条件
}

type UpdateColumn struct {
	IfScript    string `json:"ifScript"`    // 条件  满足该条件 添加
	Custom      bool   `json:"custom"`      // 是否自定义
	CustomSql   string `json:"customSql"`   // 是否自定义
	Name        string `json:"name"`        // 字段名称
	ValueScript string `json:"valueScript"` // 字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值
	Required    bool   `json:"required"`    // 必填
	AllowEmpty  bool   `json:"allowEmpty"`  // 允许空值，如果是null或空字符串则也设置值
}

func (this_ *Update) GetSqlParam(data map[string]interface{}) (sqlParam base.SqlParam, err error) {
	wrapTable := WrapTableName(this_.Database, this_.Table)

	params := []interface{}{}

	updateSql := ""

	for _, column := range this_.Columns {

		if !IfScriptValue(data, column.IfScript) {
			continue
		}
		if column.Custom {
			updateSql += column.CustomSql
			continue
		}
		if IsEmpty(column.Name) {
			continue
		}

		value := GetColumnValue(data, column.Name, column.ValueScript)
		if !column.AllowEmpty && IsEmptyObj(value) {
			continue
		}

		wrapColumn := WrapColumnName("", column.Name)
		updateSql += wrapColumn + " = ?, "
		params = append(params, value)

	}

	updateSql = strings.TrimSuffix(updateSql, ", ")

	sql := fmt.Sprint("UPDATE ", wrapTable, " SET ", updateSql)

	whereSql := ""
	var whereParams []interface{}
	whereSql, whereParams, err = getWhereSqlParam(data, this_.Wheres)
	if err != nil {
		return
	}
	if IsNotEmpty(whereSql) {
		sql += " WHERE " + whereSql
		params = append(params, whereParams...)
	}

	sqlParam.Sql = sql
	sqlParam.Params = params
	return
}

func (this_ *Update) GetSqlParams(dataList ...map[string]interface{}) (sqlParams []base.SqlParam, err error) {
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
