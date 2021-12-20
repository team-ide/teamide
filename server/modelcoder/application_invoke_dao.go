package modelcoder

func invokeDaoSqlSelectOne(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	application.logger.Debug("invoke dao sql select one start")
	sqlSelectOne := interface{}(dao).(*DaoSqlSelectOne)
	println(dao)
	println(sqlSelectOne)
	application.logger.Debug("invoke dao sql select one end")
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
