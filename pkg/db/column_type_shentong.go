package db

var (

	/** 数值类型 **/

	ShenTongBIT       = addShenTongColumnType(&ColumnTypeInfo{Name: "BIT", TypeFormat: "NUMBER($l, $d)", HasLength: false, IsNumber: true})
	ShenTongTINYINT   = addShenTongColumnType(&ColumnTypeInfo{Name: "TINYINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	ShenTongSMALLINT  = addShenTongColumnType(&ColumnTypeInfo{Name: "SMALLINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	ShenTongMEDIUMINT = addShenTongColumnType(&ColumnTypeInfo{Name: "MEDIUMINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	ShenTongINT       = addShenTongColumnType(&ColumnTypeInfo{Name: "INT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	ShenTongINTEGER   = addShenTongColumnType(&ColumnTypeInfo{Name: "INTEGER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	ShenTongBIGINT    = addShenTongColumnType(&ColumnTypeInfo{Name: "BIGINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 小数 **/

	ShenTongFLOAT   = addShenTongColumnType(&ColumnTypeInfo{Name: "FLOAT", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	ShenTongDOUBLE  = addShenTongColumnType(&ColumnTypeInfo{Name: "DOUBLE", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	ShenTongDEC     = addShenTongColumnType(&ColumnTypeInfo{Name: "DEC", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	ShenTongDECIMAL = addShenTongColumnType(&ColumnTypeInfo{Name: "DECIMAL", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	ShenTongNUMBER  = addShenTongColumnType(&ColumnTypeInfo{Name: "NUMBER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 日期/时间类型 **/

	ShenTongYEAR      = addShenTongColumnType(&ColumnTypeInfo{Name: "YEAR", TypeFormat: "YEAR", IsDateTime: true})
	ShenTongTIME      = addShenTongColumnType(&ColumnTypeInfo{Name: "TIME", TypeFormat: "TIME", IsDateTime: true})
	ShenTongDATE      = addShenTongColumnType(&ColumnTypeInfo{Name: "DATE", TypeFormat: "DATE", IsDateTime: true})
	ShenTongDATETIME  = addShenTongColumnType(&ColumnTypeInfo{Name: "DATETIME", TypeFormat: "DATETIME", IsDateTime: true})
	ShenTongTIMESTAMP = addShenTongColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", TypeFormat: "TIMESTAMP", IsDateTime: true})

	/** 字符串类型 **/

	ShenTongCHAR       = addShenTongColumnType(&ColumnTypeInfo{Name: "CHAR", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	ShenTongVARCHAR    = addShenTongColumnType(&ColumnTypeInfo{Name: "VARCHAR", TypeFormat: "VARCHAR($l)", HasLength: true, IsString: true})
	ShenTongTINYTEXT   = addShenTongColumnType(&ColumnTypeInfo{Name: "TINYTEXT", TypeFormat: "TINYTEXT($l)", HasLength: true, IsString: true})
	ShenTongTEXT       = addShenTongColumnType(&ColumnTypeInfo{Name: "TEXT", TypeFormat: "TEXT($l)", HasLength: true, IsString: true})
	ShenTongMEDIUMTEXT = addShenTongColumnType(&ColumnTypeInfo{Name: "MEDIUMTEXT", TypeFormat: "MEDIUMTEXT($l)", HasLength: true, IsString: true})
	ShenTongLONGTEXT   = addShenTongColumnType(&ColumnTypeInfo{Name: "LONGTEXT", TypeFormat: "LONGTEXT($l)", HasLength: true, IsString: true})
	ShenTongENUM       = addShenTongColumnType(&ColumnTypeInfo{Name: "ENUM", TypeFormat: "ENUM($l)", HasLength: true, IsString: true})
	ShenTongTINYBLOB   = addShenTongColumnType(&ColumnTypeInfo{Name: "TINYBLOB", TypeFormat: "TINYBLOB($l)", HasLength: true, IsString: true})
	ShenTongBLOB       = addShenTongColumnType(&ColumnTypeInfo{Name: "BLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})
	ShenTongMEDIUMBLOB = addShenTongColumnType(&ColumnTypeInfo{Name: "MEDIUMBLOB", TypeFormat: "MEDIUMBLOB($l)", HasLength: true, IsString: true})
	ShenTongLONGBLOB   = addShenTongColumnType(&ColumnTypeInfo{Name: "LONGBLOB", TypeFormat: "LONGBLOB($l)", HasLength: true, IsString: true})

	ShenTongSET = addShenTongColumnType(&ColumnTypeInfo{Name: "SET", TypeFormat: "SET($l)", HasLength: true, IsString: true})
)

func addShenTongColumnType(columnTypeInfo *ColumnTypeInfo) *ColumnTypeInfo {
	AppendColumnTypeInfo(DatabaseTypeShenTong, columnTypeInfo)
	return columnTypeInfo
}
