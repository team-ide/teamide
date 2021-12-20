package modelcoder

type DaoModel interface {
	GetName() string        // 名称，同一个应用中唯一
	GetType() *DaoModelType // 类型
}

type DaoModelType struct {
	Value   string
	Text    string
	Execute func(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error)
}

var (
	daoModelTypes []*DaoModelType

	DAO_SQL_SELECT_ONE   = newDaoModelType("SQL_SELECT_ONE", "Sql-查询单个", invokeDaoSqlSelectOne)
	DAO_SQL_SELECT_LIST  = newDaoModelType("SQL_SELECT_LIST", "Sql-查询列表", invokeDaoSqlSelectList)
	DAO_SQL_SELECT_PAGE  = newDaoModelType("SQL_SELECT_PAGE", "Sql-分页查询", invokeDaoSqlSelectPage)
	DAO_SQL_SELECT_COUNT = newDaoModelType("SQL_SELECT_COUNT", "Sql-统计查询", invokeDaoSqlSelectCount)
	DAO_SQL_INSERT       = newDaoModelType("SQL_INSERT", "Sql-新增", invokeDaoSqlInsert)
	DAO_SQL_UPDATE       = newDaoModelType("SQL_UPDATE", "Sql-更新", invokeDaoSqlUpdate)
	DAO_SQL_DELETE       = newDaoModelType("SQL_DELETE", "Sql-删除", invokeDaoSqlDelete)

	DAO_HTTP_GET    = newDaoModelType("HTTP_GET", "Http-Get", invokeDaoHttpGet)
	DAO_HTTP_POST   = newDaoModelType("HTTP_POST", "Http-Post", invokeDaoHttpPost)
	DAO_HTTP_HEAD   = newDaoModelType("HTTP_HEAD", "Http-Head", invokeDaoHttpHead)
	DAO_HTTP_PUT    = newDaoModelType("HTTP_PUT", "Http-Put", invokeDaoHttpPut)
	DAO_HTTP_PATCH  = newDaoModelType("HTTP_PATCH", "Http-Patch", invokeDaoHttpPatch)
	DAO_HTTP_DELETE = newDaoModelType("HTTP_DELETE", "Http-Delete", invokeDaoHttpDelete)
)

func newDaoModelType(value, text string, execute func(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error)) *DaoModelType {
	res := &DaoModelType{
		Value:   value,
		Text:    text,
		Execute: execute,
	}
	daoModelTypes = append(daoModelTypes, res)
	return res
}

func GetDaoModelType(value string) *DaoModelType {
	for _, one := range daoModelTypes {
		if one.Value == value {
			return one
		}
	}
	return nil
}
