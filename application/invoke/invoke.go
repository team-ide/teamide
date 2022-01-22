package invoke

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"teamide/application/base"
	"teamide/application/common"
	"teamide/application/model"
)

type InvokeInfo struct {
	App             common.IApplication     `json:"-"`
	InvokeNamespace *common.InvokeNamespace `json:"-"`
}

var (
	invokeCallMap map[string]func(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error)
)

func init() {
	invokeCallMap = make(map[string]func(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error))
	invokeCallMap[".lock"] = invokeLock
	invokeCallMap[".unlock"] = invokeUnlock
	invokeCallMap[".push"] = invokePush
	invokeCallMap[".pushs"] = invokePushs
	invokeCallMap[".sqlSelect"] = invokeSqlSelect
	invokeCallMap[".sqlSelectOne"] = invokeSqlSelectOne
	invokeCallMap[".sqlSelectCount"] = invokeSqlSelectCount
	invokeCallMap[".sqlSelectPage"] = invokeSqlSelectPage
	invokeCallMap[".sqlInsert"] = invokeSqlInsert
	invokeCallMap[".sqlUpdate"] = invokeSqlUpdate
	invokeCallMap[".sqlDelete"] = invokeSqlDelete
	invokeCallMap[".redisSet"] = invokeRedisSet
	invokeCallMap[".redisGet"] = invokeRedisGet
	invokeCallMap[".redisDel"] = invokeRedisDel
	invokeCallMap["service"] = invokeService
	invokeCallMap["newVariable"] = invokeNewVariable
	invokeCallMap["addDataInfo"] = invokeAddDataInfo
	invokeCallMap["validataRequestToken"] = invokeValidataRequestToken
	invokeCallMap["getRequestData"] = invokeGetRequestData
	invokeCallMap["fileSave"] = invokeFileSave
	invokeCallMap["fileGet"] = invokeFileGet
}

func invokeValidataRequestToken(invokeInfo *InvokeInfo, _ string, args []interface{}) (res interface{}, err error) {
	if invokeInfo.InvokeNamespace.RequestContext == nil {
		err = base.NewError("", "request context not defind")
		return
	}
	if invokeInfo.InvokeNamespace.ServerWebToken == nil {
		return
	}
	terverWebToken := invokeInfo.InvokeNamespace.ServerWebToken

	res, err = invokeWebApiValidateToken(invokeInfo.App, terverWebToken, invokeInfo.InvokeNamespace.RequestContext)
	if err != nil {
		return
	}
	if terverWebToken.VariableName != "" {
		err = invokeInfo.InvokeNamespace.SetData(terverWebToken.VariableName, res, nil)
		if err != nil {
			return
		}
	}
	return
}

