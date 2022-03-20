package common

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"teamide/pkg/application/base"
	model2 "teamide/pkg/application/model"
)

func GetJavascriptMethodName(name string) (methodName string) {
	methodName = strings.ReplaceAll(name, ".", "_")
	chars := strings.Split(methodName, "")
	needToUp := false
	methodName = ""
	for _, char := range chars {
		if char == "/" {
			needToUp = true
		} else {
			if char == "." {
				char = "_"
			}
			if needToUp {
				char = strings.ToUpper(char)
				needToUp = false
			}
			methodName += char
		}
	}

	return
}

func GetActionJavascriptByAction(app IApplication, action *model2.ActionModel) (javascript string, err error) {
	methodName := GetJavascriptMethodName(action.Name)
	javascript += ""
	javascript += "function action_" + methodName + "("
	for _, inVariable := range action.InVariables {
		javascript += inVariable.Name + ", "
	}
	javascript = strings.TrimSuffix(javascript, ", ")

	javascript += ") {"

	javascript += "\n"

	var stepsJavascript string
	stepsJavascript, err = GetJavascriptBySteps(app, action.Steps, 1)
	if err != nil {
		return
	}
	if base.IsNotEmpty(stepsJavascript) {
		javascript += stepsJavascript
	}

	if action.OutVariable != nil {
		base.AppendLine(&javascript, "return "+action.OutVariable.Name, 1)
	}

	javascript += "}"
	// fmt.Println(javascript)
	return
}

func GetJavascriptBySteps(app IApplication, steps []model2.ActionStep, tab int) (javascript string, err error) {
	if len(steps) == 0 {
		return
	}
	for _, step := range steps {
		var stepJavascript string
		stepJavascript, err = GetJavascriptByStep(app, step, tab)
		if err != nil {
			return
		}
		if base.IsNotEmpty(stepJavascript) {
			javascript += stepJavascript
			javascript += "\n"
		}
	}
	return
}

func GetJavascriptByStep(app IApplication, step model2.ActionStep, tab int) (javascript string, err error) {
	if base.IsNotEmpty(step.GetBase().Comment) {
		base.AppendLine(&javascript, "// "+step.GetBase().Comment, tab)
	}
	var if_ string = step.GetBase().If
	if base.IsNotEmpty(if_) {
		base.AppendLine(&javascript, "if ("+if_+") { ", tab)
		tab++
	}
	var validatasJavascript string
	validatasJavascript, err = getJavascriptByValidatas(app, step.GetBase().Validates, tab)
	if err != nil {
		return
	}
	if base.IsNotEmpty(validatasJavascript) {
		javascript += validatasJavascript
	}
	var variablesJavascript string
	variablesJavascript, err = getJavascriptByVariables(app, step.GetBase().Variables, tab)
	if err != nil {
		return
	}
	if base.IsNotEmpty(variablesJavascript) {
		javascript += variablesJavascript
	}
	var stepJavascript string
	switch step_ := step.(type) {
	case *model2.ActionStepLock:
		stepJavascript, err = getJavascriptByStepLock(app, step_, tab)
	case *model2.ActionStepUnlock:
		stepJavascript, err = getJavascriptByStepUnlock(app, step_, tab)
	case *model2.ActionStepError:
		stepJavascript, err = getJavascriptByStepError(app, step_, tab)
	case *model2.ActionStepSqlSelect:
		stepJavascript, err = getJavascriptByStepSqlSelect(app, step_, tab)
	case *model2.ActionStepSqlInsert:
		stepJavascript, err = getJavascriptByStepSqlInsert(app, step_, tab)
	case *model2.ActionStepSqlUpdate:
		stepJavascript, err = getJavascriptByStepSqlUpdate(app, step_, tab)
	case *model2.ActionStepSqlDelete:
		stepJavascript, err = getJavascriptByStepSqlDelete(app, step_, tab)
	case *model2.ActionStepRedisSet:
		stepJavascript, err = getJavascriptByStepRedisSet(app, step_, tab)
	case *model2.ActionStepRedisGet:
		stepJavascript, err = getJavascriptByStepRedisGet(app, step_, tab)
	case *model2.ActionStepRedisDel:
		stepJavascript, err = getJavascriptByStepRedisDel(app, step_, tab)
	case *model2.ActionStepRedisExpire:
		stepJavascript, err = getJavascriptByStepRedisExpire(app, step_, tab)
	case *model2.ActionStepRedisExpireat:
		stepJavascript, err = getJavascriptByStepRedisExpireat(app, step_, tab)
	case *model2.ActionStepAction:
		stepJavascript, err = getJavascriptByStepAction(app, step_, tab)
	case *model2.ActionStepFileSave:
		stepJavascript, err = getJavascriptByStepFileSave(app, step_, tab)
	case *model2.ActionStepFileGet:
		stepJavascript, err = getJavascriptByStepFileGet(app, step_, tab)
	case *model2.ActionStepBase:
	default:
		err = errors.New(fmt.Sprint("GetJavascriptByStep step type not match:", reflect.TypeOf(step).Elem().Name()))
		return
	}
	if err != nil {
		return
	}
	if base.IsNotEmpty(stepJavascript) {
		javascript += stepJavascript
	}

	var stepsJavascript string
	stepsJavascript, err = GetJavascriptBySteps(app, step.GetBase().Steps, tab)
	if err != nil {
		return
	}
	if base.IsNotEmpty(stepsJavascript) {
		javascript += stepsJavascript
	}

	if step.GetBase().Return {
		if base.IsNotEmpty(step.GetBase().ReturnVariableName) {
			base.AppendLine(&javascript, "return "+step.GetBase().ReturnVariableName, tab)
		} else {
			base.AppendLine(&javascript, "return", tab)
		}
	}
	if base.IsNotEmpty(if_) {
		tab--
		base.AppendLine(&javascript, "} ", tab)
	}
	return
}

