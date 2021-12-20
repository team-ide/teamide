package modelcoder

type Application struct {
	context *applicationContext
	logger  logger
}
type logger interface {
	OutDebug() bool
	OutInfo() bool
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

func (this_ *Application) OutDebug() bool {
	return this_.logger.OutDebug()
}
func (this_ *Application) Debug(args ...interface{}) {
	if this_.OutDebug() {
		this_.logger.Debug(args...)
	}
}
func (this_ *Application) OutInfo() bool {
	return this_.logger.OutInfo()
}
func (this_ *Application) Info(args ...interface{}) {
	if this_.OutInfo() {
		this_.logger.Info(args...)
	}
}
func (this_ *Application) Warn(args ...interface{}) {
	this_.logger.Warn(args...)
}
func (this_ *Application) Error(args ...interface{}) {
	this_.logger.Error(args...)
}

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
		this_.Error("invoke service error , name:", name, ", variable:", ToJSON(variable), ", error:", err)
		return
	}
	if this_.OutDebug() {
		this_.Debug("invoke service end , name:", name, ", variable:", ToJSON(variable), ", result:", ToJSON(res))
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
		this_.Error("invoke service error , service:", ToJSON(service), ", variable:", ToJSON(variable), ", error:", err)
		return
	}
	if this_.OutDebug() {
		this_.Debug("invoke service end , service:", ToJSON(service), ", variable:", ToJSON(variable), ", result:", ToJSON(res))
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
		this_.Error("invoke dao error , name:", name, ", variable:", ToJSON(variable), ", error:", err)
		return
	}
	if this_.OutDebug() {
		this_.Debug("invoke dao end , name:", name, ", variable:", ToJSON(variable), ", result:", ToJSON(res))
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
		this_.Error("invoke dao error , dao:", ToJSON(dao), ", variable:", ToJSON(variable), ", error:", err)
		return
	}
	if this_.OutDebug() {
		this_.Debug("invoke dao end , dao:", ToJSON(dao), ", variable:", ToJSON(variable), ", result:", ToJSON(res))
	}
	return
}
