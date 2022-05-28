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
	DateFunction             string `json:"dateFunction" column:"dateFunction"`
	OpenTransaction          bool   `json:"openTransaction"`
	ErrorContinue            bool   `json:"errorContinue"`
}

func ToDatabaseDDL(param *GenerateParam, database *DatabaseModel) (sqlList []string, err error) {
	databaseType := GetDatabaseType(param.DatabaseType)
	if databaseType == nil {
		err = errors.New("数据库类型[" + param.DatabaseType + "]暂不支持")
		return
	}
	dialect, err := GetDialect(databaseType)
	if err != nil {
		return
	}
	sqlList, err = dialect.DatabaseDDL(param, database)

	return
}

func ToTableDDL(param *GenerateParam, database string, table *TableModel) (sqlList []string, err error) {
	databaseType := GetDatabaseType(param.DatabaseType)
	if databaseType == nil {
		err = errors.New("数据库类型[" + param.DatabaseType + "]暂不支持")
		return
	}
	dialect, err := GetDialect(databaseType)
	if err != nil {
		return
	}
	sqlList, err = dialect.TableDDL(param, database, table)

	return
}

func ToTableUpdateDDL(param *GenerateParam, database string, table *TableModel) (sqlList []string, err error) {
	databaseType := GetDatabaseType(param.DatabaseType)
	if databaseType == nil {
		err = errors.New("数据库类型[" + param.DatabaseType + "]暂不支持")
		return
	}
	dialect, err := GetDialect(databaseType)
	if err != nil {
		return
	}

	if table == nil {
		err = errors.New("请传入表信息")
		return
	}
	var sqlList_ []string
	if table.Comment != table.OldComment {
		sqlList_, err = dialect.TableCommentDDL(param, database, table.Name, table.Comment)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}
	var primaryKeys []string
	var oldPrimaryKeys []string
	if len(table.ColumnList) > 0 {
		for _, column := range table.ColumnList {
			if column.OldPrimaryKey {
				oldPrimaryKeys = append(oldPrimaryKeys, column.OldName)
			}
			if column.Deleted {
				sqlList_, err = dialect.TableColumnDeleteDDL(param, database, table.Name, column.Name)
				if err != nil {
					return
				}
				sqlList = append(sqlList, sqlList_...)
				continue
			}
			if column.PrimaryKey {
				primaryKeys = append(primaryKeys, column.Name)
			}
			if column.OldName == "" {
				sqlList_, err = dialect.TableColumnAddDDL(param, database, table.Name, column)
				if err != nil {
					return
				}
				sqlList = append(sqlList, sqlList_...)
				continue
			}
			if column.Name != column.OldName ||
				column.Type != column.OldType ||
				column.Length != column.OldLength ||
				column.Decimal != column.OldDecimal ||
				column.NotNull != column.OldNotNull ||
				column.Default != column.OldDefault ||
				column.Comment != column.OldComment ||
				column.BeforeColumn != "" {
				sqlList_, err = dialect.TableColumnUpdateDDL(param, database, table.Name, column)
				if err != nil {
					return
				}
				sqlList = append(sqlList, sqlList_...)
			}
		}
	}
	if strings.Join(primaryKeys, ",") != strings.Join(oldPrimaryKeys, ",") {
		if len(oldPrimaryKeys) > 0 {
			sqlList_, err = dialect.TablePrimaryKeyDeleteDDL(param, database, table.Name, oldPrimaryKeys)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
		if len(primaryKeys) > 0 {
			sqlList_, err = dialect.TablePrimaryKeyAddDDL(param, database, table.Name, primaryKeys)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
	}
	if len(table.IndexList) > 0 {
		for _, index := range table.IndexList {
			if index.Deleted {
				sqlList_, err = dialect.TableIndexDeleteDDL(param, database, table.Name, index.Name)
				if err != nil {
					return
				}
				sqlList = append(sqlList, sqlList_...)
				continue
			}
			if index.OldName == "" {
				sqlList_, err = dialect.TableIndexAddDDL(param, database, table.Name, index)
				if err != nil {
					return
				}
				sqlList = append(sqlList, sqlList_...)
				continue
			}
			if index.Name != index.OldName ||
				index.Type != index.OldType ||
				index.Comment != index.OldComment ||
				strings.Join(index.Columns, ",") != strings.Join(index.OldColumns, ",") {
				sqlList_, err = dialect.TableIndexUpdateDDL(param, database, table.Name, index)
				if err != nil {
					return
				}
				sqlList = append(sqlList, sqlList_...)
			}
		}
	}

	return
}

func ToDatabaseDeleteDDL(param *GenerateParam, database string) (sqlList []string, err error) {
	databaseType := GetDatabaseType(param.DatabaseType)
	if databaseType == nil {
		err = errors.New("数据库类型[" + param.DatabaseType + "]暂不支持")
		return
	}
	dialect, err := GetDialect(databaseType)
	if err != nil {
		return
	}
	sqlList, err = dialect.DatabaseDeleteDDL(param, database)

	return
}

func ToTableDeleteDDL(param *GenerateParam, database string, table string) (sqlList []string, err error) {
	databaseType := GetDatabaseType(param.DatabaseType)
	if databaseType == nil {
		err = errors.New("数据库类型[" + param.DatabaseType + "]暂不支持")
		return
	}
	dialect, err := GetDialect(databaseType)
	if err != nil {
		return
	}
	sqlList, err = dialect.TableDeleteDDL(param, database, table)

	return
}

func (param *GenerateParam) PackingCharacterDatabase(value string) string {
	return param.packingCharacterDatabase(value)
}

func (param *GenerateParam) packingCharacterDatabase(value string) string {
	if param.DatabasePackingCharacter == "" {
		return value
	}
	return param.DatabasePackingCharacter + value + param.DatabasePackingCharacter
}

func (param *GenerateParam) PackingCharacterTable(value string) string {
	return param.packingCharacterTable(value)
}

func (param *GenerateParam) packingCharacterTable(value string) string {
	if param.TablePackingCharacter == "" {
		return value
	}
	return param.TablePackingCharacter + value + param.TablePackingCharacter
}

func (param *GenerateParam) PackingCharacterColumn(value string) string {
	return param.packingCharacterColumn(value)
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

func (param *GenerateParam) PackingCharacterColumns(value string) string {
	return param.packingCharacterColumns(value)
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

func (param *GenerateParam) PackingCharacterColumnStringValue(tableColumn *TableColumnModel, value interface{}) string {
	return param.packingCharacterColumnStringValue(tableColumn, value)
}

func (param *GenerateParam) packingCharacterColumnStringValue(tableColumn *TableColumnModel, value interface{}) string {
	var formatColumnValue = param.formatColumnValue(tableColumn, value)
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
		if param.DateFunction != "" {
			return strings.ReplaceAll(param.DateFunction, "$value", valueString)
		}
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

func (param *GenerateParam) FormatColumnValue(tableColumn *TableColumnModel, value interface{}) interface{} {
	return param.formatColumnValue(tableColumn, value)
}
func (param *GenerateParam) formatColumnValue(tableColumn *TableColumnModel, value interface{}) interface{} {

	var IsDateTime bool
	var IsNumber bool
	var Decimal int
	if tableColumn != nil {

		columnTypeInfo := DatabaseTypeMySql.GetColumnTypeInfo(tableColumn.Type)
		if columnTypeInfo == nil {
			util.Logger.Warn("字段类型[" + tableColumn.Type + "]未引射信息")
			return value
		}
		IsDateTime = columnTypeInfo.IsDateTime
		IsNumber = columnTypeInfo.IsNumber
		Decimal = tableColumn.Decimal
	}

	if value == nil {
		return value
	}
	var stringValue = GetStringValue(value)
	if IsNumber {
		if stringValue == "" {
			return nil
		}
		if Decimal > 0 {
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
	if IsDateTime {
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
