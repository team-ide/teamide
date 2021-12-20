package modelcoder

func invokeDaoSqlSelectOne(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke dao sql select one start")
	}
	sqlSelectOne := interface{}(dao).(*DaoSqlSelectOneModel)
	println(dao)
	println(sqlSelectOne)
	if application.OutDebug() {
		application.Debug("invoke dao sql select one end")
	}
	return
}

func invokeDaoSqlSelectList(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	return
}

func invokeDaoSqlSelectPage(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	return
}

func invokeDaoSqlSelectCount(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	return
}

func invokeDaoSqlInsert(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	return
}

func invokeDaoSqlUpdate(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	return
}

func invokeDaoSqlDelete(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	return
}
