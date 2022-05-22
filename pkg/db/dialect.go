package db

import (
	"errors"
	"strings"
)

type Dialect interface {
	GetDatabaseType() (databaseType *DatabaseType)
	GetColumnTypeInfo(name string) (columnTypeInfo *ColumnTypeInfo, err error)
	GetColumnType(name string, length, decimal int) (columnType string, err error)
	DatabaseDDL(param *GenerateParam, database *DatabaseModel) (sqlList []string, err error)
	DatabaseDeleteDDL(param *GenerateParam, database string) (sqlList []string, err error)
	TableDDL(param *GenerateParam, database string, table *TableModel) (sqlList []string, err error)
	TableCommentDDL(param *GenerateParam, database string, table string, comment string) (sqlList []string, err error)
	TableRenameDDL(param *GenerateParam, database string, table string, rename string) (sqlList []string, err error)
	TableDeleteDDL(param *GenerateParam, database string, table string) (sqlList []string, err error)
	TableColumnDeleteDDL(param *GenerateParam, database string, table string, column string) (sqlList []string, err error)
	TableColumnUpdateDDL(param *GenerateParam, database string, table string, column *TableColumnModel) (sqlList []string, err error)
	TableColumnAddDDL(param *GenerateParam, database string, table string, column *TableColumnModel) (sqlList []string, err error)
	TableIndexAddDDL(param *GenerateParam, database string, table string, index *TableIndexModel) (sqlList []string, err error)
	TableIndexUpdateDDL(param *GenerateParam, database string, table string, index *TableIndexModel) (sqlList []string, err error)
	TableIndexDeleteDDL(param *GenerateParam, database string, table string, index string) (sqlList []string, err error)
	TableIndexRenameDDL(param *GenerateParam, database string, table string, index string, rename string) (sqlList []string, err error)
	TablePrimaryKeyDeleteDDL(param *GenerateParam, database string, table string, primaryKeys []string) (sqlList []string, err error)
	TablePrimaryKeyAddDDL(param *GenerateParam, database string, table string, primaryKeys []string) (sqlList []string, err error)
}

func GetDialect(databaseType *DatabaseType) (dialect Dialect, err error) {
	BaseDialect := &BaseDialect{
		databaseType: databaseType,
	}
	switch databaseType {
	case DatabaseTypeMySql:
		dialect = &DatabaseMySqlDialect{
			BaseDialect: BaseDialect,
		}
		break
	case DatabaseTypeSqlite:
		dialect = &DatabaseSqliteDialect{
			BaseDialect: BaseDialect,
		}
		break
	case DatabaseTypeOracle:
		dialect = &DatabaseOracleDialect{
			BaseDialect: BaseDialect,
		}
		break
	case DatabaseTypeShenTong:
		dialect = &DatabaseShenTongDialect{
			BaseDialect: BaseDialect,
		}
		break
	case DatabaseTypeDM:
		dialect = &DatabaseDMDialect{
			BaseDialect: BaseDialect,
		}
		break
	case DatabaseTypeKingBase:
		dialect = &DatabaseKingBaseDialect{
			BaseDialect: BaseDialect,
		}
		break
	case DatabaseTypeGBase:
		dialect = &DatabaseGBaseDialect{
			BaseDialect: BaseDialect,
		}
		break
	case DatabaseTypeKunLun:
		dialect = &DatabaseKunLunDialect{
			BaseDialect: BaseDialect,
		}
		break
	case nil:
		err = errors.New("数据库类型[" + databaseType.DBType + "]暂不支持")
		break
	}
	return
}

type BaseDialect struct {
	databaseType *DatabaseType
}

func (this_ *BaseDialect) GetDatabaseType() (databaseType *DatabaseType) {
	return this_.databaseType
}

func (this_ *BaseDialect) GetColumnTypeInfo(name string) (columnTypeInfo *ColumnTypeInfo, err error) {
	databaseType := this_.GetDatabaseType()
	if databaseType == nil {
		err = errors.New("未设置数据库类型")
		return
	}
	columnTypeInfo = this_.GetDatabaseType().GetColumnTypeInfo(name)
	if columnTypeInfo == nil {
		err = errors.New("驱动[" + databaseType.DBType + "]字段类型[" + name + "]未映射!")
		return
	}
	return
}

func (this_ *BaseDialect) GetColumnType(name string, length, decimal int) (columnType string, err error) {
	columnTypeInfo, err := this_.GetColumnTypeInfo(name)
	if err != nil {
		return
	}
	columnType = columnTypeInfo.FormatColumnType(length, decimal)
	return
}

func (this_ *BaseDialect) DatabaseDDL(param *GenerateParam, database *DatabaseModel) (sqlList []string, err error) {
	var sql string
	sql = `CREATE DATABASE `
	sql += param.packingCharacterDatabase(database.Name)
	return
}

