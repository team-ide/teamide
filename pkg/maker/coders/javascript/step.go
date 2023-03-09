package javascript

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"strings"
	"teamide/pkg/base"
	"teamide/pkg/maker/coders/common"
	"teamide/pkg/maker/modelers"
)

type stepCoder struct {
	*appCoder
}

func (this_ *stepCoder) Gen(code *common.Code, model *modelers.StepModel) (err error) {

	return
}

func GetJavascriptBySteps(app *modelers.Application, steps []interface{}, tab int) (javascript string, err error) {
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

func GetJavascriptByStep(app *modelers.Application, step interface{}, tab int) (javascript string, err error) {
	if step == nil {
		err = errors.New("GetJavascriptByStep step is null")
		return
	}
	var stepJavascript string
	var stepModel *modelers.StepModel
	var hasIf bool
	switch step_ := step.(type) {
	case *modelers.StepModel:
		stepModel = step_
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

	case *modelers.StepErrorModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepError(app, step_, tab)
	case *modelers.StepVarModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepVar(app, step_, tab)
	case *modelers.StepLockModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepLock(app, step_, tab)
	case *modelers.StepDbModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepDb(app, step_, tab)
	case *modelers.StepRedisModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepRedis(app, step_, tab)
	case *modelers.StepEsModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepEs(app, step_, tab)
	case *modelers.StepZkModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepZk(app, step_, tab)
	case *modelers.StepFileModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepFile(app, step_, tab)
	case *modelers.StepCommandModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepCommand(app, step_, tab)
	case *modelers.StepCacheModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepCache(app, step_, tab)
	case *modelers.StepMqModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepMq(app, step_, tab)
	case *modelers.StepHttpModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepHttp(app, step_, tab)
	case *modelers.StepServiceModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepService(app, step_, tab)
	case *modelers.StepDaoModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepDao(app, step_, tab)
	case *modelers.StepScriptModel:
		stepModel = step_.StepModel
		javascript, hasIf = getJavascriptByStep(app, stepModel, &tab)

		stepJavascript, err = getJavascriptByStepScript(app, step_, tab)
	default:
		err = errors.New("GetJavascriptByStep step type [" + base.GetRefType(step).String() + "] is not support")
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
			base.AppendLine(&javascript, "return "+stepModel.Return, tab)
		} else {
			base.AppendLine(&javascript, "return", tab)
		}
	}
	if hasIf {
		tab--
		base.AppendLine(&javascript, "} ", tab)
	}
	return
}

func getJavascriptByStep(app *modelers.Application, step *modelers.StepModel, tab *int) (javascript string, hasIf bool) {

	if util.IsNotEmpty(step.Note) {
		base.AppendLine(&javascript, "// "+step.Note, *tab)
	} else if util.IsNotEmpty(step.Comment) {
		base.AppendLine(&javascript, "// "+step.Comment, *tab)
	}
	if util.IsNotEmpty(step.If) {
		hasIf = true
		base.AppendLine(&javascript, "if ("+step.If+") { ", *tab)
		*tab++
	}
	return
}

func getJavascriptByStepError(app *modelers.Application, step *modelers.StepErrorModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// 异常 操作", tab)

	errorModel := app.GetError(step.Error)
	if errorModel == nil {
		base.AppendLine(&javascript, "throw new Error(\""+step.Error+"\")", tab)
	} else {
		base.AppendLine(&javascript, "throw app.Errors."+errorModel.Name+"", tab)
	}
	return
}

func getJavascriptByStepVar(app *modelers.Application, step *modelers.StepVarModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// 定义变量 "+step.Var+" ", tab)

	return
}

func getJavascriptByStepLock(app *modelers.Application, step *modelers.StepLockModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// 锁 "+step.Lock+" 操作", tab)

	return
}

