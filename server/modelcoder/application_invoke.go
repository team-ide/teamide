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
	err = invokeModelBefore(application, variable)
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
	err = invokeModelAfter(application, variable)
	if err != nil {
		return
	}
	return
}

func invokeModelBefore(application *Application, variable *invokeVariable) (err error) {
	if application.OutDebug() {
		application.Debug("invoke model before")
	}

	return
}

func invokeModelAfter(application *Application, variable *invokeVariable) (err error) {
	if application.OutDebug() {
		application.Debug("invoke model after")
	}

	return
}
