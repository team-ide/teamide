package sql

import "strings"

type Select struct {
	Database     string          `json:"database"`     // 库名
	Table        string          `json:"table"`        // 表名
	Distinct     bool            `json:"distinct"`     // SELECT DISTINCT 语句  列出不同的值
	Columns      []*SelectColumn `json:"columns"`      // 字段
	Intos        []*Into         `json:"intos"`        // SELECT INTO 复制表
	InnerJoins   []*InnerJoin    `json:"innerJoins"`   // INNER JOIN
	LeftJoin     []*LeftJoin     `json:"leftJoin"`     // LEFT JOIN
	RightJoin    []*RightJoin    `json:"rightJoin"`    // RIGHT JOIN
	FullJoin     []*FullJoin     `json:"fullJoin"`     // FULL JOIN
	Wheres       []*Where        `json:"wheres"`       // 条件
	Orders       []*Order        `json:"orders"`       // 排序
	Groups       []*Group        `json:"groups"`       // 分组
	Havings      []*Having       `json:"havings"`      // 筛选分组后的各组数据
	UnionSelects []*Select       `json:"unionSelects"` // UNION 操作符 合并两个或多个 SELECT 语句的结果集
}

type SelectColumn struct {
	IfScript   string `json:"ifScript"`   // 条件  满足该条件 添加
	Custom     bool   `json:"custom"`     // 是否自定义
	CustomSql  string `json:"customSql"`  // 是否自定义
	TableAlias string `json:"tableAlias"` // 表别名
	Name       string `json:"name"`       // 字段名称
	Alias      string `json:"alias"`      // 字段别名
}

type Into struct {
	IfScript  string `json:"ifScript"`  // 条件  满足该条件 添加
	Custom    bool   `json:"custom"`    // 是否自定义
	CustomSql string `json:"customSql"` // 是否自定义
	Table     string `json:"table"`     // 表名
}

type InnerJoin struct {
	IfScript  string `json:"ifScript"`  // 条件  满足该条件 添加
	Custom    bool   `json:"custom"`    // 是否自定义
	CustomSql string `json:"customSql"` // 是否自定义
	Table     string `json:"table"`     // 表名
	Alias     string `json:"alias"`     // 别名
	On        string `json:"on"`        // ON
}

type LeftJoin struct {
	IfScript  string `json:"ifScript"`  // 条件  满足该条件 添加
	Custom    bool   `json:"custom"`    // 是否自定义
	CustomSql string `json:"customSql"` // 是否自定义
	Table     string `json:"table"`     // 表名
	Alias     string `json:"alias"`     // 别名
	On        string `json:"on"`        // ON
}

type RightJoin struct {
	IfScript  string `json:"ifScript"`  // 条件  满足该条件 添加
	Custom    bool   `json:"custom"`    // 是否自定义
	CustomSql string `json:"customSql"` // 是否自定义
	Table     string `json:"table"`     // 表名
	Alias     string `json:"alias"`     // 别名
	On        string `json:"on"`        // ON
}

type FullJoin struct {
	IfScript  string `json:"ifScript"`  // 条件  满足该条件 添加
	Custom    bool   `json:"custom"`    // 是否自定义
	CustomSql string `json:"customSql"` // 是否自定义
	Table     string `json:"table"`     // 表名
	Alias     string `json:"alias"`     // 别名
	On        string `json:"on"`        // ON
}

type Order struct {
	IfScript   string `json:"ifScript"`   // 条件  满足该条件 添加
	Custom     bool   `json:"custom"`     // 是否自定义
	CustomSql  string `json:"customSql"`  // 是否自定义
	TableAlias string `json:"tableAlias"` // 表别名
	Name       string `json:"name"`       // 字段名称
	AscDesc    string `json:"ascDesc"`    // 默认按照升序 升序ASC 降序DESC
}

type Group struct {
	IfScript   string `json:"ifScript"`   // 条件  满足该条件 添加
	Custom     bool   `json:"custom"`     // 是否自定义
	CustomSql  string `json:"customSql"`  // 是否自定义
	TableAlias string `json:"tableAlias"` // 表别名
	Name       string `json:"name"`       // 字段名称
}

