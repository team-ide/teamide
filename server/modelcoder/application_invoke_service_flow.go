package modelcoder

func invokeServiceFlow(application *Application, service ServiceModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke service flow start")
	}
	serviceFlow := interface{}(service).(*ServiceFlowModel)
	start := serviceFlow.GetStartStep()

	if start == nil {
		application.Warn("invoke service flow start step is null")
		return
	}

	if application.OutDebug() {
		application.Debug("invoke service flow start step:", ToJSON(start))
	}
	res, err = start.GetType().Execute(application, serviceFlow, start, variable)
	if application.OutDebug() {
		application.Debug("invoke service flow end")
	}
	return
}

func invokeServiceFlowStepStart(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke service flow step [start] start")
	}
	stepStart := interface{}(step).(*ServiceFlowStepStartModel)

	next := flow.GetStep(stepStart.Next)
	res, err = invokeServiceFlowStep(application, flow, next, variable)
	if application.OutDebug() {
		application.Debug("invoke service flow step [start] end")
	}
	return
}

func invokeServiceFlowStepData(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke service flow step [data] start")
	}
	stepData := interface{}(step).(*ServiceFlowStepDataModel)

	next := flow.GetStep(stepData.Next)
	res, err = invokeServiceFlowStep(application, flow, next, variable)
	if application.OutDebug() {
		application.Debug("invoke service flow step [data] end")
	}
	return
}

func invokeServiceFlowStepDecision(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke service flow step [decision] start")
	}
	stepDecision := interface{}(step).(*ServiceFlowStepDecisionModel)

	next := flow.GetStep(stepDecision.IfNext)
	res, err = invokeServiceFlowStep(application, flow, next, variable)
	if application.OutDebug() {
		application.Debug("invoke service flow step [decision] end")
	}
	return
}

func invokeServiceFlowStepDao(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke service flow step [dao] start")
	}
	stepDao := interface{}(step).(*ServiceFlowStepDaoModel)

	dao := application.context.GetDao(stepDao.DaoName)

	if dao == nil {
		err = newErrorDaoIsNull("invoke dao model [", stepDao.DaoName, "] is null")
		return
	}
	res, err = invokeModel(application, dao, variable)
	if err != nil {
		return
	}

	next := flow.GetStep(stepDao.Next)
	res, err = invokeServiceFlowStep(application, flow, next, variable)
	if application.OutDebug() {
		application.Debug("invoke service flow step [dao] end")
	}
	return
}

func invokeServiceFlowStepService(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke service flow step [service] start")
	}
	stepService := interface{}(step).(*ServiceFlowStepServiceModel)

	service := application.context.GetDao(stepService.ServiceName)

	if service == nil {
		err = newErrorServiceIsNull("invoke service model [", stepService.ServiceName, "] is null")
		return
	}
	res, err = invokeModel(application, service, variable)
	if err != nil {
		return
	}

	next := flow.GetStep(stepService.Next)
	res, err = invokeServiceFlowStep(application, flow, next, variable)
	if application.OutDebug() {
		application.Debug("invoke service flow step [service] end")
	}
	return
}

func invokeServiceFlowStepEnd(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke service flow step [end] start")
	}
	if application.OutDebug() {
		application.Debug("invoke service flow step [end] end")
	}
	return
}

func invokeServiceFlowStepError(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke service flow step [error] start")
	}
	if application.OutDebug() {
		application.Debug("invoke service flow step [error] end")
	}
	return
}

func invokeServiceFlowStep(application *Application, flow *ServiceFlowModel, next ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error) {
	if next == nil {
		return
	}
	if application.OutDebug() {
		application.Debug("invoke service flow step [next] start")
	}

	if application.OutDebug() {
		application.Debug("invoke service flow next step:", ToJSON(next))
	}
	res, err = next.GetType().Execute(application, flow, next, variable)
	if application.OutDebug() {
		application.Debug("invoke service flow step [next] end")
	}
	return
}
