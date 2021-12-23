package modelcoder

type DaoSqlSelectOneModel struct {
	Name       string           `json:"name,omitempty"`       // 名称，同一个应用中唯一
	Type       string           `json:"type,omitempty"`       // 类型
	Parameters []*VariableModel `json:"parameters,omitempty"` // 参数配置
	Result     *VariableModel   `json:"result,omitempty"`     // 结果配置
}

func (this_ *DaoSqlSelectOneModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectOneModel) GetType() *DaoModelType {
	return DAO_SQL_SELECT_ONE
}

func (this_ *DaoSqlSelectOneModel) GetParameters() []*VariableModel {
	return this_.Parameters
}

func (this_ *DaoSqlSelectOneModel) GetResult() *VariableModel {
	return this_.Result
}

type DaoSqlSelectListModel struct {
	Name       string           `json:"name,omitempty"`       // 名称，同一个应用中唯一
	Type       string           `json:"type,omitempty"`       // 类型
	Parameters []*VariableModel `json:"parameters,omitempty"` // 参数配置
	Result     *VariableModel   `json:"result,omitempty"`     // 结果配置
}

func (this_ *DaoSqlSelectListModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectListModel) GetType() *DaoModelType {
	return DAO_SQL_SELECT_LIST
}

func (this_ *DaoSqlSelectListModel) GetParameters() []*VariableModel {
	return this_.Parameters
}

func (this_ *DaoSqlSelectListModel) GetResult() *VariableModel {
	return this_.Result
}

type DaoSqlSelectPageModel struct {
	Name       string           `json:"name,omitempty"`       // 名称，同一个应用中唯一
	Type       string           `json:"type,omitempty"`       // 类型
	Parameters []*VariableModel `json:"parameters,omitempty"` // 参数配置
	Result     *VariableModel   `json:"result,omitempty"`     // 结果配置
}

func (this_ *DaoSqlSelectPageModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectPageModel) GetType() *DaoModelType {
	return DAO_SQL_SELECT_PAGE
}

func (this_ *DaoSqlSelectPageModel) GetParameters() []*VariableModel {
	return this_.Parameters
}

func (this_ *DaoSqlSelectPageModel) GetResult() *VariableModel {
	return this_.Result
}

type DaoSqlSelectCountModel struct {
	Name       string           `json:"name,omitempty"`       // 名称，同一个应用中唯一
	Type       string           `json:"type,omitempty"`       // 类型
	Parameters []*VariableModel `json:"parameters,omitempty"` // 参数配置
	Result     *VariableModel   `json:"result,omitempty"`     // 结果配置
}

func (this_ *DaoSqlSelectCountModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlSelectCountModel) GetType() *DaoModelType {
	return DAO_SQL_SELECT_COUNT
}

type DaoSqlInsertModel struct {
	Name       string                `json:"name,omitempty"`       // 名称，同一个应用中唯一
	Type       string                `json:"type,omitempty"`       // 类型
	Database   string                `json:"database,omitempty"`   // 库名
	Table      string                `json:"table,omitempty"`      // 表名
	Columns    []*DaoSqlInsertColumn `json:"columns,omitempty"`    // 新增字段
	Parameters []*VariableModel      `json:"parameters,omitempty"` // 参数配置
	Result     *VariableModel        `json:"result,omitempty"`     // 结果配置
}

type DaoSqlInsertColumn struct {
	IfScript      string `json:"ifScript,omitempty"`      // 条件  满足该条件 添加
	Name          string `json:"name,omitempty"`          // 字段名称
	ValueScript   string `json:"valueScript,omitempty"`   // 字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值
	Required      bool   `json:"required,omitempty"`      // 必填
	AutoIncrement bool   `json:"autoIncrement,omitempty"` // 自增列
	AllowEmpty    bool   `json:"allowEmpty,omitempty"`    // 允许空值，如果是null或空字符串则也设置值
}

func (this_ *DaoSqlInsertModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlInsertModel) GetType() *DaoModelType {
	return DAO_SQL_INSERT
}

func (this_ *DaoSqlInsertModel) GetParameters() []*VariableModel {
	return this_.Parameters
}

func (this_ *DaoSqlInsertModel) GetResult() *VariableModel {
	return this_.Result
}

type DaoSqlUpdateModel struct {
	Name       string           `json:"name,omitempty"`       // 名称，同一个应用中唯一
	Type       string           `json:"type,omitempty"`       // 类型
	Parameters []*VariableModel `json:"parameters,omitempty"` // 参数配置
	Result     *VariableModel   `json:"result,omitempty"`     // 结果配置
}

func (this_ *DaoSqlUpdateModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlUpdateModel) GetType() *DaoModelType {
	return DAO_SQL_UPDATE
}

func (this_ *DaoSqlUpdateModel) GetParameters() []*VariableModel {
	return this_.Parameters
}

func (this_ *DaoSqlUpdateModel) GetResult() *VariableModel {
	return this_.Result
}

type DaoSqlDeleteModel struct {
	Name       string           `json:"name,omitempty"`       // 名称，同一个应用中唯一
	Type       string           `json:"type,omitempty"`       // 类型
	Parameters []*VariableModel `json:"parameters,omitempty"` // 参数配置
	Result     *VariableModel   `json:"result,omitempty"`     // 结果配置
}

func (this_ *DaoSqlDeleteModel) GetName() string {
	return this_.Name
}

func (this_ *DaoSqlDeleteModel) GetType() *DaoModelType {
	return DAO_SQL_DELETE
}

func (this_ *DaoSqlDeleteModel) GetParameters() []*VariableModel {
	return this_.Parameters
}

func (this_ *DaoSqlDeleteModel) GetResult() *VariableModel {
	return this_.Result
}
