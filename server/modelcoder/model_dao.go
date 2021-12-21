package modelcoder

type DaoModel interface {
	GetName() string        // 名称，同一个应用中唯一
	GetType() *DaoModelType // 类型
	GetParams() []*ParamModel
}

type DaoModelType struct {
	Value   string                                                                                              `json:"value,omitempty"`
	Text    string                                                                                              `json:"text,omitempty"`
	Execute func(application *Application, dao DaoModel, variable *invokeVariable) (res interface{}, err error) `json:"-"`
}

var (
	daoModelTypes []*DaoModelType

	DAO_SQL_SELECT_ONE   = newDaoModelType("sql_select_one", "Sql-查询单个", invokeDaoSqlSelectOne)
	DAO_SQL_SELECT_LIST  = newDaoModelType("sql_select_list", "Sql-查询列表", invokeDaoSqlSelectList)
	DAO_SQL_SELECT_PAGE  = newDaoModelType("sql_select_page", "Sql-分页查询", invokeDaoSqlSelectPage)
	DAO_SQL_SELECT_COUNT = newDaoModelType("sql_select_count", "Sql-统计查询", invokeDaoSqlSelectCount)
	DAO_SQL_INSERT       = newDaoModelType("sql_insert", "Sql-新增", invokeDaoSqlInsert)
	DAO_SQL_UPDATE       = newDaoModelType("sql_update", "Sql-更新", invokeDaoSqlUpdate)
	DAO_SQL_DELETE       = newDaoModelType("sql_delete", "Sql-删除", invokeDaoSqlDelete)

	DAO_HTTP_GET    = newDaoModelType("http_get", "Http-Get", invokeDaoHttpGet)
	DAO_HTTP_POST   = newDaoModelType("http_post", "Http-Post", invokeDaoHttpPost)
	DAO_HTTP_HEAD   = newDaoModelType("http_head", "Http-Head", invokeDaoHttpHead)
	DAO_HTTP_PUT    = newDaoModelType("http_put", "Http-Put", invokeDaoHttpPut)
	DAO_HTTP_PATCH  = newDaoModelType("http_patch", "Http-Patch", invokeDaoHttpPatch)
	DAO_HTTP_DELETE = newDaoModelType("http_delete", "Http-Delete", invokeDaoHttpDelete)
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
