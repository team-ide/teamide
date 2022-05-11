package db

type DatabaseModel struct {
	Name    string `json:"name" column:"name"`
	Comment string `json:"comment" column:"comment"`
}

type TableModel struct {
	Name       string              `json:"name" column:"name"`
	Comment    string              `json:"comment" column:"comment"`
	ColumnList []*TableColumnModel `json:"columnList"`
	IndexList  []*TableIndexModel  `json:"indexList"`
}

type TableColumnModel struct {
	Name       string `json:"name" column:"name"`
	Comment    string `json:"comment" column:"comment"`
	Type       string `json:"type" column:"type"`
	Length     int    `json:"length"`
	Decimal    int    `json:"decimal"`
	PrimaryKey bool   `json:"primaryKey"`
	NotNull    bool   `json:"notNull"`
	Default    string `json:"default" column:"default"`
	ISNullable string `json:"-" column:"IS_NULLABLE"`
}

type TableIndexModel struct {
	Name      string `json:"name" column:"name"`
	Type      string `json:"type" column:"type"`
	Columns   string `json:"columns" column:"columns"`
	Comment   string `json:"comment" column:"comment"`
	NONUnique string `json:"-" column:"NON_UNIQUE"`
}