func getJavascriptByValidatas(app IApplication, validatas []*model2.ValidateModel, tab int) (javascript string, err error) {
	var errorModel *model2.ErrorModel
	var errorJavascript string
	for _, one := range validatas {
		if base.IsNotEmpty(one.Comment) {
			base.AppendLine(&javascript, "// "+one.Comment, tab)
		}
		errorModel, err = GetErrorModel(app, one.Error, one.ErrorCode, one.ErrorMsg)
		if err != nil {
			return
		}
		errorJavascript, err = getJavascriptByValidataRule(app, one.Name, &model2.ValidateRuleModel{
			Required:  one.Required,
			MinLength: one.MinLength,
			MaxLength: one.MaxLength,
			Min:       one.Min,
			Max:       one.Max,
			Pattern:   one.Pattern,
			Error:     one.Error,
			ErrorCode: one.ErrorCode,
			ErrorMsg:  one.ErrorMsg,
		}, errorModel, tab)
		if err != nil {
			return
		}
		if base.IsNotEmpty(errorJavascript) {
			javascript += errorJavascript
		}
		for _, rule := range one.Rules {
			errorJavascript, err = getJavascriptByValidataRule(app, one.Name, rule, errorModel, tab)
			if err != nil {
				return
			}
			if base.IsNotEmpty(errorJavascript) {
				javascript += errorJavascript
			}
		}
	}
	return
}

