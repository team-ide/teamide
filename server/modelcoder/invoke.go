package modelcoder

func (this_ *Application) InvokeServiceByName(name string, variable *invokeVariable) (res interface{}, err error) {
	if this_.OutDebug() {
		this_.Debug("invoke service start , name:", name, ", variable:", ToJSON(variable))
	}
	if variable == nil {
		err = newErrorVariableIsNull("invoke service variable is null")
		return
	}
	service := this_.context.GetService(name)
	if service == nil {
		err = newErrorServiceIsNull("invoke service model is null")
		return
	}
	res, err = invokeModel(this_, service, variable)
	if err != nil {
		this_.Error("invoke service error , name:", name, ", error:", err)
		return
	}
	if this_.OutDebug() {
		this_.Debug("invoke service end , name:", name, ", result:", ToJSON(res))
	}
	return
}

func (this_ *Application) InvokeService(service ServiceModel, variable *invokeVariable) (res interface{}, err error) {
	if this_.OutDebug() {
		this_.Debug("invoke service start , service:", ToJSON(service), ", variable:", ToJSON(variable))
	}
	if service == nil {
		err = newErrorServiceIsNull("invoke service model is null")
		return
	}
	if variable == nil {
		err = newErrorVariableIsNull("invoke service variable is null")
		return
	}
	res, err = invokeModel(this_, service, variable)
	if err != nil {
		this_.Error("invoke service error , service:", ToJSON(service), ", error:", err)
		return
	}
	if this_.OutDebug() {
		this_.Debug("invoke service end , service:", ToJSON(service), ", result:", ToJSON(res))
	}
	return
}

func (this_ *Application) InvokeDaoByName(name string, variable *invokeVariable) (res interface{}, err error) {
	if this_.OutDebug() {
		this_.Debug("invoke dao start , name:", name, ", variable:", ToJSON(variable))
	}
	if variable == nil {
		err = newErrorVariableIsNull("invoke dao variable is null")
		return
	}
	dao := this_.context.GetDao(name)
	if dao == nil {
		err = newErrorDaoIsNull("invoke dao model is null")
		return
	}
	res, err = invokeModel(this_, dao, variable)
	if err != nil {
		this_.Error("invoke dao error , name:", name, ", error:", err)
		return
	}
	if this_.OutDebug() {
		this_.Debug("invoke dao end , name:", name, ", result:", ToJSON(res))
	}
	return
}

func (this_ *Application) InvokeDao(dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	if this_.OutDebug() {
		this_.Debug("invoke dao start , dao:", ToJSON(dao), ", variable:", ToJSON(variable))
	}
	if dao == nil {
		err = newErrorDaoIsNull("invoke dao model is null")
		return
	}
	if variable == nil {
		err = newErrorVariableIsNull("invoke dao variable is null")
		return
	}
	res, err = invokeModel(this_, dao, variable)
	if err != nil {
		this_.Error("invoke dao error , dao:", ToJSON(dao), ", error:", err)
		return
	}
	if this_.OutDebug() {
		this_.Debug("invoke dao end , dao:", ToJSON(dao), ", result:", ToJSON(res))
	}
	return
}

func invokeModel(application *Application, model interface{}, variable *invokeVariable) (res interface{}, err error) {
	service, isService := interface{}(model).(ServiceModel)
	var parameters []*VariableModel
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
		parameters = service.GetParameters()
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
		parameters = dao.GetParameters()
	}

	err = processVariableDatas(application, parameters, variable)

	if err != nil {
		return
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
