package db

import (
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"strings"
	"teamide/pkg/base"
)

type modelGen struct {
	modelType string
	tables    []*dialect.TableModel
	content   string
}

func (this_ *modelGen) gen() (content string, err error) {
	this_.content = ""

	for _, table := range this_.tables {
		err = this_.append(table)
		if err != nil {
			return
		}
		this_.content += "\n"
	}

	content = this_.content
	return
}

func (this_ *modelGen) append(table *dialect.TableModel) (err error) {
	switch this_.modelType {
	case "table":
		err = this_.appendTableModel(table)
		break

	}
	return
}
func (this_ *modelGen) appendTableModel(table *dialect.TableModel) (err error) {
	name := strings.ToLower(table.TableName)
	name = base.Marshal(name)
	this_.content += "" + name + "Model = &dialect.TableModel{" + "\n"
	this_.content += "\t" + "TableName: `" + table.TableName + "`," + "\n"
	this_.content += "\t" + "TableComment: `" + table.TableComment + "`," + "\n"
	this_.content += "\t" + "ColumnList: []*dialect.ColumnModel{" + "\n"

	for _, column := range table.ColumnList {
		this_.content += "\t\t{"
		this_.content += "ColumnName: `" + column.ColumnName + "`, "
		this_.content += "ColumnDataType: `" + column.ColumnDataType + "`, "
		if column.ColumnLength > 0 {
			this_.content += "ColumnLength: " + fmt.Sprintf("%d", column.ColumnLength) + ", "
		}
		if column.ColumnPrecision > 0 {
			this_.content += "ColumnPrecision: " + fmt.Sprintf("%d", column.ColumnPrecision) + ", "
		}
		if column.ColumnScale > 0 {
			this_.content += "ColumnScale: " + fmt.Sprintf("%d", column.ColumnScale) + ", "
		}
		if column.ColumnNotNull {
			this_.content += "ColumnNotNull: true, "
		}
		if column.PrimaryKey {
			this_.content += "PrimaryKey: true, "
		}
		this_.content += "ColumnComment: `" + column.ColumnComment + "`"
		this_.content += "}," + "\n"
	}

	this_.content += "\t" + "}," + "\n"
	this_.content += "\t" + "IndexList: []*dialect.IndexModel{" + "\n"
	for _, index := range table.IndexList {
		this_.content += "\t\t{"
		this_.content += "IndexName: `" + index.IndexName + "`, "
		this_.content += "IndexType: `" + index.IndexType + "`, "
		this_.content += "ColumnNames: []string{\"" + strings.Join(index.ColumnNames, "\", \"") + "\"}, "
		this_.content += "IndexComment: `" + index.IndexComment + "`"
		this_.content += "}," + "\n"
	}
	this_.content += "\t" + "}," + "\n"
	this_.content += "}" + "\n"
	return
}