func getJavascriptByValidataRule(app IApplication, name string, rule *model2.ValidateRuleModel, parentErrorModel *model2.ErrorModel, tab int) (javascript string, err error) {
	var errorModel *model2.ErrorModel
	var errorJavascript string

	errorModel, err = GetErrorModel(app, rule.Error, rule.ErrorCode, rule.ErrorMsg)
	if err != nil {
		return
	}
	if parentErrorModel != nil {
		if base.IsEmpty(errorModel.Code) && base.IsEmpty(errorModel.Msg) {
			errorModel = parentErrorModel
		}
	}
	fieldComment := name
	if strings.Contains(name, ".") {
		// parentName := name[:strings.Index(name, ".")]
		// fieldName := name[strings.Index(name, ".")+1:]
		// for _, in := range action.InVariables {
		// 	if in.Name == parentName {
		// 		dataType := app.GetContext().GetVariableDataType(in.DataType)
		// 		if dataType.DataStruct != nil {
		// 			field := dataType.DataStruct.GetField(fieldName)
		// 			if field != nil {
		// 				fieldComment = field.Comment
		// 			}
		// 		}
		// 		break
		// 	}
		// }
	}
	if rule.Required {
		base.AppendLine(&javascript, "if (isEmpty("+name+")) { ", tab)
		errorJavascript = getJavascriptByErrorModel(app, errorModel, fieldComment+"不能为空", tab+1)
		javascript += errorJavascript
		base.AppendLine(&javascript, "} ", tab)
	}
	if rule.MinLength > 0 {
		base.AppendLine(&javascript, "if (length("+name+") < "+fmt.Sprint(rule.MinLength)+") { ", tab)
		errorJavascript = getJavascriptByErrorModel(app, errorModel, fieldComment+"长度不能小于 "+fmt.Sprint(rule.MinLength), tab+1)
		javascript += errorJavascript
		base.AppendLine(&javascript, "} ", tab)
	}
	if rule.MaxLength > 0 {
		base.AppendLine(&javascript, "if (length("+name+") > "+fmt.Sprint(rule.MaxLength)+") { ", tab)
		errorJavascript = getJavascriptByErrorModel(app, errorModel, fieldComment+"长度不能大于 "+fmt.Sprint(rule.MaxLength), tab+1)
		javascript += errorJavascript
		base.AppendLine(&javascript, "} ", tab)
	}
	if base.IsNotEmpty(rule.Pattern) {
		base.AppendLine(&javascript, "if (notMatch(`"+rule.Pattern+"`, "+name+")) { ", tab)
		errorJavascript = getJavascriptByErrorModel(app, errorModel, fieldComment+"格式不正确", tab+1)
		javascript += errorJavascript
		base.AppendLine(&javascript, "} ", tab)
	}
	return
}

func getJavascriptByVariables(app IApplication, variables []*model2.VariableModel, tab int) (javascript string, err error) {
	for _, one := range variables {
		if base.IsNotEmpty(one.Comment) {
			base.AppendLine(&javascript, "// "+one.Comment, tab)
		}
		base.AppendLine(&javascript, `addDataInfo("`+one.Name+`", "`+one.DataType+`", "`+one.Comment+`", "`+one.Value+`", `+fmt.Sprint(one.IsList)+`, `+fmt.Sprint(one.IsPage)+`)`, tab)
		if base.IsEmpty(one.Value) {
			dataType := app.GetContext().GetVariableDataType(one.DataType)
			valueStr := `""`
			if dataType != nil && dataType.DataStruct == nil {
				switch dataType {
				case model2.DATA_TYPE_LONG, model2.DATA_TYPE_INT, model2.DATA_TYPE_SHORT, model2.DATA_TYPE_BYTE:
					valueStr = `0`
				case model2.DATA_TYPE_BOOLEAN:
					valueStr = `false`
				case model2.DATA_TYPE_DOUBLE, model2.DATA_TYPE_FLOAT:
					valueStr = `0.0`
				case model2.DATA_TYPE_MAP:
					valueStr = `{}`
				}
				base.AppendLine(&javascript, one.Name+" = "+"{}", tab)
				base.AppendLine(&javascript, one.Name+" = "+valueStr, tab)
			} else {
				base.AppendLine(&javascript, `newVariable("`+one.Name+`", "`+one.DataType+`")`, tab)
			}
		} else {
			base.AppendLine(&javascript, one.Name+" = "+one.Value, tab)
		}
	}

	return
}

