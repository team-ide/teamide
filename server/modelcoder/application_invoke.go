package modelcoder

func invokeModel(application *Application, model interface{}, variable *invokeVariable) (res interface{}, err error) {
	service, isService := interface{}(model).(ServiceModel)
	var serviceType *ServiceModelType
	if isService {
		if application.OutDebug() {
			application.Debug("invoke model is service type [", ToJSON(service.GetType()), "]")
		}
		serviceType = service.GetType()
		if serviceType == nil {
			err = newErrorServiceTypeIsWrong("service [", service.GetName(), "] model type [", service.GetType(), "] parsing failed")
			return
		}
	}
	dao, isDao := interface{}(model).(DaoModel)
	var daoType *DaoModelType
	if isDao {
		if application.OutDebug() {
			application.Debug("invoke model is dao type [", ToJSON(dao.GetType()), "]")
		}
		daoType = dao.GetType()
		if daoType == nil {
			err = newErrorDaoTypeIsWrong("dao [", dao.GetName(), "] model type [", dao.GetType(), "] parsing failed")
			return
		}
	}
	if isService {
		res, err = serviceType.Execute(application, service, variable)
	}
	if isDao {
		res, err = daoType.Execute(application, dao, variable)
	}
	if err != nil {
		return
	}
	return
}

func processParams(application *Application, params []*ParamModel, variable *invokeVariable) (err error) {
	if len(params) == 0 {
		return
	}
	for _, one := range params {
		err = processParam(application, one, variable)
		if err != nil {
			return
		}
	}
	return
}

func newCallVnvokeVariable(application *Application, variable *invokeVariable, callParams []string, targetParams []*ParamModel) (callVariable *invokeVariable, err error) {
	// 调用外部Model 需要重置invokeVariable
	callVariable = variable.Clone()

	// 根据Call传参 解析当前变量载体应该传输哪些

	//
	if len(callParams) == 0 {
		if len(targetParams) > 0 {
			for _, targetParam := range targetParams {
				callParamData := variable.GetParamData(targetParam.Name)
				if callParamData != nil {
					callParams = append(callParams, targetParam.Name)
				}
			}
		}
	}
	if len(callParams) > 0 {
		for _, callParamName := range callParams {
			callParamData := variable.GetParamData(callParamName)
			if callParamData == nil {
				err = newErrorParamIsNull("call param [" + callParamName + "] not defind")
				return
			}
			callVariable.AddParamData(callParamData)
		}
	}
	return
}

func processParam(application *Application, param *ParamModel, variable *invokeVariable) (err error) {
	if param == nil {
		return
	}
	if application.OutDebug() {
		application.Debug("process param:", ToJSON(param))
	}
	paramData := variable.GetParamData(param.Name)
	if paramData == nil {
		err = newErrorParamIsNull("param [" + param.Name + "] not defind")
		return
	}

	data := paramData.Data

	dataType := GetParamModelDataType(param.DataType)

	// 数据类型为空 表示结构体
	if dataType == nil {
		structModel := application.context.GetStruct(param.DataType)
		if structModel == nil {
			err = newErrorStructIsNull("")
			return
		}
	}

	// 根据类型转换位新数据
	paramData = &ParamData{
		Name: param.Name,
		Data: data,
	}

	if application.OutDebug() {
		application.Debug("process param [", paramData.Name, "] data:", ToJSON(paramData.Data))
	}
	variable.AddParamData(paramData)

	return
}