func getJavascriptByStepDb(app *modelers.Application, step *modelers.StepDbModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// DB "+step.Db+" 操作", tab)

	base.AppendLine(&javascript, "sql = \"\"", tab)
	base.AppendLine(&javascript, "sqlParams = []", tab)

	var isSelect bool
	var isInsert bool
	var isUpdate bool
	var isDelete bool
	switch step.Db {
	case "select":
		isSelect = true
		base.AppendLine(&javascript, fmt.Sprintf(`sql += "SELECT %s"`, step.Table), tab)
		break
	case "selectOne":
		isSelect = true
		base.AppendLine(&javascript, fmt.Sprintf(`sql += "SELECT %s"`, step.Table), tab)
		break
	case "selectPage":
		isSelect = true
		base.AppendLine(&javascript, fmt.Sprintf(`sql += "SELECT %s"`, step.Table), tab)
		break
	case "insert":
		isInsert = true
		base.AppendLine(&javascript, fmt.Sprintf(`sql += "INSERT %s"`, step.Table), tab)
		break
	case "update":
		isUpdate = true
		base.AppendLine(&javascript, fmt.Sprintf(`sql += "UPDATE %s"`, step.Table), tab)
		break
	case "delete":
		isDelete = true
		base.AppendLine(&javascript, fmt.Sprintf(`sql += "DELETE %s"`, step.Table), tab)
		break
	}
	if isSelect {
		base.AppendLine(&javascript, fmt.Sprintf(`%s = dbHandler.select(sql, sqlParams)`, step.SetVar), tab)
	} else if isInsert {
		base.AppendLine(&javascript, fmt.Sprintf(`%s = dbHandler.insert(sql, sqlParams)`, step.SetVar), tab)
	} else if isUpdate {
		base.AppendLine(&javascript, fmt.Sprintf(`%s = dbHandler.update(sql, sqlParams)`, step.SetVar), tab)
	} else if isDelete {
		base.AppendLine(&javascript, fmt.Sprintf(`%s = dbHandler.delete(sql, sqlParams)`, step.SetVar), tab)
	}
	return
}

func getJavascriptByStepRedis(app *modelers.Application, step *modelers.StepRedisModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// Redis "+step.Redis+" 操作", tab)
	switch strings.ToLower(step.Redis) {
	case "get":
		base.AppendLine(&javascript, fmt.Sprintf(`%s = redisHandler.get("%s")`, step.SetVar, step.Key), tab)
		break
	case "set":
		base.AppendLine(&javascript, fmt.Sprintf(`redisHandler.set("%s", "%s")`, step.Key, step.Value), tab)
		break
	}
	return
}

func getJavascriptByStepEs(app *modelers.Application, step *modelers.StepEsModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// ES "+step.Es+" 操作", tab)

	return
}

func getJavascriptByStepZk(app *modelers.Application, step *modelers.StepZkModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// ZK "+step.Zk+" 操作", tab)

	return
}

func getJavascriptByStepCache(app *modelers.Application, step *modelers.StepCacheModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// Cache "+step.Cache+" 操作", tab)

	return
}

func getJavascriptByStepFile(app *modelers.Application, step *modelers.StepFileModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// File "+step.File+" 操作", tab)

	return
}

func getJavascriptByStepCommand(app *modelers.Application, step *modelers.StepCommandModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// Command "+step.Command+" 操作", tab)

	return
}

func getJavascriptByStepService(app *modelers.Application, step *modelers.StepServiceModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// Service "+step.Service+" 操作", tab)

	return
}

func getJavascriptByStepDao(app *modelers.Application, step *modelers.StepDaoModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// Dao "+step.Dao+" 操作", tab)

	return
}

func getJavascriptByStepScript(app *modelers.Application, step *modelers.StepScriptModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// Script "+step.Script+" 操作", tab)

	return
}

func getJavascriptByStepMq(app *modelers.Application, step *modelers.StepMqModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// MQ "+step.Mq+" 操作", tab)

	return
}

func getJavascriptByStepHttp(app *modelers.Application, step *modelers.StepHttpModel, tab int) (javascript string, err error) {
	base.AppendLine(&javascript, "// Http "+step.Http+" 操作", tab)

	return
}
