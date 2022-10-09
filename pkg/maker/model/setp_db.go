package model

type StepDbModel struct {
	*StepModel

	Db        string `json:"db,omitempty"`        // 数据库操作
	Table     string `json:"table,omitempty"`     // 表
	Columns   string `json:"columns,omitempty"`   // 字段
	Wheres    string `json:"wheres,omitempty"`    // 条件
	LeftJoins string `json:"leftJoins,omitempty"` // 左连接查询
	Orders    string `json:"orders,omitempty"`    // 排序
	Groups    string `json:"groups,omitempty"`    // 分组
}
