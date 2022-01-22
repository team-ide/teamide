package model

type SqlSelect struct {
	Database     string                `json:"database,omitempty" yaml:"database,omitempty"`         // 库名
	Table        string                `json:"table,omitempty" yaml:"table,omitempty"`               // 表名
	Alias        string                `json:"alias,omitempty" yaml:"alias,omitempty"`               // 表别名
	Distinct     bool                  `json:"distinct,omitempty" yaml:"distinct,omitempty"`         // SELECT DISTINCT 语句  列出不同的值
	SelectCount  bool                  `json:"selectCount,omitempty" yaml:"selectCount,omitempty"`   // 是统计查询
	SelectOne    bool                  `json:"selectOne,omitempty" yaml:"selectOne,omitempty"`       // 是查询单个
	SelectPage   bool                  `json:"selectPage,omitempty" yaml:"selectPage,omitempty"`     // 是分页查询
	Columns      []*SqlSelectColumn    `json:"columns,omitempty" yaml:"columns,omitempty"`           // 字段
	Intos        []*SqlSelectInto      `json:"intos,omitempty" yaml:"intos,omitempty"`               // SELECT INTO 复制表
	InnerJoins   []*SqlSelectInnerJoin `json:"innerJoins,omitempty" yaml:"innerJoins,omitempty"`     // INNER JOIN
	LeftJoin     []*SqlSelectLeftJoin  `json:"leftJoin,omitempty" yaml:"leftJoin,omitempty"`         // LEFT JOIN
	RightJoin    []*SqlSelectRightJoin `json:"rightJoin,omitempty" yaml:"rightJoin,omitempty"`       // RIGHT JOIN
	FullJoin     []*SqlSelectFullJoin  `json:"fullJoin,omitempty" yaml:"fullJoin,omitempty"`         // FULL JOIN
	Wheres       []*SqlWhere           `json:"wheres,omitempty" yaml:"wheres,omitempty"`             // 条件
	Orders       []*SqlSelectOrder     `json:"orders,omitempty" yaml:"orders,omitempty"`             // 排序
	Groups       []*SqlSelectGroup     `json:"groups,omitempty" yaml:"groups,omitempty"`             // 分组
	Havings      []*SqlSelectHaving    `json:"havings,omitempty" yaml:"havings,omitempty"`           // 筛选分组后的各组数据
	UnionSelects []*SqlSelect          `json:"unionSelects,omitempty" yaml:"unionSelects,omitempty"` // UNION 操作符 合并两个或多个 SELECT 语句的结果集
}

type SqlSelectColumn struct {
	IfScript   string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`     // 条件  满足该条件 添加
	Custom     bool   `json:"custom,omitempty" yaml:"custom,omitempty"`         // 是否自定义
	CustomSql  string `json:"customSql,omitempty" yaml:"customSql,omitempty"`   // 是否自定义
	TableAlias string `json:"tableAlias,omitempty" yaml:"tableAlias,omitempty"` // 表别名
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`             // 字段名称
	Alias      string `json:"alias,omitempty" yaml:"alias,omitempty"`           // 字段别名
}

type ServiceStepSqlSelect struct {
	Base *ServiceStepBase

	SqlSelect        *SqlSelect `json:"sqlSelect,omitempty" yaml:"sqlSelect,omitempty"`               // 执行 SQL SELECT 操作
	VariableName     string     `json:"variableName,omitempty" yaml:"variableName,omitempty"`         // 变量名称
	VariableDataType string     `json:"variableDataType,omitempty" yaml:"variableDataType,omitempty"` // 变量数据类型
}

func (this_ *ServiceStepSqlSelect) GetBase() *ServiceStepBase {
	return this_.Base
}

func (this_ *ServiceStepSqlSelect) SetBase(v *ServiceStepBase) {
	this_.Base = v
}

