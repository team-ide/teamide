package javascript

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"teamide/pkg/maker/model"
	"teamide/pkg/util"
)

func GetJavascriptBySteps(app *model.Application, steps []interface{}, tab int) (javascript string, err error) {
	if len(steps) == 0 {
		return
	}
	for _, step := range steps {
		var stepJavascript string
		stepJavascript, err = GetJavascriptByStep(app, step, tab)
		if err != nil {
			util.Logger.Error("GetJavascriptBySteps error", zap.Error(err))
			return
		}
		if util.IsNotEmpty(stepJavascript) {
			javascript += stepJavascript
			javascript += "\n"
		}
	}
	return
}

func GetJavascriptByStep(app *model.Application, step interface{}, tab int) (javascript string, err error) {
	if step == nil {
		err = errors.New("GetJavascriptByStep step is null")
		return
	}
	var stepJavascript string
	var stepModel *model.StepModel
	var hasIf bool
	switch step_ := step.(type) {
	case *model.StepModel:
		stepModel = step_
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

	case *model.StepErrorModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepError(app, step_, tab)
	case *model.StepVarModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepVar(app, step_, tab)
	case *model.StepLockModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepLock(app, step_, tab)
	case *model.StepDbModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepDb(app, step_, tab)
	case *model.StepRedisModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepRedis(app, step_, tab)
	case *model.StepEsModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepEs(app, step_, tab)
	case *model.StepZkModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepZk(app, step_, tab)
	case *model.StepFileModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepFile(app, step_, tab)
	case *model.StepCommandModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepCommand(app, step_, tab)
	case *model.StepCacheModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepCache(app, step_, tab)
	case *model.StepMqModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepMq(app, step_, tab)
	case *model.StepHttpModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepHttp(app, step_, tab)
	case *model.StepServiceModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepService(app, step_, tab)
	case *model.StepDaoModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepDao(app, step_, tab)
	case *model.StepScriptModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepScript(app, step_, tab)
	default:
		err = errors.New("GetJavascriptByStep step type [" + util.GetRefType(step).String() + "] is not support")
		util.Logger.Error("GetJavascriptByStep error", zap.Any("step", step), zap.Error(err))
		return
	}
	if err != nil {
		util.Logger.Error("GetJavascriptByStep error", zap.Any("step", step), zap.Error(err))
		return
	}
	if util.IsNotEmpty(stepJavascript) {
		javascript += stepJavascript
	}

	var stepsJavascript string
	stepsJavascript, err = GetJavascriptBySteps(app, stepModel.Steps, tab)
	if err != nil {
		return
	}
	if util.IsNotEmpty(stepsJavascript) {
		javascript += stepsJavascript
	}

	if util.IsNotEmpty(stepModel.Return) {
		if stepModel.Return != "-" {
			util.AppendLine(&javascript, "return "+stepModel.Return, tab)
		} else {
			util.AppendLine(&javascript, "return", tab)
		}
	}
	if hasIf {
		tab--
		util.AppendLine(&javascript, "} ", tab)
	}
	return
}

func getJavascriptByStep(app *model.Application, step *model.StepModel, tab *int) (javascript string, hasIf bool) {

	if util.IsNotEmpty(step.Note) {
		util.AppendLine(&javascript, "// "+step.Note, *tab)
	} else if util.IsNotEmpty(step.Comment) {
		util.AppendLine(&javascript, "// "+step.Comment, *tab)
	}
	if util.IsNotEmpty(step.If) {
		hasIf = true
		util.AppendLine(&javascript, "if ("+step.If+") { ", *tab)
		*tab++
	}
	return
}

func getJavascriptByStepError(app *model.Application, step *model.StepErrorModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// 异常 操作", tab)

	errorModel := app.GetError(step.Error)
	if errorModel == nil {
		util.AppendLine(&javascript, "throw new Error(\""+step.Error+"\")", tab)
	} else {
		util.AppendLine(&javascript, "throw app.Errors."+errorModel.Name+"", tab)
	}
	return
}

func getJavascriptByStepVar(app *model.Application, step *model.StepVarModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// 定义变量 "+step.Var+" ", tab)

	return
}

