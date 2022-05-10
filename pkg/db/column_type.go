package db

import (
	"errors"
	"github.com/wxnacy/wgo/arrays"
)

var (
	columnTypeNames []string
)

type ColumnTypeInfo struct {
	Name                   string   `json:"name,omitempty" column:"name"`
	HasLength              bool     `json:"hasLength,omitempty" column:"hasLength"`
	HasDecimal             bool     `json:"hasDecimal,omitempty" column:"hasDecimal"`
	DataByte               string   `json:"dataByte,omitempty" column:"dataByte"`
	IsNumber               string   `json:"isNumber,omitempty" column:"isNumber"`
	IsString               string   `json:"isString,omitempty" column:"isString"`
	IsDateTime             string   `json:"isDateTime,omitempty" column:"isDateTime"`
	NumberSymbolRangeMin   int      `json:"numberSymbolRangeMin,omitempty" column:"numberSymbolRangeMin"`
	NumberSymbolRangeMax   int      `json:"numberSymbolRangeMax,omitempty" column:"numberSymbolRangeMax"`
	NumberNoSymbolRangeMin int      `json:"numberNoSymbolRangeMin,omitempty" column:"numberNoSymbolRangeMin"`
	NumberNoSymbolRangeMax int      `json:"numberNoSymbolRangeMax,omitempty" column:"numberNoSymbolRangeMax"`
	MinLength              int      `json:"minLength,omitempty" column:"minLength"`
	MaxLength              int      `json:"maxLength,omitempty" column:"maxLength"`
	MaxByte                int      `json:"maxByte,omitempty" column:"maxByte"`
	MatchNames             []string `json:"matchNames,omitempty" column:"matchNames"`
}

func CheckColumnType() (err error) {
	for _, databaseType := range DatabaseTypes {
		err = checkDatabaseColumnType(databaseType)
		if err != nil {
			return
		}
	}
	return

}
func checkDatabaseColumnType(databaseType *DatabaseType) (err error) {
	for _, columnTypeName := range columnTypeNames {
		var find *ColumnTypeInfo
		for _, one := range databaseType.ColumnTypeInfos {
			if columnTypeName == one.Name {
				find = one
				break
			}
		}
		if find == nil {
			err = errors.New("驱动[" + databaseType.DriverName + "]字段类型[" + columnTypeName + "]未做映射")
			return
		}
	}
	return
}

func AppendColumnTypeInfo(databaseType *DatabaseType, columnTypeInfo *ColumnTypeInfo) {
	if arrays.ContainsString(columnTypeNames, columnTypeInfo.Name) < 0 {
		columnTypeNames = append(columnTypeNames, columnTypeInfo.Name)
	}
	databaseType.ColumnTypeInfos = append(databaseType.ColumnTypeInfos, columnTypeInfo)
}
