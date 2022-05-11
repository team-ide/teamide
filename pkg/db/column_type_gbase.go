package db

var (

	/** 数值类型 **/

	GBaseBIT       = addGBaseColumnType(&ColumnTypeInfo{Name: "BIT", TypeFormat: "NUMBER($l, $d)", HasLength: false, IsNumber: true})
	GBaseTINYINT   = addGBaseColumnType(&ColumnTypeInfo{Name: "TINYINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	GBaseSMALLINT  = addGBaseColumnType(&ColumnTypeInfo{Name: "SMALLINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	GBaseMEDIUMINT = addGBaseColumnType(&ColumnTypeInfo{Name: "MEDIUMINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	GBaseINT       = addGBaseColumnType(&ColumnTypeInfo{Name: "INT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	GBaseINTEGER   = addGBaseColumnType(&ColumnTypeInfo{Name: "INTEGER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	GBaseBIGINT    = addGBaseColumnType(&ColumnTypeInfo{Name: "BIGINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 小数 **/

	GBaseFLOAT   = addGBaseColumnType(&ColumnTypeInfo{Name: "FLOAT", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	GBaseDOUBLE  = addGBaseColumnType(&ColumnTypeInfo{Name: "DOUBLE", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	GBaseDEC     = addGBaseColumnType(&ColumnTypeInfo{Name: "DEC", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	GBaseDECIMAL = addGBaseColumnType(&ColumnTypeInfo{Name: "DECIMAL", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	GBaseNUMBER  = addGBaseColumnType(&ColumnTypeInfo{Name: "NUMBER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 日期/时间类型 **/

	GBaseYEAR      = addGBaseColumnType(&ColumnTypeInfo{Name: "YEAR", TypeFormat: "YEAR", IsDateTime: true})
	GBaseTIME      = addGBaseColumnType(&ColumnTypeInfo{Name: "TIME", TypeFormat: "TIME", IsDateTime: true})
	GBaseDATE      = addGBaseColumnType(&ColumnTypeInfo{Name: "DATE", TypeFormat: "DATE", IsDateTime: true})
	GBaseDATETIME  = addGBaseColumnType(&ColumnTypeInfo{Name: "DATETIME", TypeFormat: "DATETIME", IsDateTime: true})
	GBaseTIMESTAMP = addGBaseColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", TypeFormat: "TIMESTAMP", IsDateTime: true})

	/** 字符串类型 **/

	GBaseCHAR       = addGBaseColumnType(&ColumnTypeInfo{Name: "CHAR", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	GBaseVARCHAR    = addGBaseColumnType(&ColumnTypeInfo{Name: "VARCHAR", TypeFormat: "VARCHAR($l)", HasLength: true, IsString: true})
	GBaseTINYTEXT   = addGBaseColumnType(&ColumnTypeInfo{Name: "TINYTEXT", TypeFormat: "TINYTEXT($l)", HasLength: true, IsString: true})
	GBaseTEXT       = addGBaseColumnType(&ColumnTypeInfo{Name: "TEXT", TypeFormat: "TEXT($l)", HasLength: true, IsString: true})
	GBaseMEDIUMTEXT = addGBaseColumnType(&ColumnTypeInfo{Name: "MEDIUMTEXT", TypeFormat: "MEDIUMTEXT($l)", HasLength: true, IsString: true})
	GBaseLONGTEXT   = addGBaseColumnType(&ColumnTypeInfo{Name: "LONGTEXT", TypeFormat: "LONGTEXT($l)", HasLength: true, IsString: true})
	GBaseENUM       = addGBaseColumnType(&ColumnTypeInfo{Name: "ENUM", TypeFormat: "ENUM($l)", HasLength: true, IsString: true})
	GBaseTINYBLOB   = addGBaseColumnType(&ColumnTypeInfo{Name: "TINYBLOB", TypeFormat: "TINYBLOB($l)", HasLength: true, IsString: true})
	GBaseBLOB       = addGBaseColumnType(&ColumnTypeInfo{Name: "BLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})
	GBaseMEDIUMBLOB = addGBaseColumnType(&ColumnTypeInfo{Name: "MEDIUMBLOB", TypeFormat: "MEDIUMBLOB($l)", HasLength: true, IsString: true})
	GBaseLONGBLOB   = addGBaseColumnType(&ColumnTypeInfo{Name: "LONGBLOB", TypeFormat: "LONGBLOB($l)", HasLength: true, IsString: true})

	GBaseSET = addGBaseColumnType(&ColumnTypeInfo{Name: "SET", TypeFormat: "SET($l)", HasLength: true, IsString: true})
)

func addGBaseColumnType(columnTypeInfo *ColumnTypeInfo) *ColumnTypeInfo {
	AppendColumnTypeInfo(DatabaseTypeGBase, columnTypeInfo)
	return columnTypeInfo
}
