package modelers

import (
	"strings"
)

type StepDbModel struct {
	*StepModel `json:",inline"`

	Db       string `json:"db,omitempty"`       // 数据库操作
	Database string `json:"database,omitempty"` // 库名
	Sql      string `json:"sql,omitempty"`      // 自定义SQL
	Table    string `json:"table,omitempty"`    // 表
	Alias    string `json:"alias,omitempty"`    // 表别名
	Distinct bool   `json:"distinct,omitempty"` // SELECT DISTINCT 语句  列出不同的值

	Columns   []*StepDbColumn    `json:"columns,omitempty"`   // 字段
	Into      []*StepDbInto      `json:"into,omitempty"`      // SELECT INTO 复制表
	InnerJoin []*StepDbInnerJoin `json:"innerJoin,omitempty"` // INNER JOIN
	LeftJoin  []*StepDbLeftJoin  `json:"leftJoin,omitempty"`  // LEFT JOIN
	RightJoin []*StepDbRightJoin `json:"rightJoin,omitempty"` // RIGHT JOIN
	FullJoin  []*StepDbFullJoin  `json:"fullJoin,omitempty"`  // FULL JOIN
	Wheres    []*StepDbWhere     `json:"wheres,omitempty"`    // WHERE
	Order     []*StepDbOrder     `json:"order,omitempty"`     // ORDER BY
	Group     []*StepDbGroup     `json:"group,omitempty"`     // GROUP BY
	Having    []*StepDbHaving    `json:"having,omitempty"`    // HAVING 筛选分组后的各组数据
	Union     []*StepDbModel     `json:"union,omitempty"`     // UNION 操作符 合并两个或多个 SELECT 语句的结果集

	SetVar     string `json:"setVar,omitempty"`
	SetVarType string `json:"setVarType,omitempty"`
}

func (this_ *StepDbModel) IsSelect() bool {
	return strings.Contains(this_.Db, "select")
}

func (this_ *StepDbModel) IsSelectOne() bool {
	return this_.Db == "selectOne"
}

func (this_ *StepDbModel) IsSelectPage() bool {
	return this_.Db == "selectPage"
}

func (this_ *StepDbModel) GetType() *StepDbType {
	for _, one := range StepDbTypes {
		if strings.EqualFold(one.Value, this_.Db) {
			return one
		}
	}
	return nil
}

type StepDbType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	StepDbTypes []*StepDbType
	DbGet       = appendStepDbType("get", "")
)

func appendStepDbType(value string, text string) *StepDbType {
	res := &StepDbType{
		Value: value,
		Text:  text,
	}
	StepDbTypes = append(StepDbTypes, res)
	return res
}

