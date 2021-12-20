package modelcoder

type DaoSqlSelectOne struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectOne) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectOne) GetType() *DaoModelType {
	return DAO_SQL_SELECT_ONE
}

type DaoSqlSelectList struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectList) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectList) GetType() *DaoModelType {
	return DAO_SQL_SELECT_LIST
}

type DaoSqlSelectPage struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectPage) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectPage) GetType() *DaoModelType {
	return DAO_SQL_SELECT_PAGE
}

type DaoSqlSelectCount struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectCount) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectCount) GetType() *DaoModelType {
	return DAO_SQL_SELECT_COUNT
}

type DaoSqlSelectInsert struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectInsert) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectInsert) GetType() *DaoModelType {
	return DAO_SQL_INSERT
}

type DaoSqlSelectUpdate struct {
	Name string `json:"name,omitempty"` // 名称，同一个应用中唯一
	Type string `json:"type,omitempty"` // 类型
}

func (this_ *DaoSqlSelectUpdate) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectUpdate) GetType() *DaoModelType {
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
