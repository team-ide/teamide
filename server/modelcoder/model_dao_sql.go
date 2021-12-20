package modelcoder

type DaoSqlSelectOneModel struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectOneModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectOneModel) GetType() *DaoModelType {
	return DAO_SQL_SELECT_ONE
}

type DaoSqlSelectListModel struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectListModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectListModel) GetType() *DaoModelType {
	return DAO_SQL_SELECT_LIST
}

type DaoSqlSelectPageModel struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectPageModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectPageModel) GetType() *DaoModelType {
	return DAO_SQL_SELECT_PAGE
}

type DaoSqlSelectCountModel struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectCountModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectCountModel) GetType() *DaoModelType {
	return DAO_SQL_SELECT_COUNT
}

type DaoSqlSelectInsertModel struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectInsertModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectInsertModel) GetType() *DaoModelType {
	return DAO_SQL_INSERT
}

type DaDaoSqlSelectUpdateModel struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaDaoSqlSelectUpdateModel) GetName() string {
	return this_.Name
}

func (this_ *DaDaoSqlSelectUpdateModel) GetType() *DaoModelType {
	return DAO_SQL_UPDATE
}

type DaoSqlSelectDelete struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectDelete) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectDelete) GetType() *DaoModelType {
	return DAO_SQL_DELETE
}
