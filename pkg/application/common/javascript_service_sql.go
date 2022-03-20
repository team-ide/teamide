package common

import (
	"regexp"
	"strings"
	"teamide/pkg/application/base"
	model2 "teamide/pkg/application/model"
)

func getJavascriptBySqlInsert(app IApplication, sqlInsert *model2.SqlInsert, tab int) (javascript string, err error) {

	wrapTable := app.GetScript().WrapTableName(sqlInsert.Database, sqlInsert.Table)
	base.AppendLine(&javascript, `$invoke_temp.sql = "INSERT INTO `+wrapTable+` "`, tab)
	base.AppendLine(&javascript, `$invoke_temp.params = []`, tab)
	base.AppendLine(&javascript, `$invoke_temp.columnsSql = ""`, tab)
	base.AppendLine(&javascript, `$invoke_temp.valuesSql = ""`, tab)
	javascript += "\n"

	for _, column := range sqlInsert.Columns {
		if app.GetScript().IsEmpty(column.Name) {
			continue
		}
		var if_ string = column.IfScript
		if base.IsNotEmpty(if_) {
			base.AppendLine(&javascript, `if (`+if_+`) { `, tab)
			tab++
		}
		if column.AutoIncrement {
			continue
		}
		valueScript := column.Value
		if base.IsEmpty(valueScript) {
			valueScript = column.Name
		}

		wrapColumn := app.GetScript().WrapColumnName("", column.Name)
		if column.IgnoreEmpty {
			base.AppendLine(&javascript, `if (isNotEmpty(`+valueScript+`)) { `, tab)
			base.AppendLine(&javascript, `$invoke_temp.columnsSql = $invoke_temp.columnsSql + "`+wrapColumn+`, "`, tab+1)
			base.AppendLine(&javascript, `$invoke_temp.valuesSql = $invoke_temp.valuesSql+ "?, "`, tab+1)
			base.AppendLine(&javascript, `$invoke_temp.params.push(`+valueScript+`)`, tab+1)
			base.AppendLine(&javascript, `} `, tab)
		} else {
			base.AppendLine(&javascript, `$invoke_temp.columnsSql = $invoke_temp.columnsSql + "`+wrapColumn+`, "`, tab)
			base.AppendLine(&javascript, `$invoke_temp.valuesSql = $invoke_temp.valuesSql+ "?, "`, tab)
			base.AppendLine(&javascript, `$invoke_temp.params.push(`+valueScript+`)`, tab)
		}

		if base.IsNotEmpty(if_) {
			tab--
			base.AppendLine(&javascript, `} `, tab)
		}
		javascript += "\n"
	}
	base.AppendLine(&javascript, `// 去除多余的符号`, tab)
	base.AppendLine(&javascript, `$invoke_temp.columnsSql = trimSuffix($invoke_temp.columnsSql, ", ")`, tab)
	base.AppendLine(&javascript, `$invoke_temp.valuesSql = trimSuffix($invoke_temp.valuesSql, ", ")`, tab)
	base.AppendLine(&javascript, `// 组合SQL`, tab)
	base.AppendLine(&javascript, `$invoke_temp.sql = $invoke_temp.sql + "(" + $invoke_temp.columnsSql + ") VALUES (" + $invoke_temp.valuesSql + ")"`, tab)

	return
}

