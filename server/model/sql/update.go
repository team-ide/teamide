package sqlModel

import (
	"fmt"
	"server/base"
	"strings"
)

type Update struct {
	Database string          `json:"database,omitempty"` // 库名
	Table    string          `json:"table,omitempty"`    // 表名
	Columns  []*UpdateColumn `json:"columns,omitempty"`  // 字段
	Wheres   []*Where        `json:"wheres,omitempty"`   // 条件
}

type UpdateColumn struct {
	IfScript    string `json:"ifScript,omitempty"`    // 条件  满足该条件 添加
	Custom      bool   `json:"custom,omitempty"`      // 是否自定义
	CustomSql   string `json:"customSql,omitempty"`   // 是否自定义
	Name        string `json:"name,omitempty"`        // 字段名称
	ValueScript string `json:"valueScript,omitempty"` // 字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值
	Required    bool   `json:"required,omitempty"`    // 必填
	AllowEmpty  bool   `json:"allowEmpty,omitempty"`  // 允许空值，如果是null或空字符串则也设置值
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

func (this_ *Update) GetTableColumns() (tableColumns map[string][]string) {
	tableColumns = map[string][]string{}

	wrapTable := WrapTableName(this_.Database, this_.Table)

	var columns []string
	for _, column := range this_.Columns {
		if IsEmpty(column.Name) {
			continue
		}
		wrapColumn := WrapColumnName("", column.Name)
		columns = append(columns, wrapColumn)
	}

	appendWhereColumns(this_.Wheres, &columns)

	tableColumns[wrapTable] = columns

	return
}
