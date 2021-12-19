package organizationService

import sqlModel "server/model/sql"

var (
	selectSqlModels = []*sqlModel.Select{}
	insertSqlModels = []*sqlModel.Insert{}
	updateSqlModels = []*sqlModel.Update{}
	deleteSqlModels = []*sqlModel.Delete{}
)

func newSelectSql(model *sqlModel.Select) *sqlModel.Select {
	selectSqlModels = append(selectSqlModels, model)
	return model
}

func newInsertSql(model *sqlModel.Insert) *sqlModel.Insert {
	insertSqlModels = append(insertSqlModels, model)
	return model
}

func newUpdateSql(model *sqlModel.Update) *sqlModel.Update {
	updateSqlModels = append(updateSqlModels, model)
	return model
}

func newDeleteSql(model *sqlModel.Delete) *sqlModel.Delete {
	deleteSqlModels = append(deleteSqlModels, model)
	return model
}

func (this_ *OrganizationService) GetSelectSqls() []*sqlModel.Select {

	return selectSqlModels
}

func (this_ *OrganizationService) GetInsertSqls() []*sqlModel.Insert {

	return insertSqlModels
}

func (this_ *OrganizationService) GetUpdateSqls() []*sqlModel.Update {

	return updateSqlModels
}

func (this_ *OrganizationService) GetDeleteSqls() []*sqlModel.Delete {

	return deleteSqlModels
}