func getJavascriptBySqlUpdate(app IApplication, sqlUpdate *model2.SqlUpdate, tab int) (javascript string, err error) {

	wrapTable := app.GetScript().WrapTableName(sqlUpdate.Database, sqlUpdate.Table)
	base.AppendLine(&javascript, `$invoke_temp.sql = "UPDATE `+wrapTable+` SET "`, tab)
	base.AppendLine(&javascript, `$invoke_temp.params = []`, tab)
	javascript += "\n"

	for _, column := range sqlUpdate.Columns {
		if app.GetScript().IsEmpty(column.Name) {
			continue
		}
		var if_ string = column.IfScript
		if base.IsNotEmpty(if_) {
			base.AppendLine(&javascript, `if (`+if_+`) { `, tab)
			tab++
		}
		valueScript := column.Value
		if base.IsEmpty(valueScript) {
			valueScript = column.Name
		}

		wrapColumn := app.GetScript().WrapColumnName("", column.Name)
		if column.IgnoreEmpty {
			base.AppendLine(&javascript, `if (isNotEmpty(`+valueScript+`)) { `, tab)
			base.AppendLine(&javascript, `$invoke_temp.sql = $invoke_temp.sql + "`+wrapColumn+` = ?, "`, tab+1)
			base.AppendLine(&javascript, `$invoke_temp.params.push(`+valueScript+`)`, tab+1)
			base.AppendLine(&javascript, `} `, tab)
		} else {
			base.AppendLine(&javascript, `$invoke_temp.sql = $invoke_temp.sql + "`+wrapColumn+` = ?, "`, tab)
			base.AppendLine(&javascript, `$invoke_temp.params.push(`+valueScript+`)`, tab)
		}

		if base.IsNotEmpty(if_) {
			tab--
			base.AppendLine(&javascript, "} ", tab)
		}
		javascript += "\n"
	}
	base.AppendLine(&javascript, `// 去除多余的符号`, tab)
	base.AppendLine(&javascript, `$invoke_temp.sql = trimSuffix($invoke_temp.sql, ", ")`, tab)

	var javascript_ string
	javascript_, err = getJavascriptBySqlWheres(app, sqlUpdate.Wheres, tab)
	if err != nil {
		return
	}
	if base.IsNotEmpty(javascript_) {
		javascript += "\n"
		javascript += javascript_
		javascript += "\n"
		base.AppendLine(&javascript, `// 组合条件`, tab)
		base.AppendLine(&javascript, `$invoke_temp.sql = $invoke_temp.sql + $invoke_temp.whereSql`, tab)
		base.AppendLine(&javascript, `$invoke_temp.params.pushs($invoke_temp.params, $invoke_temp.whereParams)`, tab)
	}
	return
}

func getJavascriptBySqlDelete(app IApplication, sqlDelete *model2.SqlDelete, tab int) (javascript string, err error) {
	wrapTable := app.GetScript().WrapTableName(sqlDelete.Database, sqlDelete.Table)

	base.AppendLine(&javascript, `$invoke_temp.sql = "DELETE FROM `+wrapTable+` "`, tab)
	base.AppendLine(&javascript, `$invoke_temp.params = []`, tab)

	var javascript_ string
	javascript_, err = getJavascriptBySqlWheres(app, sqlDelete.Wheres, tab)
	if err != nil {
		return
	}
	if base.IsNotEmpty(javascript_) {
		javascript += "\n"
		javascript += javascript_
		javascript += "\n"
		base.AppendLine(&javascript, `// 组合条件`, tab)
		base.AppendLine(&javascript, `$invoke_temp.sql = $invoke_temp.sql + $invoke_temp.whereSql`, tab)
		base.AppendLine(&javascript, `$invoke_temp.params.pushs($invoke_temp.params, $invoke_temp.whereParams)`, tab)
	}
	return
}