func invokeGetRequestData(invokeInfo *InvokeInfo, _ string, args []interface{}) (res interface{}, err error) {
	if invokeInfo.InvokeNamespace.RequestContext == nil {
		err = base.NewError("", "request context not defind")
		return
	}
	dataPlace := args[0].(string)
	name := args[1].(string)
	valueScript := args[2].(string)
	var value interface{}
	switch model.GetVariableDataPlace(dataPlace) {
	case model.DATA_PLACE_PATH:
		var pathValue string
		if base.IsNotEmpty(valueScript) {
			pathValue = invokeInfo.InvokeNamespace.RequestContext.Param(valueScript)
		} else {
			pathValue = invokeInfo.InvokeNamespace.RequestContext.Param(name)
		}
		if strings.Index(pathValue, "/") == 0 {
			pathValue = pathValue[1:]
		}
		value = pathValue
	case model.DATA_PLACE_HEADER:
		if base.IsNotEmpty(valueScript) {
			value = invokeInfo.InvokeNamespace.RequestContext.GetHeader(valueScript)
		} else {
			value = invokeInfo.InvokeNamespace.RequestContext.GetHeader(name)
		}
	case model.DATA_PLACE_PARAM:
		if base.IsNotEmpty(valueScript) {
			value = invokeInfo.InvokeNamespace.RequestContext.Query(valueScript)
		} else {
			value = invokeInfo.InvokeNamespace.RequestContext.Query(name)
		}
	case model.DATA_PLACE_FILE:
		var fileHeader *multipart.FileHeader
		if base.IsNotEmpty(valueScript) {
			fileHeader, err = invokeInfo.InvokeNamespace.RequestContext.FormFile(valueScript)
		} else {
			fileHeader, err = invokeInfo.InvokeNamespace.RequestContext.FormFile(name)
		}
		if err != nil {
			return
		}
		fileInfo := map[string]interface{}{
			"name":       fileHeader.Filename,
			"fileHeader": fileHeader,
		}
		value = fileInfo
	case model.DATA_PLACE_FORM:
		if base.IsNotEmpty(valueScript) {
			value = invokeInfo.InvokeNamespace.RequestContext.PostForm(valueScript)
		} else {
			value = invokeInfo.InvokeNamespace.RequestContext.PostForm(name)
		}
	default:
		if invokeInfo.InvokeNamespace.RequestBody == nil {
			invokeInfo.InvokeNamespace.RequestBody = map[string]interface{}{}
			err = invokeInfo.InvokeNamespace.RequestContext.Bind(&invokeInfo.InvokeNamespace.RequestBody)
			if err != nil {
				return
			}
		}
		value = invokeInfo.InvokeNamespace.RequestBody
		if base.IsNotEmpty(valueScript) {
			switch m := invokeInfo.InvokeNamespace.RequestBody.(type) {
			case map[string]interface{}:
				value = m[valueScript]
			default:
				err = base.NewError("", "request body can not to map[string]interface{}")
				return
			}
		}
	}

	err = invokeInfo.InvokeNamespace.SetData(name, value, nil)
	if err != nil {
		return
	}
	var invokeData *common.InvokeData
	invokeData, err = invokeInfo.InvokeNamespace.GetData(name)
	if err != nil {
		return
	}
	res = invokeData.Value
	return
}

