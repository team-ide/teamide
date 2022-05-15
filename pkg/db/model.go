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
	OldName    string              `json:"oldName" column:"oldName"`
	OldComment string              `json:"oldComment" column:"oldComment"`
}

func (this_ *TableModel) FindColumnByName(name string) *TableColumnModel {
	if len(this_.ColumnList) > 0 {
		for _, one := range this_.ColumnList {
			if one.Name == name {
				return one
			}
		}
	}
	return nil
}
func (this_ *TableModel) FindColumnByOldName(oldName string) *TableColumnModel {
	if len(this_.ColumnList) > 0 {
		for _, one := range this_.ColumnList {
			if one.OldName == oldName {
				return one
			}
		}
	}
	return nil
}

func (this_ *TableModel) FindIndexByName(name string) *TableIndexModel {
	if len(this_.IndexList) > 0 {
		for _, one := range this_.IndexList {
			if one.Name == name {
				return one
			}
		}
	}
	return nil
}
func (this_ *TableModel) FindIndexByOldName(oldName string) *TableIndexModel {
	if len(this_.IndexList) > 0 {
		for _, one := range this_.IndexList {
			if one.OldName == oldName {
				return one
			}
		}
	}
	return nil
}

type TableColumnModel struct {
	Name          string      `json:"name" column:"name"`
	Comment       string      `json:"comment" column:"comment"`
	Type          string      `json:"type" column:"type"`
	Length        int         `json:"length"`
	Decimal       int         `json:"decimal"`
	PrimaryKey    bool        `json:"primaryKey"`
	NotNull       bool        `json:"notNull"`
	Default       interface{} `json:"default" column:"default"`
	ISNullable    string      `json:"-" column:"IS_NULLABLE"`
	OldName       string      `json:"oldName" column:"oldName"`
	OldComment    string      `json:"oldComment" column:"oldComment"`
	OldType       string      `json:"oldType" column:"oldType"`
	OldLength     int         `json:"oldLength" column:"oldLength"`
	OldDecimal    int         `json:"oldDecimal" column:"oldDecimal"`
	OldPrimaryKey bool        `json:"oldPrimaryKey" column:"oldPrimaryKey"`
	OldNotNull    bool        `json:"oldNotNull" column:"oldNotNull"`
	OldDefault    interface{} `json:"oldDefault" column:"oldDefault"`
	BeforeColumn  string      `json:"beforeColumn" column:"beforeColumn"`
	Deleted       bool        `json:"deleted" column:"deleted"`
}

type TableIndexModel struct {
	Name       string   `json:"name" column:"name"`
	Type       string   `json:"type" column:"type"`
	Columns    []string `json:"columns" column:"columns"`
	Comment    string   `json:"comment" column:"comment"`
	NONUnique  string   `json:"-" column:"NON_UNIQUE"`
	COLUMNName string   `json:"-" column:"COLUMN_NAME"`
	OldName    string   `json:"oldName" column:"oldName"`
	OldComment string   `json:"oldComment" column:"oldComment"`
	OldType    string   `json:"oldType" column:"oldType"`
	OldColumns []string `json:"oldColumns" column:"oldColumns"`
	Deleted    bool     `json:"deleted" column:"deleted"`
}
