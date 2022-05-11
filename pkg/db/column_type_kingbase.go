package db

var (

	/** 数值类型 **/

	KingBaseBIT       = addKingBaseColumnType(&ColumnTypeInfo{Name: "BIT", TypeFormat: "NUMBER($l, $d)", HasLength: false, IsNumber: true})
	KingBaseTINYINT   = addKingBaseColumnType(&ColumnTypeInfo{Name: "TINYINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KingBaseSMALLINT  = addKingBaseColumnType(&ColumnTypeInfo{Name: "SMALLINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KingBaseMEDIUMINT = addKingBaseColumnType(&ColumnTypeInfo{Name: "MEDIUMINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KingBaseINT       = addKingBaseColumnType(&ColumnTypeInfo{Name: "INT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KingBaseINTEGER   = addKingBaseColumnType(&ColumnTypeInfo{Name: "INTEGER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KingBaseBIGINT    = addKingBaseColumnType(&ColumnTypeInfo{Name: "BIGINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 小数 **/

	KingBaseFLOAT   = addKingBaseColumnType(&ColumnTypeInfo{Name: "FLOAT", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	KingBaseDOUBLE  = addKingBaseColumnType(&ColumnTypeInfo{Name: "DOUBLE", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	KingBaseDEC     = addKingBaseColumnType(&ColumnTypeInfo{Name: "DEC", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KingBaseDECIMAL = addKingBaseColumnType(&ColumnTypeInfo{Name: "DECIMAL", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	KingBaseNUMBER  = addKingBaseColumnType(&ColumnTypeInfo{Name: "NUMBER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 日期/时间类型 **/

	KingBaseYEAR      = addKingBaseColumnType(&ColumnTypeInfo{Name: "YEAR", TypeFormat: "YEAR", IsDateTime: true})
	KingBaseTIME      = addKingBaseColumnType(&ColumnTypeInfo{Name: "TIME", TypeFormat: "TIME", IsDateTime: true})
	KingBaseDATE      = addKingBaseColumnType(&ColumnTypeInfo{Name: "DATE", TypeFormat: "DATE", IsDateTime: true})
	KingBaseDATETIME  = addKingBaseColumnType(&ColumnTypeInfo{Name: "DATETIME", TypeFormat: "DATETIME", IsDateTime: true})
	KingBaseTIMESTAMP = addKingBaseColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", TypeFormat: "TIMESTAMP", IsDateTime: true})

	/** 字符串类型 **/

	KingBaseCHAR       = addKingBaseColumnType(&ColumnTypeInfo{Name: "CHAR", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	KingBaseVARCHAR    = addKingBaseColumnType(&ColumnTypeInfo{Name: "VARCHAR", TypeFormat: "VARCHAR($l)", HasLength: true, IsString: true})
	KingBaseTINYTEXT   = addKingBaseColumnType(&ColumnTypeInfo{Name: "TINYTEXT", TypeFormat: "TINYTEXT($l)", HasLength: true, IsString: true})
	KingBaseTEXT       = addKingBaseColumnType(&ColumnTypeInfo{Name: "TEXT", TypeFormat: "TEXT($l)", HasLength: true, IsString: true})
	KingBaseMEDIUMTEXT = addKingBaseColumnType(&ColumnTypeInfo{Name: "MEDIUMTEXT", TypeFormat: "MEDIUMTEXT($l)", HasLength: true, IsString: true})
	KingBaseLONGTEXT   = addKingBaseColumnType(&ColumnTypeInfo{Name: "LONGTEXT", TypeFormat: "LONGTEXT($l)", HasLength: true, IsString: true})
	KingBaseENUM       = addKingBaseColumnType(&ColumnTypeInfo{Name: "ENUM", TypeFormat: "ENUM($l)", HasLength: true, IsString: true})
	KingBaseTINYBLOB   = addKingBaseColumnType(&ColumnTypeInfo{Name: "TINYBLOB", TypeFormat: "TINYBLOB($l)", HasLength: true, IsString: true})
	KingBaseBLOB       = addKingBaseColumnType(&ColumnTypeInfo{Name: "BLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})
	KingBaseMEDIUMBLOB = addKingBaseColumnType(&ColumnTypeInfo{Name: "MEDIUMBLOB", TypeFormat: "MEDIUMBLOB($l)", HasLength: true, IsString: true})
	KingBaseLONGBLOB   = addKingBaseColumnType(&ColumnTypeInfo{Name: "LONGBLOB", TypeFormat: "LONGBLOB($l)", HasLength: true, IsString: true})

	KingBaseSET = addKingBaseColumnType(&ColumnTypeInfo{Name: "SET", TypeFormat: "SET($l)", HasLength: true, IsString: true})
)

func addKingBaseColumnType(columnTypeInfo *ColumnTypeInfo) *ColumnTypeInfo {
	AppendColumnTypeInfo(DatabaseTypeKingBase, columnTypeInfo)
	return columnTypeInfo
}