type Having struct {
	IfScript  string `json:"ifScript"`  // 条件  满足该条件 添加
	Custom    bool   `json:"custom"`    // 是否自定义
	CustomSql string `json:"customSql"` // 是否自定义
	Having    string `json:"having"`    // 筛选分组后的各组数据
}

func (this_ *Select) GetSqlParam(data map[string]interface{}) (sqlParam SqlParam, err error) {

	var sql string
	var params []interface{}

	sql, params, err = getSelectSqlParam(data, this_)
	if err != nil {
		return
	}
	sqlParam.Sql = sql
	sqlParam.Params = params
	return
}

func (this_ *Select) GetSqlParams(dataList ...map[string]interface{}) (sqlParams []SqlParam, err error) {
	if len(dataList) == 0 {
		return
	}
	for _, data := range dataList {
		var sqlParam SqlParam
		sqlParam, err = this_.GetSqlParam(data)
		if err != nil {
			return
		}
		sqlParams = append(sqlParams, sqlParam)
	}

	return
}

func getSelectSqlParam(data map[string]interface{}, model *Select) (sql string, params []interface{}, err error) {
	wrapTable := WrapTableName(model.Database, model.Table)

	params = []interface{}{}

	sql = "SELECT"

	if model.Distinct {
		sql += " DISTINCT"
	}

	var sql_ string
	var param_ []interface{}

	sql_, param_, err = getSelectColumnSqlParam(data, model.Columns)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " " + sql_
		params = append(params, param_...)
	}

	sql_, param_, err = getIntoSqlParam(data, model.Intos)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " INTO " + sql_
		params = append(params, param_...)
	}

	if IsNotEmpty(wrapTable) {
		sql += " FROM " + wrapTable
	}

	sql_, param_, err = getInnerJoinSqlParam(data, model.InnerJoins)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " INNER JOIN " + sql_
		params = append(params, param_...)
	}

	sql_, param_, err = getLeftJoinSqlParam(data, model.LeftJoin)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " LEFT JOIN " + sql_
		params = append(params, param_...)
	}

	sql_, param_, err = getRightJoinSqlParam(data, model.RightJoin)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " RIGHT JOIN " + sql_
		params = append(params, param_...)
	}

	sql_, param_, err = getFullJoinSqlParam(data, model.FullJoin)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " FULL JOIN " + sql_
		params = append(params, param_...)
	}

	sql_, param_, err = getWhereSqlParam(data, model.Wheres)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " WHERE " + sql_
		params = append(params, param_...)
	}

	sql_, param_, err = getOrderSqlParam(data, model.Orders)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " ORDER BY " + sql_
		params = append(params, param_...)
	}

	sql_, param_, err = getGroupSqlParam(data, model.Groups)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " GROUP BY " + sql_
		params = append(params, param_...)
	}

	sql_, param_, err = getHavingSqlParam(data, model.Havings)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " HAVING " + sql_
		params = append(params, param_...)
	}

	sql_, param_, err = getUnionSqlParam(data, model.UnionSelects)
	if err != nil {
		return
	}
	if IsNotEmpty(sql_) {
		sql += " UNION " + sql_
		params = append(params, param_...)
	}

	return
}
func getSelectColumnSqlParam(data map[string]interface{}, columns []*SelectColumn) (sql string, params []interface{}, err error) {
	if len(columns) == 0 {
		return
	}
	for _, one := range columns {
		if !IfScriptValue(data, one.IfScript) {
			continue
		}
		if one.Custom {
			sql += one.CustomSql + ", "
			continue
		}
		if IsEmpty(one.Name) {
			continue
		}
		wrapColumn := WrapColumnName(one.TableAlias, one.Name)
		sql += wrapColumn
		if IsNotEmpty(one.Alias) {
			sql += " AS " + one.Alias
		}
		sql += ", "

	}
	sql = strings.TrimSuffix(sql, ", ")

	return
}
func getIntoSqlParam(data map[string]interface{}, list []*Into) (sql string, params []interface{}, err error) {
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		if !IfScriptValue(data, one.IfScript) {
			continue
		}
		if one.Custom {
			sql += one.CustomSql + "INTO "
			continue
		}
		wrapTable := WrapTableName("", one.Table)
		if IsNotEmpty(wrapTable) {
			sql += wrapTable
		}
		sql += "INTO "

	}
	sql = strings.TrimSuffix(sql, "INTO ")

	return
}
func getInnerJoinSqlParam(data map[string]interface{}, list []*InnerJoin) (sql string, params []interface{}, err error) {
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		if !IfScriptValue(data, one.IfScript) {
			continue
		}
		if one.Custom {
			sql += one.CustomSql + "INNER JOIN "
			continue
		}
		wrapTable := WrapTableName("", one.Table)
		sql += wrapTable
		if IsNotEmpty(one.On) {
			sql += " " + one.On
		}
		sql += "INNER JOIN "

	}
	sql = strings.TrimSuffix(sql, "INNER JOIN ")
	return
}
func getLeftJoinSqlParam(data map[string]interface{}, list []*LeftJoin) (sql string, params []interface{}, err error) {
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		if !IfScriptValue(data, one.IfScript) {
			continue
		}
		if one.Custom {
			sql += one.CustomSql + "LEFT JOIN "
			continue
		}
		wrapTable := WrapTableName("", one.Table)
		sql += wrapTable
		if IsNotEmpty(one.On) {
			sql += " " + one.On
		}
		sql += "LEFT JOIN "

	}
	sql = strings.TrimSuffix(sql, "LEFT JOIN ")
	return
}
func getRightJoinSqlParam(data map[string]interface{}, list []*RightJoin) (sql string, params []interface{}, err error) {
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		if !IfScriptValue(data, one.IfScript) {
			continue
		}
		if one.Custom {
			sql += one.CustomSql + "RIGHT JOIN "
			continue
		}
		wrapTable := WrapTableName("", one.Table)
		sql += wrapTable
		if IsNotEmpty(one.On) {
			sql += " " + one.On
		}
		sql += "RIGHT JOIN "

	}
	sql = strings.TrimSuffix(sql, "RIGHT JOIN ")
	return
}
func getFullJoinSqlParam(data map[string]interface{}, list []*FullJoin) (sql string, params []interface{}, err error) {
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		if !IfScriptValue(data, one.IfScript) {
			continue
		}
		if one.Custom {
			sql += one.CustomSql + "FULL JOIN "
			continue
		}
		wrapTable := WrapTableName("", one.Table)
		sql += wrapTable
		if IsNotEmpty(one.On) {
			sql += " " + one.On
		}
		sql += "FULL JOIN "

	}
	sql = strings.TrimSuffix(sql, "FULL JOIN ")
	return
}
func getOrderSqlParam(data map[string]interface{}, list []*Order) (sql string, params []interface{}, err error) {
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		if !IfScriptValue(data, one.IfScript) {
			continue
		}
		if one.Custom {
			sql += one.CustomSql + ", "
			continue
		}
		wrapColumn := WrapColumnName(one.TableAlias, one.Name)
		sql += wrapColumn
		if IsNotEmpty(one.AscDesc) {
			sql += " " + one.AscDesc
		}
		sql += ", "

	}
	sql = strings.TrimSuffix(sql, ", ")
	return
}
func getGroupSqlParam(data map[string]interface{}, list []*Group) (sql string, params []interface{}, err error) {
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		if !IfScriptValue(data, one.IfScript) {
			continue
		}
		if one.Custom {
			sql += one.CustomSql + ", "
			continue
		}
		wrapColumn := WrapColumnName(one.TableAlias, one.Name)
		sql += wrapColumn
		sql += ", "

	}
	sql = strings.TrimSuffix(sql, ", ")
	return
}
func getHavingSqlParam(data map[string]interface{}, list []*Having) (sql string, params []interface{}, err error) {
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		if !IfScriptValue(data, one.IfScript) {
			continue
		}
		if one.Custom {
			sql += one.CustomSql + ", "
			continue
		}
		sql += one.Having
		sql += ", "

	}
	sql = strings.TrimSuffix(sql, ", ")
	return
}
func getUnionSqlParam(data map[string]interface{}, list []*Select) (sql string, params []interface{}, err error) {
	if len(list) == 0 {
		return
	}
	for _, one := range list {
		var selectSql string
		var selectParams []interface{}
		selectSql, selectParams, err = getSelectSqlParam(data, one)
		if err != nil {
			return
		}
		if IsEmpty(selectSql) {
			continue
		}
		sql += selectSql
		params = append(params, selectParams...)
		sql += "UNION "

	}
	sql = strings.TrimSuffix(sql, "UNION ")
	return
}
