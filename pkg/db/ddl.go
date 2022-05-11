package db

import (
	"errors"
	"strings"
)

type GenerateParam struct {
	DatabaseType             string `json:"databaseType" column:"databaseType"`
	GenerateDatabase         bool   `json:"generateDatabase" column:"generateDatabase"`
	AppendDatabase           bool   `json:"appendDatabase" column:"appendDatabase"`
	CharacterSet             string `json:"characterSet" column:"characterSet"`
	Collate                  string `json:"collate" column:"collate"`
	DatabasePackingCharacter string `json:"databasePackingCharacter" column:"databasePackingCharacter"`
	TablePackingCharacter    string `json:"tablePackingCharacter" column:"tablePackingCharacter"`
	ColumnPackingCharacter   string `json:"columnPackingCharacter" column:"columnPackingCharacter"`
	StringPackingCharacter   string `json:"stringPackingCharacter" column:"stringPackingCharacter"`
}

func ToDatabaseDDL(param *GenerateParam, database *DatabaseModel) (sqlList []string, err error) {

	switch GetDatabaseType(param.DatabaseType) {
	case DatabaseTypeMySql:
		DatabaseMySqlDialect := &DatabaseMySqlDialect{}
		sqlList, err = DatabaseMySqlDialect.DatabaseDDL(param, database)
		break
	case DatabaseTypeSqlite:
		DatabaseSqliteDialect := &DatabaseSqliteDialect{}
		sqlList, err = DatabaseSqliteDialect.DatabaseDDL(param, database)
		break
	case DatabaseTypeOracle:
		DatabaseOracleDialect := &DatabaseOracleDialect{}
		sqlList, err = DatabaseOracleDialect.DatabaseDDL(param, database)
		break
	case DatabaseTypeShenTong:
		DatabaseShenTongDialect := &DatabaseShenTongDialect{}
		sqlList, err = DatabaseShenTongDialect.DatabaseDDL(param, database)
		break
	case DatabaseTypeDM:
		DatabaseDMDialect := &DatabaseDMDialect{}
		sqlList, err = DatabaseDMDialect.DatabaseDDL(param, database)
		break
	case DatabaseTypeKingBase:
		DatabaseKingBaseDialect := &DatabaseKingBaseDialect{}
		sqlList, err = DatabaseKingBaseDialect.DatabaseDDL(param, database)
		break
	case nil:
		err = errors.New("数据库类型[" + param.DatabaseType + "]暂不支持")
		break

	}

	return
}

func ToTableDDL(param *GenerateParam, database string, table *TableModel) (sqlList []string, err error) {

	switch GetDatabaseType(param.DatabaseType) {
	case DatabaseTypeMySql:
		DatabaseMySqlDialect := &DatabaseMySqlDialect{}
		sqlList, err = DatabaseMySqlDialect.TableDDL(param, database, table)
		break
	case DatabaseTypeSqlite:
		DatabaseSqliteDialect := &DatabaseSqliteDialect{}
		sqlList, err = DatabaseSqliteDialect.TableDDL(param, database, table)
		break
	case DatabaseTypeOracle:
		DatabaseOracleDialect := &DatabaseOracleDialect{}
		sqlList, err = DatabaseOracleDialect.TableDDL(param, database, table)
		break
	case DatabaseTypeShenTong:
		DatabaseShenTongDialect := &DatabaseShenTongDialect{}
		sqlList, err = DatabaseShenTongDialect.TableDDL(param, database, table)
		break
	case DatabaseTypeDM:
		DatabaseDMDialect := &DatabaseDMDialect{}
		sqlList, err = DatabaseDMDialect.TableDDL(param, database, table)
		break
	case DatabaseTypeKingBase:
		DatabaseKingBaseDialect := &DatabaseKingBaseDialect{}
		sqlList, err = DatabaseKingBaseDialect.TableDDL(param, database, table)
		break
	case nil:
		err = errors.New("数据库类型[" + param.DatabaseType + "]暂不支持")
		break
	}
	return
}

type DatabaseDialect struct {
}

func (this_ *DatabaseDialect) packingCharacterDatabase(param *GenerateParam, value string) string {
	if param.DatabasePackingCharacter == "" {
		return value
	}
	return param.DatabasePackingCharacter + value + param.DatabasePackingCharacter
}

func (this_ *DatabaseDialect) packingCharacterTable(param *GenerateParam, value string) string {
	if param.TablePackingCharacter == "" {
		return value
	}
	return param.TablePackingCharacter + value + param.TablePackingCharacter
}

func (this_ *DatabaseDialect) packingCharacterColumn(param *GenerateParam, value string) string {
	if param.ColumnPackingCharacter == "" {
		return value
	}
	value = strings.ReplaceAll(value, `""`, "")
	value = strings.ReplaceAll(value, `'`, "")
	value = strings.ReplaceAll(value, "`", "")
	return param.ColumnPackingCharacter + value + param.ColumnPackingCharacter
}

func (this_ *DatabaseDialect) packingCharacterColumns(param *GenerateParam, columns string) string {
	if param.ColumnPackingCharacter == "" {
		return columns
	}
	res := ""
	columnList := strings.Split(columns, ",")

	for _, column := range columnList {
		res += this_.packingCharacterColumn(param, column) + ","
	}
	res = strings.TrimSuffix(res, ",")
	return res
}

func (this_ *DatabaseDialect) packingCharacterString(param *GenerateParam, value interface{}) interface{} {
	if value == nil || param.StringPackingCharacter == "" {
		return value
	}
	valueString, ok := value.(string)
	if !ok {
		return value
	}
	return formatStringValue(param.StringPackingCharacter, valueString)
}

func formatStringValue(packingCharacter string, valueString string) string {
	if packingCharacter == "" {
		return valueString
	}
	ss := strings.Split(valueString, "")
	out := packingCharacter
	for _, s := range ss {
		switch s {
		case packingCharacter:
			out += "\\" + s
		case "\\":
			out += "\\" + s
		default:
			out += s
		}
	}
	out += packingCharacter
	return out
}