var (
	docTemplateStepDbName          = "step_db"
	docTemplateStepDbColumnName    = "step_db_column"
	docTemplateStepDbIntoName      = "step_db_into"
	docTemplateStepDbInnerJoinName = "step_db_innerJoin"
	docTemplateStepDbLeftJoinName  = "step_db_leftJoin"
	docTemplateStepDbRightJoinName = "step_db_rightJoin"
	docTemplateStepDbFullJoinName  = "step_db_fullJoin"
	docTemplateStepDbWhereName     = "step_db_where"
	docTemplateStepDbOrderName     = "step_db_order"
	docTemplateStepDbGroupName     = "step_db_group"
	docTemplateStepDbHavingName    = "step_db_having"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepDbName,
		Fields: []*docTemplateField{
			{Name: "db", Comment: "数据库操作"},
			{Name: "database", Comment: "数据库"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "table", Comment: "表"},
			{Name: "alias", Comment: "别名"},
			{Name: "distinct", Comment: "去重"},

			{Name: "columns", Comment: "字段", IsList: true, StructName: docTemplateStepDbColumnName},
			{Name: "into", Comment: "SELECT INTO", IsList: true, StructName: docTemplateStepDbIntoName},
			{Name: "innerJoin", Comment: "INNER JOIN", IsList: true, StructName: docTemplateStepDbInnerJoinName},
			{Name: "leftJoin", Comment: "LEFT JOIN", IsList: true, StructName: docTemplateStepDbLeftJoinName},
			{Name: "rightJoin", Comment: "RIGHT JOIN", IsList: true, StructName: docTemplateStepDbRightJoinName},
			{Name: "fullJoin", Comment: "FULL JOIN", IsList: true, StructName: docTemplateStepDbFullJoinName},
			{Name: "wheres", Comment: "WHERE", IsList: true, StructName: docTemplateStepDbWhereName},
			{Name: "order", Comment: "ORDER", IsList: true, StructName: docTemplateStepDbOrderName},
			{Name: "group", Comment: "GROUP", IsList: true, StructName: docTemplateStepDbGroupName},
			{Name: "having", Comment: "HAVING", IsList: true, StructName: docTemplateStepDbHavingName},
			{Name: "union", Comment: "UNION", IsList: true, StructName: docTemplateStepDbName},

			{Name: "setVar", Comment: "设置变量"},
			{Name: "setVarType", Comment: "设置变量类型"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepDbModel{}
		},
		newModels: func() interface{} {
			var vs []*StepDbModel
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbModel)
			vs = append(vs, value.(*StepDbModel))
			return vs
		},
	})
}

type StepDbColumn struct {
	If            string `json:"if,omitempty"`            // 条件  满足该条件 添加
	Custom        bool   `json:"custom,omitempty"`        // 是否自定义
	Sql           string `json:"sql,omitempty"`           // 自定义SQL
	TableAlias    string `json:"tableAlias,omitempty"`    // 表别名
	Name          string `json:"name,omitempty"`          // 字段名称
	AutoIncrement bool   `json:"autoIncrement,omitempty"` // 自增列
	Alias         string `json:"alias,omitempty"`         // 字段别名
	Value         string `json:"value,omitempty"`         // 字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值
	Required      bool   `json:"required,omitempty"`      // 必填
	IgnoreEmpty   bool   `json:"ignoreEmpty,omitempty"`   // 忽略空值，如果忽略，则值是null、空字符串、0不设值

}

func init() {
	addDocTemplate(&docTemplate{
		Name:         docTemplateStepDbColumnName,
		Abbreviation: "name",
		Fields: []*docTemplateField{
			{Name: "if", Comment: "条件"},
			{Name: "custom", Comment: "是否自定义"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "tableAlias", Comment: "表别名"},
			{Name: "name", Comment: "字段名称"},
			{Name: "autoIncrement", Comment: "自增列"},
			{Name: "alias", Comment: "字段别名"},
			{Name: "value", Comment: "字段值，可以是属性名、表达式等，如果该值为空，自动取名称相同的值"},
			{Name: "required", Comment: "必填"},
			{Name: "ignoreEmpty", Comment: "忽略空值，如果忽略，则值是null、空字符串、0不设值"},
		},
		newModel: func() interface{} {
			return &StepDbColumn{}
		},
		newModels: func() interface{} {
			var vs []*StepDbColumn
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbColumn)
			vs = append(vs, value.(*StepDbColumn))
			return vs
		},
	})
}

type StepDbInto struct {
	If       string `json:"if,omitempty"`       // 条件  满足该条件 添加
	Custom   bool   `json:"custom,omitempty"`   // 是否自定义
	Sql      string `json:"sql,omitempty"`      // 自定义SQL
	Database string `json:"database,omitempty"` // 库名称
	Table    string `json:"table,omitempty"`    // 表名
}

