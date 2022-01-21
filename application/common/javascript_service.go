package common

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"teamide/application/base"
	"teamide/application/model"
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

func GetServiceJavascriptByService(app IApplication, service *model.ServiceModel) (javascript string, err error) {
	methodName := GetJavascriptMethodName(service.Name)
	javascript += ""
	javascript += "function service_" + methodName + "("
	for _, inVariable := range service.InVariables {
		javascript += inVariable.Name + ", "
	}
	javascript = strings.TrimSuffix(javascript, ", ")

	javascript += ") {"

	javascript += "\n"

	var stepsJavascript string
	stepsJavascript, err = GetJavascriptBySteps(app, service.Steps, 1)
	if err != nil {
		return
	}
	if base.IsNotEmpty(stepsJavascript) {
		javascript += stepsJavascript
	}

	if service.OutVariable != nil {
		base.AppendLine(&javascript, "return "+service.OutVariable.Name, 1)
	}

	javascript += "}"
	// fmt.Println(javascript)
	return
}

func GetJavascriptBySteps(app IApplication, steps []model.ServiceStep, tab int) (javascript string, err error) {
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

func GetJavascriptByStep(app IApplication, step model.ServiceStep, tab int) (javascript string, err error) {
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
	case *model.ServiceStepLock:
		stepJavascript, err = getJavascriptByStepLock(app, step_, tab)
	case *model.ServiceStepUnlock:
		stepJavascript, err = getJavascriptByStepUnlock(app, step_, tab)
	case *model.ServiceStepError:
		stepJavascript, err = getJavascriptByStepError(app, step_, tab)
	case *model.ServiceStepSqlSelect:
		stepJavascript, err = getJavascriptByStepSqlSelect(app, step_, tab)
	case *model.ServiceStepSqlInsert:
		stepJavascript, err = getJavascriptByStepSqlInsert(app, step_, tab)
	case *model.ServiceStepSqlUpdate:
		stepJavascript, err = getJavascriptByStepSqlUpdate(app, step_, tab)
	case *model.ServiceStepSqlDelete:
		stepJavascript, err = getJavascriptByStepSqlDelete(app, step_, tab)
	case *model.ServiceStepRedisSet:
		stepJavascript, err = getJavascriptByStepRedisSet(app, step_, tab)
	case *model.ServiceStepRedisGet:
		stepJavascript, err = getJavascriptByStepRedisGet(app, step_, tab)
	case *model.ServiceStepRedisDel:
		stepJavascript, err = getJavascriptByStepRedisDel(app, step_, tab)
	case *model.ServiceStepRedisExpire:
		stepJavascript, err = getJavascriptByStepRedisExpire(app, step_, tab)
	case *model.ServiceStepRedisExpireat:
		stepJavascript, err = getJavascriptByStepRedisExpireat(app, step_, tab)
	case *model.ServiceStepService:
		stepJavascript, err = getJavascriptByStepService(app, step_, tab)
	case *model.ServiceStepFileSave:
		stepJavascript, err = getJavascriptByStepFileSave(app, step_, tab)
	case *model.ServiceStepFileGet:
		stepJavascript, err = getJavascriptByStepFileGet(app, step_, tab)
	case *model.ServiceStepBase:
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

func getJavascriptByValidatas(app IApplication, validatas []*model.ValidateModel, tab int) (javascript string, err error) {
	var errorModel *model.ErrorModel
	var errorJavascript string
	for _, one := range validatas {
		if base.IsNotEmpty(one.Comment) {
			base.AppendLine(&javascript, "// "+one.Comment, tab)
		}
		errorModel, err = GetErrorModel(app, one.Error, one.ErrorCode, one.ErrorMsg)
		if err != nil {
			return
		}
		errorJavascript, err = getJavascriptByValidataRule(app, one.Name, &model.ValidateRuleModel{
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

func getJavascriptByValidataRule(app IApplication, name string, rule *model.ValidateRuleModel, parentErrorModel *model.ErrorModel, tab int) (javascript string, err error) {
	var errorModel *model.ErrorModel
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
		// for _, in := range service.InVariables {
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

func getJavascriptByVariables(app IApplication, variables []*model.VariableModel, tab int) (javascript string, err error) {
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
				case model.DATA_TYPE_LONG, model.DATA_TYPE_INT, model.DATA_TYPE_SHORT, model.DATA_TYPE_BYTE:
					valueStr = `0`
				case model.DATA_TYPE_BOOLEAN:
					valueStr = `false`
				case model.DATA_TYPE_DOUBLE, model.DATA_TYPE_FLOAT:
					valueStr = `0.0`
				case model.DATA_TYPE_MAP:
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

func getJavascriptByStepLock(app IApplication, step *model.ServiceStepLock, tab int) (javascript string, err error) {
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

func getJavascriptByStepUnlock(app IApplication, step *model.ServiceStepUnlock, tab int) (javascript string, err error) {
	name := step.Unlock.Name
	base.AppendLine(&javascript, "$invoke_temp.lock_"+name+".unlock()", tab)
	return
}

func getJavascriptByStepError(app IApplication, step *model.ServiceStepError, tab int) (javascript string, err error) {
	var errorModel *model.ErrorModel
	errorModel, err = GetErrorModel(app, step.Error.Name, step.Error.Code, step.Error.Msg)
	if err != nil {
		return
	}
	var errorJavascript string = getJavascriptByErrorModel(app, errorModel, "", tab)
	javascript += errorJavascript
	return
}

func getJavascriptByErrorModel(app IApplication, errorModel *model.ErrorModel, defaultMsg string, tab int) (javascript string) {
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

func getJavascriptByStepSqlSelect(app IApplication, step *model.ServiceStepSqlSelect, tab int) (javascript string, err error) {

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

func getJavascriptByStepSqlInsert(app IApplication, step *model.ServiceStepSqlInsert, tab int) (javascript string, err error) {
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

func getJavascriptByStepSqlUpdate(app IApplication, step *model.ServiceStepSqlUpdate, tab int) (javascript string, err error) {
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

func getJavascriptByStepSqlDelete(app IApplication, step *model.ServiceStepSqlDelete, tab int) (javascript string, err error) {
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

func getJavascriptByStepRedisSet(app IApplication, step *model.ServiceStepRedisSet, tab int) (javascript string, err error) {
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

func getJavascriptByStepRedisGet(app IApplication, step *model.ServiceStepRedisGet, tab int) (javascript string, err error) {
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

func getJavascriptByStepRedisDel(app IApplication, step *model.ServiceStepRedisDel, tab int) (javascript string, err error) {
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

func getJavascriptByStepRedisExpire(app IApplication, step *model.ServiceStepRedisExpire, tab int) (javascript string, err error) {
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

func getJavascriptByStepRedisExpireat(app IApplication, step *model.ServiceStepRedisExpireat, tab int) (javascript string, err error) {
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

func getJavascriptByStepService(app IApplication, step *model.ServiceStepService, tab int) (javascript string, err error) {

	var variablesJavascript string
	variablesJavascript, err = getJavascriptByVariables(app, step.Service.CallVariables, tab)
	if err != nil {
		return
	}
	if base.IsNotEmpty(variablesJavascript) {
		javascript += variablesJavascript
	}
	callArgs := `"` + step.Service.Name + `", `
	for _, one := range step.Service.CallVariables {
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
	base.AppendLine(&javascript, variableName+"service("+callArgs+")", tab)
	return
}

func getJavascriptByStepFileSave(app IApplication, step *model.ServiceStepFileSave, tab int) (javascript string, err error) {
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

func getJavascriptByStepFileGet(app IApplication, step *model.ServiceStepFileGet, tab int) (javascript string, err error) {
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