func getJavascriptByStepLock(app *model.Application, step *model.StepLockModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// 锁 "+step.Lock+" 操作", tab)

	return
}

func getJavascriptByStepDb(app *model.Application, step *model.StepDbModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// DB "+step.Db+" 操作", tab)

	util.AppendLine(&javascript, "sql = \"\"", tab)
	util.AppendLine(&javascript, "sqlParams = []", tab)

	var isSelect bool
	var isInsert bool
	var isUpdate bool
	var isDelete bool
	switch step.Db {
	case "select":
		isSelect = true
		util.AppendLine(&javascript, fmt.Sprintf(`sql += "SELECT %s"`, step.Table), tab)
		break
	case "selectOne":
		isSelect = true
		util.AppendLine(&javascript, fmt.Sprintf(`sql += "SELECT %s"`, step.Table), tab)
		break
	case "selectPage":
		isSelect = true
		util.AppendLine(&javascript, fmt.Sprintf(`sql += "SELECT %s"`, step.Table), tab)
		break
	case "insert":
		isInsert = true
		util.AppendLine(&javascript, fmt.Sprintf(`sql += "INSERT %s"`, step.Table), tab)
		break
	case "update":
		isUpdate = true
		util.AppendLine(&javascript, fmt.Sprintf(`sql += "UPDATE %s"`, step.Table), tab)
		break
	case "delete":
		isDelete = true
		util.AppendLine(&javascript, fmt.Sprintf(`sql += "DELETE %s"`, step.Table), tab)
		break
	}
	if isSelect {
		util.AppendLine(&javascript, fmt.Sprintf(`%s = dbHandler.select(sql, sqlParams)`, step.SetVar), tab)
	} else if isInsert {
		util.AppendLine(&javascript, fmt.Sprintf(`%s = dbHandler.insert(sql, sqlParams)`, step.SetVar), tab)
	} else if isUpdate {
		util.AppendLine(&javascript, fmt.Sprintf(`%s = dbHandler.update(sql, sqlParams)`, step.SetVar), tab)
	} else if isDelete {
		util.AppendLine(&javascript, fmt.Sprintf(`%s = dbHandler.delete(sql, sqlParams)`, step.SetVar), tab)
	}
	return
}

func getJavascriptByStepRedis(app *model.Application, step *model.StepRedisModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// Redis "+step.Redis+" 操作", tab)
	switch strings.ToLower(step.Redis) {
	case "get":
		util.AppendLine(&javascript, fmt.Sprintf(`%s = redisHandler.get("%s")`, step.SetVar, step.Key), tab)
		break
	case "set":
		util.AppendLine(&javascript, fmt.Sprintf(`redisHandler.set("%s", "%s")`, step.Key, step.Value), tab)
		break
	}
	return
}

func getJavascriptByStepEs(app *model.Application, step *model.StepEsModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// ES "+step.Es+" 操作", tab)

	return
}

func getJavascriptByStepZk(app *model.Application, step *model.StepZkModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// ZK "+step.Zk+" 操作", tab)

	return
}

func getJavascriptByStepCache(app *model.Application, step *model.StepCacheModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// Cache "+step.Cache+" 操作", tab)

	return
}

func getJavascriptByStepFile(app *model.Application, step *model.StepFileModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// File "+step.File+" 操作", tab)

	return
}

func getJavascriptByStepCommand(app *model.Application, step *model.StepCommandModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// Command "+step.Command+" 操作", tab)

	return
}

func getJavascriptByStepService(app *model.Application, step *model.StepServiceModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// Service "+step.Service+" 操作", tab)

	return
}

func getJavascriptByStepDao(app *model.Application, step *model.StepDaoModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// Dao "+step.Dao+" 操作", tab)

	return
}

func getJavascriptByStepScript(app *model.Application, step *model.StepScriptModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// Script "+step.Script+" 操作", tab)

	return
}

func getJavascriptByStepMq(app *model.Application, step *model.StepMqModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// MQ "+step.Mq+" 操作", tab)

	return
}

func getJavascriptByStepHttp(app *model.Application, step *model.StepHttpModel, tab int) (javascript string, err error) {
	util.AppendLine(&javascript, "// Http "+step.Http+" 操作", tab)

	return
}
