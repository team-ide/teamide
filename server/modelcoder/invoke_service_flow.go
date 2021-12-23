package modelcoder

func invokeServiceFlow(application *Application, service ServiceModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke service [", service.GetName(), "] start, variable:", ToJSON(variable))
	}
	serviceFlow := interface{}(service).(*ServiceFlowModel)

	start := serviceFlow.GetStartStep()

	if start == nil {
		application.Warn("invoke service [", service.GetName(), "] start step is null")
		return
	}
	res, err = invokeServiceFlowStep(application, serviceFlow, start, variable)
	if application.OutDebug() {
		application.Debug("invoke service [", service.GetName(), "] end")
	}
	return
}

func invokeServiceFlowStep(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke service flow step [", step.GetName(), "] start, variable:", ToJSON(variable))
	}
	var parameters []*VariableModel

	stepStart, startOk := interface{}(step).(*ServiceFlowStepStartModel)

	stepData, dataOk := interface{}(step).(*ServiceFlowStepDataModel)

	stepDecision, decisionOk := interface{}(step).(*ServiceFlowStepDecisionModel)

	stepDao, daoOk := interface{}(step).(*ServiceFlowStepDaoModel)

	stepService, serviceOk := interface{}(step).(*ServiceFlowStepServiceModel)

	stepEnd, endOk := interface{}(step).(*ServiceFlowStepEndModel)

	stepError, errorOk := interface{}(step).(*ServiceFlowStepErrorModel)

	err = processVariableDatas(application, parameters, variable)

	if err != nil {
		return
	}
	var next string
	if startOk {
		next, res, err = invokeServiceFlowStepStart(application, flow, stepStart, variable)
	}
	if dataOk {
		next, res, err = invokeServiceFlowStepData(application, flow, stepData, variable)
	}
	if decisionOk {
		next, res, err = invokeServiceFlowStepDecision(application, flow, stepDecision, variable)
	}
	if daoOk {
		next, res, err = invokeServiceFlowStepDao(application, flow, stepDao, variable)
	}
	if serviceOk {
		next, res, err = invokeServiceFlowStepService(application, flow, stepService, variable)
	}
	if endOk {
		next, res, err = invokeServiceFlowStepEnd(application, flow, stepEnd, variable)
	}
	if errorOk {
		next, res, err = invokeServiceFlowStepError(application, flow, stepError, variable)
	}

	if next != "" {
		nextStep := flow.GetStep(next)
		if nextStep == nil {
			err = newErrorServiceStepIsWrong("invoke service flow next step [", next, "] not defind")
			return
		}
		res, err = invokeServiceFlowStep(application, flow, nextStep, variable)
	}
	if application.OutDebug() {
		application.Debug("invoke service flow step [", step.GetName(), "] end")
	}
	return
}

func invokeServiceFlowStepStart(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (next string, res interface{}, err error) {
	stepStart := interface{}(step).(*ServiceFlowStepStartModel)

	next = stepStart.Next
	return
}

func invokeServiceFlowStepData(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (next string, res interface{}, err error) {
	stepData := interface{}(step).(*ServiceFlowStepDataModel)

	next = stepData.Next
	return
}

func invokeServiceFlowStepDecision(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (next string, res interface{}, err error) {
	stepDecision := interface{}(step).(*ServiceFlowStepDecisionModel)

	next = stepDecision.IfNext
	return
}

func invokeServiceFlowStepDao(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (next string, res interface{}, err error) {
	stepDao := interface{}(step).(*ServiceFlowStepDaoModel)

	dao := application.context.GetDao(stepDao.DaoName)

	if dao == nil {
		err = newErrorDaoIsNull("invoke dao model [", stepDao.DaoName, "] is null")
		return
	}

	var callVariable *invokeVariable
	callParams := []string{}
	// 调用外部Model 需要重置invokeVariable
	callVariable, err = newCallInvokeVariable(application, variable, callParams, dao.GetParameters())
	if err != nil {
		return
	}
	res, err = invokeModel(application, dao, callVariable)
	if err != nil {
		return
	}

	next = stepDao.Next
	return
}

func invokeServiceFlowStepService(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (next string, res interface{}, err error) {

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

	next = stepService.Next
	return
}

func invokeServiceFlowStepEnd(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (next string, res interface{}, err error) {

	return
}

func invokeServiceFlowStepError(application *Application, flow *ServiceFlowModel, step ServiceFlowStepModel, variable *invokeVariable) (next string, res interface{}, err error) {

	return
}