func getJavascriptBySqlSelect(app IApplication, sqlSelect *model2.SqlSelect, tab int) (javascript string, err error) {

	wrapTable := app.GetScript().WrapTableName(sqlSelect.Database, sqlSelect.Table)

	columnHasIf := false
	noIfColumns := ""
	for _, one := range sqlSelect.Columns {
		if base.IsNotEmpty(one.IfScript) {
			columnHasIf = true
			continue
		}
		asAlias := ""
		if base.IsNotEmpty(one.Alias) {
			asAlias = " AS " + one.Alias
		}
		if one.Custom {
			noIfColumns += "(" + one.CustomSql + ")" + asAlias + ", "
		} else {
			if base.IsEmpty(one.Name) {
				continue
			}
			wrapColumn := app.GetScript().WrapColumnName(one.TableAlias, one.Name)
			noIfColumns += wrapColumn + asAlias + ", "
		}
	}

	noIfColumns = strings.TrimSuffix(noIfColumns, ", ")

	if columnHasIf {
		base.AppendLine(&javascript, `$invoke_temp.columnsSql = ""`, tab)
		for _, one := range sqlSelect.Columns {
			if base.IsEmpty(one.IfScript) {
				continue
			}
			columnHasIf = true
			base.AppendLine(&javascript, `if (isNotEmpty(`+one.IfScript+`)) { `, tab)
			tab++
			asAlias := ""
			if base.IsNotEmpty(one.Alias) {
				asAlias = " AS " + one.Alias
			}
			if one.Custom {
				base.AppendLine(&javascript, `$invoke_temp.columnsSql = $invoke_temp.columnsSql + "(`+one.CustomSql+`)`+asAlias+`, "`, tab)
			} else {
				if base.IsEmpty(one.Name) {
					continue
				}
				wrapColumn := app.GetScript().WrapColumnName(one.TableAlias, one.Name)
				base.AppendLine(&javascript, `$invoke_temp.columnsSql = $invoke_temp.columnsSql + "`+wrapColumn+``+asAlias+`, "`, tab)
				noIfColumns += wrapColumn
			}
			tab--
			base.AppendLine(&javascript, "} ", tab)
		}
		base.AppendLine(&javascript, `// 去除多余的符号`, tab)
		base.AppendLine(&javascript, `$invoke_temp.columnsSql = trimSuffix($invoke_temp.columnsSql, ", ")`, tab)
	}

	if len(sqlSelect.Columns) == 0 {
		noIfColumns = "*"
	}
	DistinctStr := ""
	if sqlSelect.Distinct {
		DistinctStr = "DISTINCT "
	}
	if !sqlSelect.SelectCount {
		asAlias := ""
		if base.IsNotEmpty(sqlSelect.Alias) {
			asAlias = " AS " + sqlSelect.Alias
		}
		if columnHasIf {
			base.AppendLine(&javascript, `$invoke_temp.sql = "SELECT`+DistinctStr+` `+noIfColumns+` " + $invoke_temp.columnsSql + " FROM `+wrapTable+asAlias+` "`, tab)
			base.AppendLine(&javascript, `$invoke_temp.params = []`, tab)
		} else {
			base.AppendLine(&javascript, `$invoke_temp.sql = "SELECT`+DistinctStr+` `+noIfColumns+` FROM `+wrapTable+asAlias+` "`, tab)
			base.AppendLine(&javascript, `$invoke_temp.params = []`, tab)
		}
	}
	if sqlSelect.SelectPage || sqlSelect.SelectCount {
		base.AppendLine(&javascript, `$invoke_temp.countSql = "SELECT`+DistinctStr+` COUNT(*) FROM `+wrapTable+` "`, tab)
		base.AppendLine(&javascript, `$invoke_temp.countParams = []`, tab)
	}

	for _, one := range sqlSelect.LeftJoin {
		if base.IsNotEmpty(one.IfScript) {
			base.AppendLine(&javascript, `if (isNotEmpty(`+one.IfScript+`)) { `, tab)
			tab++
		}
		if one.Custom {
			base.AppendLine(&javascript, `$invoke_temp.sql = $invoke_temp.sql + "LEFT JOIN `+one.CustomSql+` "`, tab)
		} else {
			asAlias := ""
			if base.IsNotEmpty(one.Alias) {
				asAlias = " AS " + one.Alias
			}
			base.AppendLine(&javascript, `$invoke_temp.sql = $invoke_temp.sql + "LEFT JOIN `+one.Table+asAlias+` ON `+one.On+` "`, tab)
		}

		if base.IsNotEmpty(one.IfScript) {
			tab--
			base.AppendLine(&javascript, "} ", tab)
		}
	}
	var javascript_ string
	javascript_, err = getJavascriptBySqlWheres(app, sqlSelect.Wheres, tab)
	if err != nil {
		return
	}
	if base.IsNotEmpty(javascript_) {
		javascript += "\n"
		javascript += javascript_
		javascript += "\n"
		if !sqlSelect.SelectCount {
			base.AppendLine(&javascript, `// 组合条件`, tab)
			base.AppendLine(&javascript, `$invoke_temp.sql = $invoke_temp.sql + $invoke_temp.whereSql`, tab)
			base.AppendLine(&javascript, `$invoke_temp.params.pushs($invoke_temp.params, $invoke_temp.whereParams)`, tab)
		}
		if sqlSelect.SelectPage || sqlSelect.SelectCount {
			base.AppendLine(&javascript, `// 组合条件`, tab)
			base.AppendLine(&javascript, `$invoke_temp.countSql = $invoke_temp.countSql + $invoke_temp.whereSql`, tab)
			base.AppendLine(&javascript, `$invoke_temp.countParams.pushs($invoke_temp.countParams, $invoke_temp.whereParams)`, tab)
		}
	}

	return
}