type SqlSelectInto struct {
	IfScript  string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`   // 条件  满足该条件 添加
	Custom    bool   `json:"custom,omitempty" yaml:"custom,omitempty"`       // 是否自定义
	CustomSql string `json:"customSql,omitempty" yaml:"customSql,omitempty"` // 是否自定义
	Table     string `json:"table,omitempty" yaml:"table,omitempty"`         // 表名
}

type SqlSelectInnerJoin struct {
	IfScript  string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`   // 条件  满足该条件 添加
	Custom    bool   `json:"custom,omitempty" yaml:"custom,omitempty"`       // 是否自定义
	CustomSql string `json:"customSql,omitempty" yaml:"customSql,omitempty"` // 是否自定义
	Table     string `json:"table,omitempty" yaml:"table,omitempty"`         // 表名
	Alias     string `json:"alias,omitempty" yaml:"alias,omitempty"`         // 别名
	On        string `json:"on,omitempty" yaml:"on,omitempty"`               // ON
}

type SqlSelectLeftJoin struct {
	IfScript  string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`   // 条件  满足该条件 添加
	Custom    bool   `json:"custom,omitempty" yaml:"custom,omitempty"`       // 是否自定义
	CustomSql string `json:"customSql,omitempty" yaml:"customSql,omitempty"` // 是否自定义
	Table     string `json:"table,omitempty" yaml:"table,omitempty"`         // 表名
	Alias     string `json:"alias,omitempty" yaml:"alias,omitempty"`         // 别名
	On        string `json:"on,omitempty" yaml:"on,omitempty"`               // ON
}

type SqlSelectRightJoin struct {
	IfScript  string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`   // 条件  满足该条件 添加
	Custom    bool   `json:"custom,omitempty" yaml:"custom,omitempty"`       // 是否自定义
	CustomSql string `json:"customSql,omitempty" yaml:"customSql,omitempty"` // 是否自定义
	Table     string `json:"table,omitempty" yaml:"table,omitempty"`         // 表名
	Alias     string `json:"alias,omitempty" yaml:"alias,omitempty"`         // 别名
	On        string `json:"on,omitempty" yaml:"on,omitempty"`               // ON
}

type SqlSelectFullJoin struct {
	IfScript  string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`   // 条件  满足该条件 添加
	Custom    bool   `json:"custom,omitempty" yaml:"custom,omitempty"`       // 是否自定义
	CustomSql string `json:"customSql,omitempty" yaml:"customSql,omitempty"` // 是否自定义
	Table     string `json:"table,omitempty" yaml:"table,omitempty"`         // 表名
	Alias     string `json:"alias,omitempty" yaml:"alias,omitempty"`         // 别名
	On        string `json:"on,omitempty" yaml:"on,omitempty"`               // ON
}

type SqlSelectOrder struct {
	IfScript   string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`     // 条件  满足该条件 添加
	Custom     bool   `json:"custom,omitempty" yaml:"custom,omitempty"`         // 是否自定义
	CustomSql  string `json:"customSql,omitempty" yaml:"customSql,omitempty"`   // 是否自定义
	TableAlias string `json:"tableAlias,omitempty" yaml:"tableAlias,omitempty"` // 表别名
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`             // 字段名称
	AscDesc    string `json:"ascDesc,omitempty" yaml:"ascDesc,omitempty"`       // 默认按照升序 升序ASC 降序DESC
}

type SqlSelectGroup struct {
	IfScript   string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`     // 条件  满足该条件 添加
	Custom     bool   `json:"custom,omitempty" yaml:"custom,omitempty"`         // 是否自定义
	CustomSql  string `json:"customSql,omitempty" yaml:"customSql,omitempty"`   // 是否自定义
	TableAlias string `json:"tableAlias,omitempty" yaml:"tableAlias,omitempty"` // 表别名
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`             // 字段名称
}

type SqlSelectHaving struct {
	IfScript  string `json:"ifScript,omitempty" yaml:"ifScript,omitempty"`   // 条件  满足该条件 添加
	Custom    bool   `json:"custom,omitempty" yaml:"custom,omitempty"`       // 是否自定义
	CustomSql string `json:"customSql,omitempty" yaml:"customSql,omitempty"` // 是否自定义
	Having    string `json:"having,omitempty" yaml:"having,omitempty"`       // 筛选分组后的各组数据
}