func init() {
	addDocTemplate(&docTemplate{
		Name:         docTemplateStepDbIntoName,
		Abbreviation: "table",
		Fields: []*docTemplateField{
			{Name: "if", Comment: "条件"},
			{Name: "custom", Comment: "是否自定义"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "database", Comment: "库名称"},
			{Name: "table", Comment: "表名"},
		},
		newModel: func() interface{} {
			return &StepDbInto{}
		},
		newModels: func() interface{} {
			var vs []*StepDbInto
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbInto)
			vs = append(vs, value.(*StepDbInto))
			return vs
		},
	})
}

type StepDbInnerJoin struct {
	If     string `json:"if,omitempty"`     // 条件  满足该条件 添加
	Custom bool   `json:"custom,omitempty"` // 是否自定义
	Sql    string `json:"sql,omitempty"`    // 自定义SQL
	Table  string `json:"table,omitempty"`  // 表名
	Alias  string `json:"alias,omitempty"`  // 别名
	On     string `json:"on,omitempty"`     // ON
}

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepDbInnerJoinName,
		Fields: []*docTemplateField{
			{Name: "if", Comment: "条件"},
			{Name: "custom", Comment: "是否自定义"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "table", Comment: "表名"},
			{Name: "alias", Comment: "别名"},
			{Name: "on", Comment: "条件"},
		},
		newModel: func() interface{} {
			return &StepDbInnerJoin{}
		},
		newModels: func() interface{} {
			var vs []*StepDbInnerJoin
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbInnerJoin)
			vs = append(vs, value.(*StepDbInnerJoin))
			return vs
		},
	})
}

type StepDbLeftJoin struct {
	If     string `json:"if,omitempty"`     // 条件  满足该条件 添加
	Custom bool   `json:"custom,omitempty"` // 是否自定义
	Sql    string `json:"sql,omitempty"`    // 自定义SQL
	Table  string `json:"table,omitempty"`  // 表名
	Alias  string `json:"alias,omitempty"`  // 别名
	On     string `json:"on,omitempty"`     // ON
}

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepDbLeftJoinName,
		Fields: []*docTemplateField{
			{Name: "if", Comment: "条件"},
			{Name: "custom", Comment: "是否自定义"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "table", Comment: "表名"},
			{Name: "alias", Comment: "别名"},
			{Name: "on", Comment: "条件"},
		},
		newModel: func() interface{} {
			return &StepDbLeftJoin{}
		},
		newModels: func() interface{} {
			var vs []*StepDbLeftJoin
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbLeftJoin)
			vs = append(vs, value.(*StepDbLeftJoin))
			return vs
		},
	})
}

type StepDbRightJoin struct {
	If     string `json:"if,omitempty"`     // 条件  满足该条件 添加
	Custom bool   `json:"custom,omitempty"` // 是否自定义
	Sql    string `json:"sql,omitempty"`    // 自定义SQL
	Table  string `json:"table,omitempty"`  // 表名
	Alias  string `json:"alias,omitempty"`  // 别名
	On     string `json:"on,omitempty"`     // ON
}

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepDbRightJoinName,
		Fields: []*docTemplateField{
			{Name: "if", Comment: "条件"},
			{Name: "custom", Comment: "是否自定义"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "table", Comment: "表名"},
			{Name: "alias", Comment: "别名"},
			{Name: "on", Comment: "条件"},
		},
		newModel: func() interface{} {
			return &StepDbRightJoin{}
		},
		newModels: func() interface{} {
			var vs []*StepDbRightJoin
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbRightJoin)
			vs = append(vs, value.(*StepDbRightJoin))
			return vs
		},
	})
}

