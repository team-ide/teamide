package modelcoder

type Application struct {
	context *applicationContext
	logger  logger
}
type logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type invokeVariable struct {
	data        interface{}
	parent      *invokeVariable
	invokeCache map[string]interface{}
}

func NewApplication(applicationModel *ApplicationModel, logger logger) *Application {
	if applicationModel == nil {
		return nil
	}
	res := &Application{
		logger: logger,
	}
	res.context = newApplicationContext(applicationModel)
	return res
}

func (this_ *Application) NewInvokeVariable(data interface{}) *invokeVariable {
	res := &invokeVariable{
		data:        data,
		parent:      nil,
		invokeCache: map[string]interface{}{},
	}
	return res
}

func (this_ *Application) InvokeServiceByName(name string, variable *invokeVariable) (res interface{}, err error) {
	this_.logger.Debug("invoke service start , name:", name, ", variable:", variable)
	service := this_.context.GetService(name)
	res, err = this_.InvokeService(service, variable)
	if err != nil {
		this_.logger.Error("invoke service error , name:", name, ", variable:", variable, ", error:", err)
		return
	}
	this_.logger.Debug("invoke service end , name:", name, ", variable:", variable, ", result:", res)
	return
}

func (this_ *Application) InvokeService(service ServiceModel, variable *invokeVariable) (res interface{}, err error) {
	this_.logger.Debug("invoke service start , service:", service, ", variable:", variable)
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
		this_.logger.Error("invoke service error , service:", service, ", variable:", variable, ", error:", err)
		return
	}
	this_.logger.Debug("invoke service end , service:", service, ", variable:", variable, ", result:", res)
	return
}

func (this_ *Application) InvokeDaoByName(name string, variable *invokeVariable) (res interface{}, err error) {
	this_.logger.Debug("invoke dao start , name:", name, ", variable:", variable)
	dao := this_.context.GetDao(name)
	res, err = this_.InvokeDao(dao, variable)
	if err != nil {
		this_.logger.Error("invoke dao error , name:", name, ", variable:", variable, ", error:", err)
		return
	}
	this_.logger.Debug("invoke dao end , name:", name, ", variable:", variable, ", result:", res)
	return
}

func (this_ *Application) InvokeDao(dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	this_.logger.Debug("invoke dao start , dao:", dao, ", variable:", variable)
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
		this_.logger.Error("invoke dao error , dao:", dao, ", variable:", variable, ", error:", err)
		return
	}
	this_.logger.Debug("invoke dao end , dao:", dao, ", variable:", variable, ", result:", res)
	return
}