func invokeAddDataInfo(invokeInfo *InvokeInfo, _ string, args []interface{}) (res interface{}, err error) {
	name := args[0].(string)
	var dataType string = args[1].(string)
	var comment string
	var value string
	var isList bool
	var isPage bool
	if len(args) > 2 {
		comment = args[2].(string)
	}
	if len(args) > 3 {
		value = args[3].(string)
	}
	if len(args) > 4 {
		isList = args[4].(bool)
	}
	if len(args) > 5 {
		isPage = args[5].(bool)
	}
	variable := &model.VariableModel{
		Name:     name,
		Comment:  comment,
		Value:    value,
		DataType: dataType,
		IsList:   isList,
		IsPage:   isPage,
	}
	err = invokeInfo.InvokeNamespace.SetDataInfo(variable)
	if err != nil {
		return
	}
	return
}
func invokeNewVariable(invokeInfo *InvokeInfo, _ string, args []interface{}) (res interface{}, err error) {
	name := args[0].(string)
	var dataInfo *common.InvokeDataInfo
	dataInfo, err = invokeInfo.InvokeNamespace.GetDataInfo(name)
	if err != nil {
		return
	}
	if dataInfo.IsList {
		err = invokeInfo.InvokeNamespace.SetData(name, []map[string]interface{}{}, nil)
	} else {
		err = invokeInfo.InvokeNamespace.SetData(name, map[string]interface{}{}, nil)
	}
	if err != nil {
		return
	}
	return
}
func invokeLock(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	// fmt.Println("invoke ", prefixName, " lock")
	return
}
func invokeUnlock(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	// fmt.Println("invoke ", prefixName, " unlock")
	return
}
func invokePush(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var value interface{}
	defer func() {
		if err == nil {
			err = base.PanicError(recover())
		}
		if err != nil {
			if invokeInfo.App.GetLogger() != nil {
				invokeInfo.App.GetLogger().Error("invoke push data [", prefixName, "] value [", value, "] error:", err)
			}
		}
	}()
	var invokeData *common.InvokeData
	invokeData, err = invokeInfo.InvokeNamespace.GetData(prefixName)
	if err != nil {
		return
	}
	value = invokeData.Value
	if value == nil {
		value = []interface{}{}
	}
	switch l := value.(type) {
	case []interface{}:
		value = append(l, args...)
	default:
		err = base.NewError("", "invoke data [", prefixName, "] value can not to []interface{}")
		return
	}
	err = invokeInfo.InvokeNamespace.SetData(prefixName, value, nil)

	// fmt.Println("invoke ", prefixName, " push:", base.ToJSON(args))
	// fmt.Println("SetData name:", prefixName)
	// fmt.Println("SetData value:", base.ToJSON(list))
	return
}
func invokePushs(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var value interface{}
	defer func() {
		if err == nil {
			err = base.PanicError(recover())
		}
		if err != nil {
			if invokeInfo.App.GetLogger() != nil {
				invokeInfo.App.GetLogger().Error("invoke pushs data [", prefixName, "] value [", value, "] error:", err)
			}
		}
	}()
	var invokeData *common.InvokeData
	invokeData, err = invokeInfo.InvokeNamespace.GetData(prefixName)
	if err != nil {
		return
	}
	value = invokeData.Value
	if value == nil {
		value = []interface{}{}
	}
	switch l := value.(type) {
	case []interface{}:
		for _, arg := range args {
			if argList, argListOk := arg.([]interface{}); argListOk {
				value = append(l, argList...)
			} else {
				err = base.NewError("", "invoke data [", prefixName, "] arg can not to []interface{}")
				return
			}
		}
	default:
		err = base.NewError("", "invoke data [", prefixName, "] value can not to []interface{}")
		return
	}
	err = invokeInfo.InvokeNamespace.SetData(prefixName, value, nil)

	return
}
func invokeSqlSelect(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var sqlExecutor common.ISqlExecutor
	sqlExecutor, err = invokeInfo.App.GetSqlExecutor(prefixName)
	if err != nil {
		return
	}
	var sql string
	var params []interface{}
	var columnFieldMap map[string]*model.StructFieldModel

	sql = args[0].(string)
	params = args[1].([]interface{})
	structName := args[2].(string)
	columnFieldMap, err = common.GetColumnFieldMap(invokeInfo.App, structName)
	if err != nil {
		return
	}
	if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
		invokeInfo.App.GetLogger().Debug("execute sql select sql   :", sql)
		invokeInfo.App.GetLogger().Debug("execute sql select params:", base.ToJSON(params))
	}
	res, err = sqlExecutor.Select(sql, params, columnFieldMap)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("execute sql select error sql   :", sql)
			invokeInfo.App.GetLogger().Error("execute sql select error params:", base.ToJSON(params))
			invokeInfo.App.GetLogger().Error("execute sql select error error :", err)
		}
		return
	}
	return
}
func invokeSqlSelectPage(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var sqlExecutor common.ISqlExecutor
	sqlExecutor, err = invokeInfo.App.GetSqlExecutor(prefixName)
	if err != nil {
		return
	}
	var sql string
	var params []interface{}
	var columnFieldMap map[string]*model.StructFieldModel

	sql = args[0].(string)
	params = args[1].([]interface{})
	pageNumber := args[2].(int64)
	pageSize := args[3].(int64)
	structName := args[4].(string)
	columnFieldMap, err = common.GetColumnFieldMap(invokeInfo.App, structName)
	if err != nil {
		return
	}
	if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
		invokeInfo.App.GetLogger().Debug("execute sql select page sql   :", sql)
		invokeInfo.App.GetLogger().Debug("execute sql select page params:", base.ToJSON(params))
	}
	res, err = sqlExecutor.SelectPage(sql, params, pageNumber, pageSize, columnFieldMap)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("execute sql select page error sql   :", sql)
			invokeInfo.App.GetLogger().Error("execute sql select page error params:", base.ToJSON(params))
			invokeInfo.App.GetLogger().Error("execute sql select page error error :", err)
		}
		return
	}
	return
}