func getJavascriptByStepLock(app IApplication, step *model2.ActionStepLock, tab int) (javascript string, err error) {
	name := step.Lock.Name
	if base.IsEmpty(name) {
		name = "$lock_" + app.GetScript().RandString(10, 10)
	}
	key := step.Lock.Key
	if base.IsEmpty(key) {
		key = `""`
	}

	base.AppendLine(&javascript, "$invoke_temp.lock_"+name+" = getLock("+key+")", tab)
	base.AppendLine(&javascript, "$invoke_temp.lock_"+name+".lock()", tab)
	base.AppendLine(&javascript, "// 埋点，防止异常或为主动释放锁，将在执行结束释放该锁", tab)
	base.AppendLine(&javascript, "$invoke_temp.defer.$invoke_temp.lock_"+name+".unlock()", tab)
	return
}

func getJavascriptByStepUnlock(app IApplication, step *model2.ActionStepUnlock, tab int) (javascript string, err error) {
	name := step.Unlock.Name
	base.AppendLine(&javascript, "$invoke_temp.lock_"+name+".unlock()", tab)
	return
}

func getJavascriptByStepError(app IApplication, step *model2.ActionStepError, tab int) (javascript string, err error) {
	var errorModel *model2.ErrorModel
	errorModel, err = GetErrorModel(app, step.Error.Name, step.Error.Code, step.Error.Msg)
	if err != nil {
		return
	}
	var errorJavascript string = getJavascriptByErrorModel(app, errorModel, "", tab)
	javascript += errorJavascript
	return
}

func getJavascriptByErrorModel(app IApplication, errorModel *model2.ErrorModel, defaultMsg string, tab int) (javascript string) {
	if errorModel == nil {
		return
	}
	msg := errorModel.Msg
	if base.IsEmpty(msg) {
		msg = defaultMsg
	}
	base.AppendLine(&javascript, "throwError(\""+errorModel.Code+"\", \""+msg+"\")", tab)
	return
}