type StepDbFullJoin struct {
	If     string `json:"if,omitempty"`     // 条件  满足该条件 添加
	Custom bool   `json:"custom,omitempty"` // 是否自定义
	Sql    string `json:"sql,omitempty"`    // 自定义SQL
	Table  string `json:"table,omitempty"`  // 表名
	Alias  string `json:"alias,omitempty"`  // 别名
	On     string `json:"on,omitempty"`     // ON
}

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepDbFullJoinName,
		Fields: []*docTemplateField{
			{Name: "if", Comment: "条件"},
			{Name: "custom", Comment: "是否自定义"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "table", Comment: "表名"},
			{Name: "alias", Comment: "别名"},
			{Name: "on", Comment: "条件"},
		},
		newModel: func() interface{} {
			return &StepDbFullJoin{}
		},
		newModels: func() interface{} {
			var vs []*StepDbFullJoin
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbFullJoin)
			vs = append(vs, value.(*StepDbFullJoin))
			return vs
		},
	})
}

type StepDbOrder struct {
	If         string `json:"if,omitempty"`         // 条件  满足该条件 添加
	Custom     bool   `json:"custom,omitempty"`     // 是否自定义
	Sql        string `json:"sql,omitempty"`        // 自定义SQL
	TableAlias string `json:"tableAlias,omitempty"` // 表别名
	Name       string `json:"name,omitempty"`       // 字段名称
	Asc        bool   `json:"asc,omitempty"`        // 升序  默认按照升序 升序ASC 降序DESC
	Desc       bool   `json:"desc,omitempty"`       // 降序
}

func (this_ *StepDbOrder) GetAscDesc() string {
	if this_.Desc {
		return "DESC"
	}
	return "ASC"
}

func init() {
	addDocTemplate(&docTemplate{
		Name:         docTemplateStepDbOrderName,
		Abbreviation: "name",
		Fields: []*docTemplateField{
			{Name: "if", Comment: "条件"},
			{Name: "custom", Comment: "是否自定义"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "tableAlias", Comment: "表别名"},
			{Name: "name", Comment: "字段名称"},
			{Name: "asc", Comment: "升序 默认"},
			{Name: "desc", Comment: "降序"},
		},
		newModel: func() interface{} {
			return &StepDbOrder{}
		},
		newModels: func() interface{} {
			var vs []*StepDbOrder
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbOrder)
			vs = append(vs, value.(*StepDbOrder))
			return vs
		},
	})
}

type StepDbGroup struct {
	If         string `json:"if,omitempty"`         // 条件  满足该条件 添加
	Custom     bool   `json:"custom,omitempty"`     // 是否自定义
	Sql        string `json:"sql,omitempty"`        // 自定义SQL
	TableAlias string `json:"tableAlias,omitempty"` // 表别名
	Name       string `json:"name,omitempty"`       // 字段名称
}

func init() {
	addDocTemplate(&docTemplate{
		Name:         docTemplateStepDbGroupName,
		Abbreviation: "name",
		Fields: []*docTemplateField{
			{Name: "if", Comment: "条件"},
			{Name: "custom", Comment: "是否自定义"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "tableAlias", Comment: "表别名"},
			{Name: "name", Comment: "字段名称"},
		},
		newModel: func() interface{} {
			return &StepDbGroup{}
		},
		newModels: func() interface{} {
			var vs []*StepDbGroup
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbGroup)
			vs = append(vs, value.(*StepDbGroup))
			return vs
		},
	})
}

type StepDbHaving struct {
	If     string `json:"if,omitempty"`     // 条件  满足该条件 添加
	Custom bool   `json:"custom,omitempty"` // 是否自定义
	Sql    string `json:"sql,omitempty"`    // 自定义SQL
	Having string `json:"having,omitempty"` // 筛选分组后的各组数据
}

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepDbHavingName,
		Fields: []*docTemplateField{
			{Name: "if", Comment: "条件"},
			{Name: "custom", Comment: "是否自定义"},
			{Name: "sql", Comment: "自定义SQL"},
			{Name: "having", Comment: "筛选分组后的各组数据"},
		},
		newModel: func() interface{} {
			return &StepDbHaving{}
		},
		newModels: func() interface{} {
			var vs []*StepDbHaving
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StepDbHaving)
			vs = append(vs, value.(*StepDbHaving))
			return vs
		},
	})
}
