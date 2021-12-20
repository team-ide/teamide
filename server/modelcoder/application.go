package modelcoder

type Application struct {
	context *ApplicationContext
}

type invokeVariable struct {
	data        interface{}
	parent      *invokeVariable
	invokeCache map[string]interface{}
}

func NewApplication(applicationModel *ApplicationModel) *Application {
	if applicationModel == nil {
		return nil
	}
	res := &Application{}
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
	service := this_.context.GetService(name)
	res, err = this_.InvokeService(service, variable)
	return
}

func (this_ *Application) InvokeService(service *ServiceModel, variable *invokeVariable) (res interface{}, err error) {
	if service == nil {
		err = newErrorServiceIsNull("服务模型为空")
		return
	}
	if variable == nil {
		err = newErrorVariableIsNull("调用变量为空")
		return
	}
	return
}

func (this_ *Application) InvokeDaoByName(name string, variable *invokeVariable) (res interface{}, err error) {
	dao := this_.context.GetDao(name)
	res, err = this_.InvokeDao(dao, variable)
	return
}

func (this_ *Application) InvokeDao(dao *DaoModel, variable *invokeVariable) (res interface{}, err error) {
	if dao == nil {
		err = newErrorDaoIsNull("数据模型为空")
		return
	}
	if variable == nil {
		err = newErrorVariableIsNull("调用变量为空")
		return
	}
	return
}
