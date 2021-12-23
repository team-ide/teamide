package modelcoder

func invokeDaoSqlSelectOne(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] start,variable:", ToJSON(variable))
	}
	sqlSelect := interface{}(dao).(*DaoSqlSelectOneModel)

	println(sqlSelect)

	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] end")
	}
	return
}

func invokeDaoSqlSelectList(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] start,variable:", ToJSON(variable))
	}
	sqlSelect := interface{}(dao).(*DaoSqlSelectListModel)
	println(sqlSelect)
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] end")
	}
	return
}

func invokeDaoSqlSelectPage(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] start,variable:", ToJSON(variable))
	}
	sqlSelect := interface{}(dao).(*DaoSqlSelectPageModel)

	println(sqlSelect)
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] end")
	}
	return
}

func invokeDaoSqlSelectCount(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] start,variable:", ToJSON(variable))
	}
	sqlSelect := interface{}(dao).(*DaoSqlSelectCountModel)

	println(sqlSelect)
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] end")
	}
	return
}

func invokeDaoSqlInsert(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] start,variable:", ToJSON(variable))
	}
	sqlInsert := interface{}(dao).(*DaoSqlInsertModel)

	var sql string
	var sqlParams []interface{}
	sql, sqlParams, err = getSqlInsertSqlParams(sqlInsert, application, variable)
	if err != nil {
		return
	}
	err = application.executeSqlInsert(sqlInsert.Database, sql, sqlParams)
	if err != nil {
		return
	}
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] end")
	}
	return
}

func invokeDaoSqlUpdate(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] start,variable:", ToJSON(variable))
	}
	sqlUpdate := interface{}(dao).(*DaoSqlUpdateModel)

	println(sqlUpdate)
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] end")
	}
	return
}

func invokeDaoSqlDelete(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) {
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] start,variable:", ToJSON(variable))
	}
	sqlDelete := interface{}(dao).(*DaoSqlDeleteModel)

	println(sqlDelete)
	if application.OutDebug() {
		application.Debug("invoke dao sql [", dao.GetName(), "] end")
	}
	return
}