func invokeSqlSelectOne(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var sqlExecutor common.ISqlExecutor
	sqlExecutor, err = invokeInfo.App.GetSqlExecutor(prefixName)
	if err != nil {
		return
	}
	var sql string
	var params []interface{}
	var columnFieldMap map[string]*model.StructFieldModel

	sql = args[0].(string)
	params = args[1].([]interface{})
	structName := args[2].(string)
	columnFieldMap, err = common.GetColumnFieldMap(invokeInfo.App, structName)
	if err != nil {
		return
	}
	if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
		invokeInfo.App.GetLogger().Debug("execute sql select one sql   :", sql)
		invokeInfo.App.GetLogger().Debug("execute sql select one params:", base.ToJSON(params))
	}
	res, err = sqlExecutor.SelectOne(sql, params, columnFieldMap)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("execute sql select one error sql   :", sql)
			invokeInfo.App.GetLogger().Error("execute sql select one error params:", base.ToJSON(params))
			invokeInfo.App.GetLogger().Error("execute sql select one error error :", err)
		}
		return
	}
	return
}
func invokeSqlSelectCount(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var sqlExecutor common.ISqlExecutor
	sqlExecutor, err = invokeInfo.App.GetSqlExecutor(prefixName)
	if err != nil {
		return
	}
	var sql string
	var params []interface{}

	sql = args[0].(string)
	params = args[1].([]interface{})
	if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
		invokeInfo.App.GetLogger().Debug("execute sql select count sql   :", sql)
		invokeInfo.App.GetLogger().Debug("execute sql select count params:", base.ToJSON(params))
	}
	res, err = sqlExecutor.SelectCount(sql, params)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("execute sql select count error sql   :", sql)
			invokeInfo.App.GetLogger().Error("execute sql select count error params:", base.ToJSON(params))
			invokeInfo.App.GetLogger().Error("execute sql select count error error :", err)
		}
		return
	}
	return
}
func invokeSqlInsert(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var sqlExecutor common.ISqlExecutor
	sqlExecutor, err = invokeInfo.App.GetSqlExecutor(prefixName)
	if err != nil {
		return
	}
	var sql string
	var params []interface{}

	sql = args[0].(string)
	params = args[1].([]interface{})
	if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
		invokeInfo.App.GetLogger().Debug("execute sql insert sql   :", sql)
		invokeInfo.App.GetLogger().Debug("execute sql insert params:", base.ToJSON(params))
	}
	res, err = sqlExecutor.Insert(sql, params)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("execute sql insert error sql   :", sql)
			invokeInfo.App.GetLogger().Error("execute sql insert error params:", base.ToJSON(params))
			invokeInfo.App.GetLogger().Error("execute sql insert error error :", err)
		}
		return
	}
	return
}
func invokeSqlUpdate(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var sqlExecutor common.ISqlExecutor
	sqlExecutor, err = invokeInfo.App.GetSqlExecutor(prefixName)
	if err != nil {
		return
	}
	var sql string
	var params []interface{}

	sql = args[0].(string)
	params = args[1].([]interface{})
	if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
		invokeInfo.App.GetLogger().Debug("execute sql update sql   :", sql)
		invokeInfo.App.GetLogger().Debug("execute sql update params:", base.ToJSON(params))
	}
	res, err = sqlExecutor.Update(sql, params)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("execute sql update error sql   :", sql)
			invokeInfo.App.GetLogger().Error("execute sql update error params:", base.ToJSON(params))
			invokeInfo.App.GetLogger().Error("execute sql update error error :", err)
		}
		return
	}
	return
}
func invokeSqlDelete(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var sqlExecutor common.ISqlExecutor
	sqlExecutor, err = invokeInfo.App.GetSqlExecutor(prefixName)
	if err != nil {
		return
	}
	var sql string
	var params []interface{}

	sql = args[0].(string)
	params = args[1].([]interface{})
	if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
		invokeInfo.App.GetLogger().Debug("execute sql delete sql   :", sql)
		invokeInfo.App.GetLogger().Debug("execute sql delete params:", base.ToJSON(params))
	}
	res, err = sqlExecutor.Delete(sql, params)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("execute sql delete error sql   :", sql)
			invokeInfo.App.GetLogger().Error("execute sql delete error params:", base.ToJSON(params))
			invokeInfo.App.GetLogger().Error("execute sql delete error error :", err)
		}
		return
	}
	return
}
func invokeRedisSet(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var key string
	var value string
	if base.IsEmptyObj(args[0]) {
		err = base.NewError("", "redis set key not defind")
	}
	if base.IsEmptyObj(args[1]) {
		err = base.NewError("", "redis set value not defind")
	}
	switch v := args[0].(type) {
	case string:
		key = v
	default:
		key = fmt.Sprint(v)
	}
	switch v := args[1].(type) {
	case string:
		value = v
	case map[string]interface{}, []interface{}, []map[string]interface{}:
		value = base.ToJSON(v)
	default:
		value = fmt.Sprint(v)
	}
	var redisExecutor common.IRedisExecutor
	redisExecutor, err = invokeInfo.App.GetRedisExecutor(prefixName)
	if err != nil {
		return
	}
	key = redisExecutor.FormatKey(key)
	if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
		invokeInfo.App.GetLogger().Debug("redis set key  :", key)
		invokeInfo.App.GetLogger().Debug("redis set value:", value)
	}
	err = redisExecutor.Set(key, value)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("redis set error key  :", key)
			invokeInfo.App.GetLogger().Error("redis set error value:", value)
			invokeInfo.App.GetLogger().Error("redis set error error:", err)
		}
		return
	}
	return
}
func invokeRedisGet(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var key string
	if base.IsEmptyObj(args[0]) {
		err = base.NewError("", "redis get key not defind")
	}
	switch v := args[0].(type) {
	case string:
		key = v
	default:
		key = fmt.Sprint(v)
	}
	var redisExecutor common.IRedisExecutor
	redisExecutor, err = invokeInfo.App.GetRedisExecutor(prefixName)
	if err != nil {
		return
	}
	key = redisExecutor.FormatKey(key)
	if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
		invokeInfo.App.GetLogger().Debug("redis get key  :", key)
	}
	res, err = redisExecutor.Get(key)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("redis get error key  :", key)
			invokeInfo.App.GetLogger().Error("redis get error error:", err)
		}
		return
	}
	return
}
func invokeRedisDel(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	var key string
	if base.IsEmptyObj(args[0]) {
		err = base.NewError("", "redis del key not defind")
	}
	switch v := args[0].(type) {
	case string:
		key = v
	default:
		key = fmt.Sprint(v)
	}
	var redisExecutor common.IRedisExecutor
	redisExecutor, err = invokeInfo.App.GetRedisExecutor(prefixName)
	if err != nil {
		return
	}
	key = redisExecutor.FormatKey(key)
	if invokeInfo.App.GetLogger() != nil && invokeInfo.App.GetLogger().OutDebug() {
		invokeInfo.App.GetLogger().Debug("redis del key  :", key)
	}
	res, err = redisExecutor.Del(key)
	if err != nil {
		if invokeInfo.App.GetLogger() != nil {
			invokeInfo.App.GetLogger().Error("redis del error key  :", key)
			invokeInfo.App.GetLogger().Error("redis del error error:", err)
		}
		return
	}
	return
}