func getJavascriptBySqlWheres(app IApplication, wheres []*model2.SqlWhere, tab int) (javascript string, err error) {
	if len(wheres) == 0 {
		return
	}
	base.AppendLine(&javascript, `$invoke_temp.whereSql = ""`, tab)
	base.AppendLine(&javascript, `$invoke_temp.whereParams = []`, tab)

	err = appendJavascriptBySqlWheres(app, &javascript, wheres, tab)

	if err != nil {
		return
	}
	base.AppendLine(&javascript, `// 去除多余的符号`, tab)
	base.AppendLine(&javascript, `$invoke_temp.whereSql = trimPrefix($invoke_temp.whereSql, "AND ")`, tab)
	base.AppendLine(&javascript, `$invoke_temp.whereSql = trimPrefix($invoke_temp.whereSql, "OR ")`, tab)
	base.AppendLine(&javascript, `$invoke_temp.whereSql = "WHERE " + $invoke_temp.whereSql`, tab)
	return
}

func appendJavascriptBySqlWheres(app IApplication, javascript *string, wheres []*model2.SqlWhere, tab int) (err error) {
	for _, one := range wheres {
		err = appendJavascriptBySqlWhere(app, javascript, one, tab)
		if err != nil {
			return
		}
	}
	return
}
func formatJavascriptByCustomSql(app IApplication, customSql string) (sql string, valueScripts []string, err error) {
	if base.IsEmpty(customSql) {
		return
	}
	sql = ""
	valueScripts = []string{}
	var re *regexp.Regexp
	re, err = regexp.Compile(`{(.+?)}`)
	if err != nil {
		return
	}
	indexsList := re.FindAllIndex([]byte(customSql), -1)
	var lastIndex int = 0
	for _, indexs := range indexsList {
		sql += customSql[lastIndex:indexs[0]]
		lastIndex = indexs[1]

		sql += "?"
		script := customSql[indexs[0]+1 : indexs[1]-1]
		valueScripts = append(valueScripts, script)
	}
	sql += customSql[lastIndex:]
	return
}
func appendJavascriptBySqlWhere(app IApplication, javascript *string, where *model2.SqlWhere, tab int) (err error) {
	if where.Piece {
		base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` ("`, tab)
		base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + ")"`, tab)
	} else {
		if where.Custom {
			var customSql string
			var valueScripts []string
			customSql, valueScripts, err = formatJavascriptByCustomSql(app, where.CustomSql)
			if err != nil {
				return
			}
			if !strings.HasPrefix(strings.ToUpper(customSql), "AND ") && !strings.HasPrefix(strings.ToUpper(customSql), "OR ") {
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+customSql+`"`, tab)
			} else {
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+customSql+`"`, tab)
			}
			for _, valueScript := range valueScripts {
				base.AppendLine(javascript, `$invoke_temp.whereParams.push(`+valueScript+`)`, tab)
			}
		} else {

			valueScript := where.Value
			if base.IsEmpty(valueScript) {
				valueScript = where.Name
			}
			wrapColumn := app.GetScript().WrapColumnName(where.TableAlias, where.Name)
			operator := where.GetOperator()
			switch operator {
			case model2.IS_NULL:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` IS NULL "`, tab)
			case model2.IS_NOT_NULL:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` IS NOT NULL "`, tab)
			case model2.IS_EMPTY:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` = '' "`, tab)
			case model2.IS_NOT_EMPTY:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` <> '' "`, tab)
			case model2.LIKE:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` LIKE ? "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push("%" + `+valueScript+` + "%")`, tab)
			case model2.NOT_LIKE:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` NOT LIKE ? "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push("%" + `+valueScript+` + "%")`, tab)
			case model2.LIKE_BEFORE:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` LIKE ? "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push(`+valueScript+` + "%")`, tab)
			case model2.NOT_LIKE_BEFORE:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` NOT LIKE ? "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push(`+valueScript+` + "%")`, tab)
			case model2.LIKE_AFTER:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` LIKE ? "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push("%" + `+valueScript+`)`, tab)
			case model2.NOT_LIKE_AFTER:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` NOT LIKE ? "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push("%" + `+valueScript+`)`, tab)
			case model2.IN:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` IN (?) "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push(`+valueScript+`)`, tab)
			case model2.NOT_IN:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` NOT IN (?) "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push(`+valueScript+`)`, tab)
			case model2.IN_LIKE:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` IN LIKE (?) "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push(`+valueScript+`)`, tab)
			case model2.NOT_IN:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` NOT IN LIKE (?) "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push(`+valueScript+`)`, tab)
			default:
				base.AppendLine(javascript, `$invoke_temp.whereSql = $invoke_temp.whereSql + "`+where.GetAndOr()+` `+wrapColumn+` `+operator.Value+` ? "`, tab)
				base.AppendLine(javascript, `$invoke_temp.whereParams.push(`+valueScript+`)`, tab)
			}
		}
	}
	return
}