func getJavascriptByStepSqlSelect(app IApplication, step *model2.ActionStepSqlSelect, tab int) (javascript string, err error) {

	var javascript_ string
	javascript_, err = getJavascriptBySqlSelect(app, step.SqlSelect, tab)
	if err != nil {
		return
	}
	javascript += javascript_
	javascript += "\n"
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	variableDataType := step.VariableDataType
	if base.IsNotEmpty(step.VariableName) {
		if step.SqlSelect.SelectCount {
			if base.IsEmpty(step.VariableDataType) {
				base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "long", "", "", false)`, tab)
			} else {
				base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`", "", "", false)`, tab)
			}
		} else if step.SqlSelect.SelectOne {
			base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`", "", "", false)`, tab)
		} else if step.SqlSelect.SelectPage {
			base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`", "", "", false, true)`, tab)
		} else {
			base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`", "", "", true)`, tab)
		}
	}
	if step.SqlSelect.SelectCount {
		base.AppendLine(&javascript, "// 执行SQL 统计查询", tab)
		base.AppendLine(&javascript, variableName+`_.sqlSelectCount($invoke_temp.countSql, $invoke_temp.countParams)`, tab)
	} else if step.SqlSelect.SelectPage {
		base.AppendLine(&javascript, "// 执行SQL 分页查询", tab)
		base.AppendLine(&javascript, `$invoke_temp.pageNumber = 1`, tab)
		base.AppendLine(&javascript, `$invoke_temp.pageSize = 1`, tab)
		base.AppendLine(&javascript, `$invoke_temp.totalPage = 0`, tab)
		base.AppendLine(&javascript, `$invoke_temp.totalSize = 0`, tab)
		base.AppendLine(&javascript, `$invoke_temp.list = []`, tab)
		base.AppendLine(&javascript, `$invoke_temp.totalSize = _.sqlSelectCount($invoke_temp.countSql, $invoke_temp.countParams)`, tab)
		base.AppendLine(&javascript, `if ($invoke_temp.pageSize > 0) {`, tab)
		base.AppendLine(&javascript, `$invoke_temp.totalPage = ($invoke_temp.totalSize + $invoke_temp.pageSize - 1) / $invoke_temp.pageSize`, tab+1)
		base.AppendLine(&javascript, `}`, tab)
		base.AppendLine(&javascript, `if ($invoke_temp.pageNumber > 0 && $invoke_temp.pageSize > 0 && $invoke_temp.pageSize <= $invoke_temp.totalPage) {`, tab)
		base.AppendLine(&javascript, `$invoke_temp.list = _.sqlSelectPage($invoke_temp.sql, $invoke_temp.params, $invoke_temp.pageNumber, $invoke_temp.pageSize, "`+variableDataType+`")`, tab)
		base.AppendLine(&javascript, `}`, tab)

		if base.IsNotEmpty(variableName) {
			base.AppendLine(&javascript, `newVariable("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
			base.AppendLine(&javascript, step.VariableName+`.pageNumber = $invoke_temp.pageNumber`, tab)
			base.AppendLine(&javascript, step.VariableName+`.pageSize = $invoke_temp.pageSize`, tab)
			base.AppendLine(&javascript, step.VariableName+`.totalPage = $invoke_temp.totalPage`, tab)
			base.AppendLine(&javascript, step.VariableName+`.totalSize = $invoke_temp.totalSize`, tab)
			base.AppendLine(&javascript, step.VariableName+`.list = $invoke_temp.list`, tab)
		}
	} else if step.SqlSelect.SelectOne {
		base.AppendLine(&javascript, "// 执行SQL 查询单个", tab)
		base.AppendLine(&javascript, variableName+`_.sqlSelectOne($invoke_temp.sql, $invoke_temp.params, "`+variableDataType+`")`, tab)
	} else {
		base.AppendLine(&javascript, "// 执行SQL 查询列表", tab)
		base.AppendLine(&javascript, variableName+`_.sqlSelect($invoke_temp.sql, $invoke_temp.params, "`+variableDataType+`")`, tab)
	}
	return
}

func getJavascriptByStepSqlInsert(app IApplication, step *model2.ActionStepSqlInsert, tab int) (javascript string, err error) {
	var javascript_ string
	javascript_, err = getJavascriptBySqlInsert(app, step.SqlInsert, tab)
	if err != nil {
		return
	}
	javascript += javascript_
	javascript += "\n"
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}
	base.AppendLine(&javascript, "// 执行SQL 新增", tab)
	base.AppendLine(&javascript, variableName+"_.sqlInsert($invoke_temp.sql, $invoke_temp.params)", tab)
	return
}

func getJavascriptByStepSqlUpdate(app IApplication, step *model2.ActionStepSqlUpdate, tab int) (javascript string, err error) {
	var javascript_ string
	javascript_, err = getJavascriptBySqlUpdate(app, step.SqlUpdate, tab)
	if err != nil {
		return
	}
	javascript += javascript_
	javascript += "\n"
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}
	base.AppendLine(&javascript, "// 执行SQL 更新", tab)
	base.AppendLine(&javascript, variableName+"_.sqlUpdate($invoke_temp.sql, $invoke_temp.params)", tab)
	return
}

func getJavascriptByStepSqlDelete(app IApplication, step *model2.ActionStepSqlDelete, tab int) (javascript string, err error) {
	var javascript_ string
	javascript_, err = getJavascriptBySqlDelete(app, step.SqlDelete, tab)
	if err != nil {
		return
	}
	javascript += javascript_
	javascript += "\n"
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}
	base.AppendLine(&javascript, "// 执行SQL 删除", tab)
	base.AppendLine(&javascript, variableName+"_.sqlDelete($invoke_temp.sql, $invoke_temp.params)", tab)
	return
}

func getJavascriptByStepRedisSet(app IApplication, step *model2.ActionStepRedisSet, tab int) (javascript string, err error) {
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}
	base.AppendLine(&javascript, "// 执行Redis 设置缓存", tab)
	base.AppendLine(&javascript, variableName+"_.redisSet("+step.RedisSet.Key+", "+step.RedisSet.Value+")", tab)
	return
}

func getJavascriptByStepRedisGet(app IApplication, step *model2.ActionStepRedisGet, tab int) (javascript string, err error) {
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}
	base.AppendLine(&javascript, "// 执行Redis 获取缓存", tab)
	base.AppendLine(&javascript, variableName+"_.redisGet("+step.RedisGet.Key+")", tab)
	return
}

func getJavascriptByStepRedisDel(app IApplication, step *model2.ActionStepRedisDel, tab int) (javascript string, err error) {
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}
	base.AppendLine(&javascript, "// 执行Redis 删除缓存", tab)
	base.AppendLine(&javascript, variableName+"_.redisDel("+step.RedisDel.Key+")", tab)
	return
}

func getJavascriptByStepRedisExpire(app IApplication, step *model2.ActionStepRedisExpire, tab int) (javascript string, err error) {
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}
	base.AppendLine(&javascript, variableName+"_.redisExpire()", tab)
	return
}

func getJavascriptByStepRedisExpireat(app IApplication, step *model2.ActionStepRedisExpireat, tab int) (javascript string, err error) {
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}
	base.AppendLine(&javascript, variableName+"_.redisExpireat()", tab)
	return
}

func getJavascriptByStepAction(app IApplication, step *model2.ActionStepAction, tab int) (javascript string, err error) {

	var variablesJavascript string
	variablesJavascript, err = getJavascriptByVariables(app, step.Action.CallVariables, tab)
	if err != nil {
		return
	}
	if base.IsNotEmpty(variablesJavascript) {
		javascript += variablesJavascript
	}
	callArgs := `"` + step.Action.Name + `", `
	for _, one := range step.Action.CallVariables {
		callArgs += one.Name + ", "
	}
	callArgs = strings.TrimSuffix(callArgs, ", ")

	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}

	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}
	base.AppendLine(&javascript, "// 调用服务方法", tab)
	base.AppendLine(&javascript, variableName+"action("+callArgs+")", tab)
	return
}

func getJavascriptByStepFileSave(app IApplication, step *model2.ActionStepFileSave, tab int) (javascript string, err error) {
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}
	callArgs := ""
	if base.IsNotEmpty(step.FileSave.Name) {
		callArgs += `` + step.FileSave.Name + `, `
	} else {
		callArgs += `null, `
	}
	if base.IsNotEmpty(step.FileSave.Dir) {
		callArgs += `` + step.FileSave.Dir + `, `
	} else {
		callArgs += `null, `
	}
	if base.IsNotEmpty(step.FileSave.Reader) {
		callArgs += `` + step.FileSave.Reader + `, `
	} else {
		callArgs += `null, `
	}
	if base.IsNotEmpty(step.FileSave.Bytes) {
		callArgs += `` + step.FileSave.Bytes + `, `
	} else {
		callArgs += `null, `
	}
	callArgs = strings.TrimSuffix(callArgs, ", ")

	base.AppendLine(&javascript, "// 文件保存", tab)
	base.AppendLine(&javascript, variableName+"fileSave("+callArgs+")", tab)
	return
}

func getJavascriptByStepFileGet(app IApplication, step *model2.ActionStepFileGet, tab int) (javascript string, err error) {
	variableName := step.VariableName
	if base.IsNotEmpty(variableName) {
		variableName += " = "
	}
	if base.IsNotEmpty(step.VariableName) {
		base.AppendLine(&javascript, `addDataInfo("`+step.VariableName+`", "`+step.VariableDataType+`")`, tab)
	}

	callArgs := ""
	if base.IsNotEmpty(step.FileGet.Path) {
		callArgs += `` + step.FileGet.Path + `, `
	} else {
		callArgs += `null, `
	}
	callArgs = strings.TrimSuffix(callArgs, ", ")

	base.AppendLine(&javascript, "// 文件获取", tab)
	base.AppendLine(&javascript, variableName+"fileGet("+callArgs+")", tab)
	return
}