func invokeService(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	callServiceName := args[0].(string)
	callService := invokeInfo.App.GetContext().GetService(callServiceName)
	if callService == nil {
		err = base.NewErrorServiceIsNull("service [", callServiceName, "] not defind")
		return
	}

	var callInvokeNamespace *common.InvokeNamespace
	callInvokeNamespace, err = common.NewInvokeNamespace(invokeInfo.App)
	if err != nil {
		return
	}

	for index, callVariable := range callService.InVariables {
		value := args[index+1]
		err = callInvokeNamespace.SetDataInfo(callVariable)
		if err != nil {
			return
		}
		err = callInvokeNamespace.SetData(callVariable.Name, value, nil)
		if err != nil {
			return
		}
	}
	// fmt.Println("callInvokeNamespace:", base.ToJSON(callInvokeNamespace))
	res, err = invokeInfo.App.InvokeService(callService, callInvokeNamespace)
	if err != nil {
		return
	}
	return
}

func invokeFileSave(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	name := args[0].(string)
	dir := args[1].(string)

	var reader io.Reader
	var bytes []byte

	if args[2] != nil {
		ioReader, ioReaderOk := args[2].(io.Reader)
		if !ioReaderOk {
			err = base.NewError("", "file reader must io.Reader")
			return
		}
		reader = ioReader
	} else if args[3] != nil {
		ioBytes, ioBytesOk := args[3].([]byte)
		if !ioBytesOk {
			err = base.NewError("", "file bytes must []byte")
			return
		}
		bytes = ioBytes
	}
	if reader != nil {
		if Closer, CloserOk := reader.(io.Closer); CloserOk {
			defer Closer.Close()
		}
	}
	if base.IsEmpty(name) {
		err = base.NewError("", "file save name not defind")
		return
	}
	if base.IsEmpty(dir) {
		err = base.NewError("", "file save dir not defind")
		return
	}
	if reader == nil && bytes == nil {
		err = base.NewError("", "file save reader or bytes not defind")
		return
	}
	// 创建目录
	var exists bool
	exists, err = base.PathExists(dir)
	if err != nil {
		return
	}
	if !exists {
		os.MkdirAll(dir, 0777)
	}
	savePath := dir + "/" + name
	exists, err = base.PathExists(savePath)
	if err != nil {
		return
	}
	if exists {
		err = base.NewError("", "file [", savePath, "] is exists")
		return
	}
	file, err := os.Create(savePath)
	if err != nil {
		return
	}
	defer file.Close()

	if reader != nil {
		_, err = io.Copy(file, reader)
	} else {
		_, err = file.Write(bytes)
	}
	if err != nil {
		return
	}
	var abs string
	abs, err = filepath.Abs(savePath)
	if err != nil {
		return
	}
	fileInfo := map[string]interface{}{
		"name":         name,
		"dir":          dir,
		"path":         savePath,
		"absolutePath": filepath.ToSlash(abs),
		"file":         file,
	}
	res = fileInfo
	return
}
func invokeFileGet(invokeInfo *InvokeInfo, prefixName string, args []interface{}) (res interface{}, err error) {
	path := args[0].(string)

	if base.IsEmpty(path) {
		err = base.NewError("", "file get path not defind")
		return
	}
	var exists bool
	exists, err = base.PathExists(path)
	if err != nil {
		return
	}
	var fileInfo map[string]interface{} = nil
	if exists {
		var file *os.File
		file, err = os.Open(path)
		if err != nil {
			return
		}
		path = strings.ReplaceAll(path, "\\", "/")
		name := path
		dir := ""
		if strings.Contains(path, "/") {
			name = path[strings.LastIndex(path, "/")+1:]
			dir = path[:strings.LastIndex(path, "/")]
		}
		var abs string
		abs, err = filepath.Abs(path)
		if err != nil {
			return
		}
		fileInfo = map[string]interface{}{
			"name":         name,
			"dir":          dir,
			"path":         path,
			"absolutePath": filepath.ToSlash(abs),
			"file":         file,
		}
	}

	res = fileInfo
	return
}