func (this_ *BaseDialect) DatabaseDeleteDDL(param *GenerateParam, database string) (sqlList []string, err error) {
	var sql string
	sql = `DROP DATABASE `
	sql += param.packingCharacterDatabase(database)
	return
}

func (this_ *BaseDialect) TableDDL(param *GenerateParam, database string, table *TableModel) (sqlList []string, err error) {
	sqlList = []string{}

	createTableSql := `CREATE TABLE `

	if param.AppendDatabase && database != "" {
		createTableSql += param.packingCharacterDatabase(database) + "."
	}
	createTableSql += param.packingCharacterTable(table.Name)

	createTableSql += `(`
	createTableSql += "\n"
	primaryKeys := ""
	if len(table.ColumnList) > 0 {
		for _, column := range table.ColumnList {
			var columnSql = param.packingCharacterColumn(column.Name)

			var columnType string
			columnType, err = this_.GetColumnType(column.Type, column.Length, column.Decimal)
			if err != nil {
				return
			}
			columnSql += " " + columnType

			if column.NotNull {
				columnSql += ` NOT NULL`
			}
			if column.Default != nil {
				columnSql += ` DEFAULT ` + formatStringValue("'", GetStringValue(column.Default))
			}

			if column.PrimaryKey {
				primaryKeys += "" + column.Name + ","
			}
			createTableSql += "\t" + columnSql + ",\n"
		}
	}
	if primaryKeys != "" {
		primaryKeys = strings.TrimSuffix(primaryKeys, ",")
		createTableSql += "\tPRIMARY KEY (" + param.packingCharacterColumns(primaryKeys) + ")"
	}

	createTableSql = strings.TrimSuffix(createTableSql, ",\n")
	createTableSql += "\n"

	createTableSql += `)`

	sqlList = append(sqlList, createTableSql)

	// 添加注释
	if table.Comment != "" {
		var sqlList_ []string
		sqlList_, err = this_.TableCommentDDL(param, database, table.Name, table.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}
	if len(table.ColumnList) > 0 {
		for _, one := range table.ColumnList {
			if one.Comment == "" {
				continue
			}
			var sqlList_ []string
			sqlList_, err = this_.TableColumnCommentDDL(param, database, table.Name, one.Name, one.Comment)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
	}

	if len(table.IndexList) > 0 {
		for _, one := range table.IndexList {
			if one.Name == "" || len(one.Columns) == 0 {
				continue
			}
			var sqlList_ []string
			sqlList_, err = this_.TableIndexAddDDL(param, database, table.Name, one)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
	}
	return
}
func (this_ *BaseDialect) TableCommentDDL(param *GenerateParam, database string, table string, comment string) (sqlList []string, err error) {
	sql := "COMMENT ON TABLE  "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)
	sql += " IS " + formatStringValue("'", comment)
	sqlList = append(sqlList, sql)
	return
}
func (this_ *BaseDialect) TableRenameDDL(param *GenerateParam, database string, table string, rename string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += param.packingCharacterTable(table)

	sql += ` RENAME `
	sql += param.packingCharacterColumn(table)
	sql += ` TO `
	sql += param.packingCharacterColumn(rename)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *BaseDialect) TableDeleteDDL(param *GenerateParam, database string, table string) (sqlList []string, err error) {

	var sql string
	sql = `DROP TABLE `
	if param.AppendDatabase {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += param.packingCharacterTable(table)
	return
}

func (this_ *BaseDialect) TableColumnDeleteDDL(param *GenerateParam, database string, table string, column string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += param.packingCharacterTable(table)

	sql += ` DROP COLUMN `
	sql += param.packingCharacterColumn(column)

	sqlList = append(sqlList, sql)

	return
}
func (this_ *BaseDialect) TableColumnRenameDDL(param *GenerateParam, database string, table string, column string, rename string) (sqlList []string, err error) {
	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += param.packingCharacterTable(table)

	sql += ` RENAME COLUMN `
	sql += param.packingCharacterColumn(column)
	sql += ` TO `
	sql += param.packingCharacterColumn(rename)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *BaseDialect) TableColumnUpdateDDL(param *GenerateParam, database string, table string, column *TableColumnModel) (sqlList []string, err error) {
	var columnType string
	columnType, err = this_.GetColumnType(column.Type, column.Length, column.Decimal)
	if err != nil {
		return
	}

	var sqlList_ []string

	if column.OldName != column.Name {
		sqlList_, err = this_.TableColumnRenameDDL(param, database, table, column.OldName, column.Name)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}

	if column.Type != column.OldType ||
		column.Length != column.OldLength ||
		column.Decimal != column.OldDecimal ||
		column.NotNull != column.OldNotNull ||
		column.Default != column.OldDefault ||
		column.BeforeColumn != "" {
		var sql string
		sql = `ALTER TABLE `

		if param.AppendDatabase {
			sql += param.packingCharacterDatabase(database) + "."
		}
		sql += param.packingCharacterTable(table)

		sql += ` MODIFY (`
		sql += param.packingCharacterColumn(column.Name)
		sql += ` ` + columnType + ``
		if column.NotNull {
			sql += ` NOT NULL`
		}
		if column.Default != nil {
			sql += ` DEFAULT ` + formatStringValue("'", GetStringValue(column.Default))
		}
		sql += `)`

		sqlList = append(sqlList, sql)
	}
	if column.Comment != column.OldComment {
		sqlList_, err = this_.TableColumnCommentDDL(param, database, table, column.Name, column.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}
	return
}
func (this_ *BaseDialect) TableColumnAddDDL(param *GenerateParam, database string, table string, column *TableColumnModel) (sqlList []string, err error) {
	var columnType string
	columnType, err = this_.GetColumnType(column.Type, column.Length, column.Decimal)
	if err != nil {
		return
	}

	var sql string
	sql = `ALTER TABLE `

	if param.AppendDatabase {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += param.packingCharacterTable(table)

	sql += ` ADD (`
	sql += param.packingCharacterColumn(column.Name)
	sql += ` ` + columnType + ``
	if column.NotNull {
		sql += ` NOT NULL`
	}
	if column.Default != nil {
		sql += ` DEFAULT ` + formatStringValue("'", GetStringValue(column.Default))
	}
	sql += `)`

	sqlList = append(sqlList, sql)

	if column.Comment != "" {
		var sqlList_ []string
		sqlList_, err = this_.TableColumnCommentDDL(param, database, table, column.Name, column.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}

	return
}

func (this_ *BaseDialect) TableColumnCommentDDL(param *GenerateParam, database string, table string, column string, comment string) (sqlList []string, err error) {
	sql := "COMMENT ON COLUMN "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)
	sql += "." + param.packingCharacterColumn(column)
	sql += " IS " + formatStringValue("'", comment)
	sqlList = append(sqlList, sql)
	return
}

func (this_ *BaseDialect) TableIndexAddDDL(param *GenerateParam, database string, table string, index *TableIndexModel) (sqlList []string, err error) {

	sql := "CREATE "
	switch strings.ToUpper(index.Type) {
	case "UNIQUE":
		sql += "UNIQUE INDEX"
	default:
		sql += "INDEX"
	}

	sql += " " + param.packingCharacterColumn(index.Name) + ""

	sql += " ON "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)

	sql += "(" + param.packingCharacterColumns(strings.Join(index.Columns, ",")) + ")"

	sqlList = append(sqlList, sql)
	return
}

func (this_ *BaseDialect) TableIndexUpdateDDL(param *GenerateParam, database string, table string, index *TableIndexModel) (sqlList []string, err error) {
	var sqlList_ []string
	var sql = " DROP INDEX " + param.packingCharacterColumn(index.OldName) + ""
	sqlList = append(sqlList, sql)

	sqlList_, err = this_.TableIndexAddDDL(param, database, table, index)
	if err != nil {
		return
	}
	sqlList = append(sqlList, sqlList_...)

	return
}

func (this_ *BaseDialect) TableIndexRenameDDL(param *GenerateParam, database string, table string, index string, rename string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)

	sql += ` RENAME INDEX `
	sql += "" + param.packingCharacterColumn(index)
	sql += ` TO `
	sql += "" + param.packingCharacterColumn(rename)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *BaseDialect) TableIndexDeleteDDL(param *GenerateParam, database string, table string, index string) (sqlList []string, err error) {
	sql := "DROP INDEX "
	sql += "" + param.packingCharacterColumn(index)

	sqlList = append(sqlList, sql)
	return
}

func (this_ *BaseDialect) TablePrimaryKeyDeleteDDL(param *GenerateParam, database string, table string, primaryKeys []string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)

	sql += ` DROP PRIMARY KEY `

	sqlList = append(sqlList, sql)
	return
}

func (this_ *BaseDialect) TablePrimaryKeyAddDDL(param *GenerateParam, database string, table string, primaryKeys []string) (sqlList []string, err error) {
	sql := "ALTER TABLE "
	if param.AppendDatabase && database != "" {
		sql += param.packingCharacterDatabase(database) + "."
	}
	sql += "" + param.packingCharacterTable(table)

	sql += ` ADD PRIMARY KEY `

	sql += "(" + param.packingCharacterColumns(strings.Join(primaryKeys, ",")) + ")"

	sqlList = append(sqlList, sql)
	return
}
