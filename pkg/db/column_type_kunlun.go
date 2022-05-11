package db

var (

	/** 数值类型 **/

	KunLunBIT       = addKunLunColumnType(&ColumnTypeInfo{Name: "BIT", TypeFormat: "NUMBER($l, $d)", HasLength: false, IsNumber: true})
	KunLunTINYINT   = addKunLunColumnType(&ColumnTypeInfo{Name: "TINYINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KunLunSMALLINT  = addKunLunColumnType(&ColumnTypeInfo{Name: "SMALLINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KunLunMEDIUMINT = addKunLunColumnType(&ColumnTypeInfo{Name: "MEDIUMINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KunLunINT       = addKunLunColumnType(&ColumnTypeInfo{Name: "INT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KunLunINTEGER   = addKunLunColumnType(&ColumnTypeInfo{Name: "INTEGER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KunLunBIGINT    = addKunLunColumnType(&ColumnTypeInfo{Name: "BIGINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 小数 **/

	KunLunFLOAT   = addKunLunColumnType(&ColumnTypeInfo{Name: "FLOAT", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	KunLunDOUBLE  = addKunLunColumnType(&ColumnTypeInfo{Name: "DOUBLE", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	KunLunDEC     = addKunLunColumnType(&ColumnTypeInfo{Name: "DEC", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KunLunDECIMAL = addKunLunColumnType(&ColumnTypeInfo{Name: "DECIMAL", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KunLunNUMBER  = addKunLunColumnType(&ColumnTypeInfo{Name: "NUMBER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 日期/时间类型 **/

	KunLunYEAR      = addKunLunColumnType(&ColumnTypeInfo{Name: "YEAR", TypeFormat: "YEAR", IsDateTime: true})
	KunLunTIME      = addKunLunColumnType(&ColumnTypeInfo{Name: "TIME", TypeFormat: "TIME", IsDateTime: true})
	KunLunDATE      = addKunLunColumnType(&ColumnTypeInfo{Name: "DATE", TypeFormat: "DATE", IsDateTime: true})
	KunLunDATETIME  = addKunLunColumnType(&ColumnTypeInfo{Name: "DATETIME", TypeFormat: "DATETIME", IsDateTime: true})
	KunLunTIMESTAMP = addKunLunColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", TypeFormat: "TIMESTAMP", IsDateTime: true})

	/** 字符串类型 **/

	KunLunCHAR       = addKunLunColumnType(&ColumnTypeInfo{Name: "CHAR", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	KunLunVARCHAR    = addKunLunColumnType(&ColumnTypeInfo{Name: "VARCHAR", TypeFormat: "VARCHAR($l)", HasLength: true, IsString: true})
	KunLunTINYTEXT   = addKunLunColumnType(&ColumnTypeInfo{Name: "TINYTEXT", TypeFormat: "TINYTEXT($l)", HasLength: true, IsString: true})
	KunLunTEXT       = addKunLunColumnType(&ColumnTypeInfo{Name: "TEXT", TypeFormat: "TEXT($l)", HasLength: true, IsString: true})
	KunLunMEDIUMTEXT = addKunLunColumnType(&ColumnTypeInfo{Name: "MEDIUMTEXT", TypeFormat: "MEDIUMTEXT($l)", HasLength: true, IsString: true})
	KunLunLONGTEXT   = addKunLunColumnType(&ColumnTypeInfo{Name: "LONGTEXT", TypeFormat: "LONGTEXT($l)", HasLength: true, IsString: true})
	KunLunENUM       = addKunLunColumnType(&ColumnTypeInfo{Name: "ENUM", TypeFormat: "ENUM($l)", HasLength: true, IsString: true})
	KunLunTINYBLOB   = addKunLunColumnType(&ColumnTypeInfo{Name: "TINYBLOB", TypeFormat: "TINYBLOB($l)", HasLength: true, IsString: true})
	KunLunBLOB       = addKunLunColumnType(&ColumnTypeInfo{Name: "BLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})
	KunLunMEDIUMBLOB = addKunLunColumnType(&ColumnTypeInfo{Name: "MEDIUMBLOB", TypeFormat: "MEDIUMBLOB($l)", HasLength: true, IsString: true})
	KunLunLONGBLOB   = addKunLunColumnType(&ColumnTypeInfo{Name: "LONGBLOB", TypeFormat: "LONGBLOB($l)", HasLength: true, IsString: true})

	KunLunSET = addKunLunColumnType(&ColumnTypeInfo{Name: "SET", TypeFormat: "SET($l)", HasLength: true, IsString: true})
)

func addKunLunColumnType(columnTypeInfo *ColumnTypeInfo) *ColumnTypeInfo {
	AppendColumnTypeInfo(DatabaseTypeKunLun, columnTypeInfo)
	return columnTypeInfo
}
