package db

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"teamide/pkg/util"
	"time"
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
	AppendSqlValue           bool   `json:"appendSqlValue" column:"appendSqlValue"`
	OpenTransaction          bool   `json:"openTransaction"`
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

func (param *GenerateParam) packingCharacterDatabase(value string) string {
	if param.DatabasePackingCharacter == "" {
		return value
	}
	return param.DatabasePackingCharacter + value + param.DatabasePackingCharacter
}

func (param *GenerateParam) packingCharacterTable(value string) string {
	if param.TablePackingCharacter == "" {
		return value
	}
	return param.TablePackingCharacter + value + param.TablePackingCharacter
}

func (param *GenerateParam) packingCharacterColumn(value string) string {
	if param.ColumnPackingCharacter == "" {
		return value
	}
	value = strings.ReplaceAll(value, `""`, "")
	value = strings.ReplaceAll(value, `'`, "")
	value = strings.ReplaceAll(value, "`", "")
	return param.ColumnPackingCharacter + value + param.ColumnPackingCharacter
}

func (param *GenerateParam) packingCharacterColumns(columns string) string {
	if param.ColumnPackingCharacter == "" {
		return columns
	}
	res := ""
	columnList := strings.Split(columns, ",")

	for _, column := range columnList {
		res += param.packingCharacterColumn(column) + ","
	}
	res = strings.TrimSuffix(res, ",")
	return res
}

func (param *GenerateParam) packingCharacterColumnStringValue(column *TableColumnModel, value interface{}) string {
	var formatColumnValue = param.formatColumnValue(column, value)
	if formatColumnValue == nil {
		return "NULL"
	}
	var valueString string
	switch v := formatColumnValue.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case time.Time:
		if v.IsZero() {
			return "NULL"
		}
		valueString = v.Format("2006-01-02 15:04:05")
		break
	case string:
		valueString = v
		break
	case []byte:
		valueString = string(v)
	default:
		newValue, _ := json.Marshal(value)
		valueString = string(newValue)
		break
	}
	if param.StringPackingCharacter == "" {
		return valueString
	}
	return formatStringValue(param.StringPackingCharacter, valueString)
}
func (param *GenerateParam) formatColumnValue(column *TableColumnModel, value interface{}) interface{} {
	if value == nil {
		return value
	}
	columnTypeInfo := DatabaseTypeMySql.GetColumnTypeInfo(column.Type)
	if columnTypeInfo == nil {
		util.Logger.Warn("字段类型[" + column.Type + "]未引射信息")
		return value
	}
	var stringValue = GetStringValue(value)
	if columnTypeInfo.IsNumber {
		if stringValue == "" {
			return nil
		}
		if column.Decimal > 0 {
			f64, err := strconv.ParseFloat(stringValue, 64)
			if err != nil {
				util.Logger.Error("值["+stringValue+"]转化float64异常", zap.Error(err))
				return value
			}
			return f64
		} else {
			i64, err := strconv.ParseInt(stringValue, 10, 64)
			if err != nil {
				util.Logger.Error("值["+stringValue+"]转化int64异常", zap.Error(err))
				return value
			}
			return i64
		}
	}
	if columnTypeInfo.IsDateTime {
		if stringValue == "" {
			return nil
		}
		format := "2006-01-02 15:04:05.000"
		valueLen := len(stringValue)
		if valueLen >= len("2006-01-02 15:04:05.000") {
			format = "2006-01-02 15:04:05.000"
		} else if valueLen >= len("2006-01-02 15:04:05") {
			format = "2006-01-02 15:04:05"
		} else if valueLen >= len("2006-01-02 15:04") {
			format = "2006-01-02 15:04"
		} else if valueLen >= len("2006-01-02 15") {
			format = "2006-01-02 15"
		} else if valueLen >= len("2006-01-02") {
			format = "2006-01-02"
		} else if valueLen >= len("15:04:05") {
			format = "15:04:05"
		} else if valueLen >= len("15:04") {
			format = "15:04"
		} else if valueLen >= len("2006") {
			format = "2006"
		}
		timeValue, err := time.ParseInLocation(format, stringValue, time.Local)
		if err != nil {
			util.Logger.Error("值["+stringValue+"]转化time异常", zap.Error(err))
			return value
		}
		return timeValue
	}
	return value
}

func GetStringValue(value interface{}) string {

	var valueString string
	switch v := value.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case time.Time:
		if v.IsZero() {
			return "NULL"
		}
		valueString = v.Format("2006-01-02 15:04:05")
		break
	case string:
		valueString = v
		break
	case []byte:
		valueString = string(v)
	default:
		newValue, _ := json.Marshal(value)
		valueString = string(newValue)
		break
	}
	return valueString
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
