package db

import (
	"errors"
	"strconv"
	"strings"
	"teamide/pkg/util"
)

var (
	columnTypeNames []string
)

type ColumnTypeInfo struct {
	Name       string `json:"name,omitempty" column:"name"`
	TypeFormat string `json:"typeFormat,omitempty" column:"typeFormat"`
	HasLength  bool   `json:"hasLength,omitempty" column:"hasLength"`
	HasDecimal bool   `json:"hasDecimal,omitempty" column:"hasDecimal"`
	IsNumber   bool   `json:"isNumber,omitempty" column:"isNumber"`
	IsString   bool   `json:"isString,omitempty" column:"isString"`
	IsDateTime bool   `json:"isDateTime,omitempty" column:"isDateTime"`
	IsBytes    bool   `json:"isBytes,omitempty" column:"isBytes"`
	MinLength  int    `json:"minLength,omitempty" column:"minLength"`
	MaxLength  int    `json:"maxLength,omitempty" column:"maxLength"`
}

func (this_ *ColumnTypeInfo) FormatColumnType(length int, decimal int) (columnType string) {
	if this_.TypeFormat == "VARCHAR2($l)" {
		if length <= 0 {
			length = 4000
		}
	}
	columnType = this_.TypeFormat
	lStr := ""
	dStr := ""
	if length > 0 {
		lStr = strconv.Itoa(length)
		if decimal > 0 {
			dStr = strconv.Itoa(decimal)
		}
	}
	columnType = strings.ReplaceAll(columnType, "$l", lStr)
	columnType = strings.ReplaceAll(columnType, "$d", dStr)
	columnType = strings.ReplaceAll(columnType, " ", "")
	columnType = strings.ReplaceAll(columnType, ",)", ")")
	columnType = strings.TrimSuffix(columnType, "(,)")
	columnType = strings.TrimSuffix(columnType, "()")
	return
}

func CheckColumnType() (err error) {
	//for _, databaseType := range DatabaseTypes {
	//	err = checkDatabaseColumnType(databaseType)
	//	if err != nil {
	//		return
	//	}
	//}
	return

}
func checkDatabaseColumnType(databaseType *DatabaseType) (err error) {
	for _, columnTypeName := range columnTypeNames {
		var find = databaseType.ColumnTypeInfoMap[columnTypeName]
		if find == nil {
			err = errors.New("驱动[" + databaseType.DriverName + "]字段类型[" + columnTypeName + "]未做映射")
			return
		}
	}
	return
}

func AppendColumnTypeInfo(databaseType *DatabaseType, columnTypeInfo *ColumnTypeInfo) {
	if util.ContainsString(columnTypeNames, columnTypeInfo.Name) < 0 {
		columnTypeNames = append(columnTypeNames, columnTypeInfo.Name)
	}
	databaseType.ColumnTypeInfoMap[columnTypeInfo.Name] = columnTypeInfo
}
