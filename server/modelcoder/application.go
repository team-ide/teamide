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
	Parent      *invokeVariable        `json:"-"`
	ParamDatas  []*ParamData           `json:"paramDatas,omitempty"`
	InvokeCache map[string]interface{} `json:"invokeCache,omitempty"`
}

func (this_ *invokeVariable) GetParamData(name string) *ParamData {
	if len(this_.ParamDatas) == 0 {
		return nil
	}
	for _, one := range this_.ParamDatas {
		if one.Name == name {
			return one
		}
	}
	return nil
}

func (this_ *invokeVariable) AddParamData(paramDatas ...*ParamData) *invokeVariable {
	this_.ParamDatas = append(this_.ParamDatas, paramDatas...)
	return this_
}

func (this_ *invokeVariable) Clone() *invokeVariable {
	res := &invokeVariable{
		ParamDatas:  []*ParamData{},
		Parent:      this_,
		InvokeCache: map[string]interface{}{},
	}
	return res
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

func (this_ *Application) NewInvokeVariable(paramDatas ...*ParamData) *invokeVariable {
	res := &invokeVariable{
		ParamDatas:  []*ParamData{},
		Parent:      nil,
		InvokeCache: map[string]interface{}{},
	}
	if len(paramDatas) > 0 {
		res.AddParamData(paramDatas...)
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
func (this_ *Application) executeSqlInsert(database string, sql string, sqlParams []interface{}) (err error) {
	if this_.OutDebug() {
		this_.Debug("execute sql insert sql   :", sql)
		this_.Debug("execute sql insert params:", ToJSON(sqlParams))
	}

	return
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
