package db

var (

	/** 数值类型 **/

	OracleBIT       = addOracleColumnType(&ColumnTypeInfo{Name: "BIT", TypeFormat: "NUMBER($l, $d)", HasLength: false, IsNumber: true})
	OracleTINYINT   = addOracleColumnType(&ColumnTypeInfo{Name: "TINYINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	OracleSMALLINT  = addOracleColumnType(&ColumnTypeInfo{Name: "SMALLINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	OracleMEDIUMINT = addOracleColumnType(&ColumnTypeInfo{Name: "MEDIUMINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	OracleINT       = addOracleColumnType(&ColumnTypeInfo{Name: "INT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	OracleINTEGER   = addOracleColumnType(&ColumnTypeInfo{Name: "INTEGER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	OracleBIGINT    = addOracleColumnType(&ColumnTypeInfo{Name: "BIGINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 小数 **/

	OracleFLOAT   = addOracleColumnType(&ColumnTypeInfo{Name: "FLOAT", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	OracleDOUBLE  = addOracleColumnType(&ColumnTypeInfo{Name: "DOUBLE", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	OracleDEC     = addOracleColumnType(&ColumnTypeInfo{Name: "DEC", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	OracleDECIMAL = addOracleColumnType(&ColumnTypeInfo{Name: "DECIMAL", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	OracleNUMBER  = addOracleColumnType(&ColumnTypeInfo{Name: "NUMBER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 日期/时间类型 **/

	OracleYEAR      = addOracleColumnType(&ColumnTypeInfo{Name: "YEAR", TypeFormat: "DATE", IsDateTime: true})
	OracleTIME      = addOracleColumnType(&ColumnTypeInfo{Name: "TIME", TypeFormat: "DATE", IsDateTime: true})
	OracleDATE      = addOracleColumnType(&ColumnTypeInfo{Name: "DATE", TypeFormat: "DATE", IsDateTime: true})
	OracleDATETIME  = addOracleColumnType(&ColumnTypeInfo{Name: "DATETIME", TypeFormat: "DATE", IsDateTime: true})
	OracleTIMESTAMP = addOracleColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", TypeFormat: "TIMESTAMP", IsDateTime: true})

	/** 字符串类型 **/

	OracleCHAR       = addOracleColumnType(&ColumnTypeInfo{Name: "CHAR", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	OracleVARCHAR    = addOracleColumnType(&ColumnTypeInfo{Name: "VARCHAR", TypeFormat: "VARCHAR2($l)", HasLength: true, IsString: true})
	OracleTINYTEXT   = addOracleColumnType(&ColumnTypeInfo{Name: "TINYTEXT", TypeFormat: "VARCHAR2($l)", HasLength: true, IsString: true})
	OracleTEXT       = addOracleColumnType(&ColumnTypeInfo{Name: "TEXT", TypeFormat: "VARCHAR2($l)", HasLength: true, IsString: true})
	OracleMEDIUMTEXT = addOracleColumnType(&ColumnTypeInfo{Name: "MEDIUMTEXT", TypeFormat: "CLOB($l)", HasLength: true, IsString: true})
	OracleLONGTEXT   = addOracleColumnType(&ColumnTypeInfo{Name: "LONGTEXT", TypeFormat: "CLOB($l)", HasLength: true, IsString: true})
	OracleENUM       = addOracleColumnType(&ColumnTypeInfo{Name: "ENUM", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	OracleTINYBLOB   = addOracleColumnType(&ColumnTypeInfo{Name: "TINYBLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})
	OracleBLOB       = addOracleColumnType(&ColumnTypeInfo{Name: "BLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})
	OracleMEDIUMBLOB = addOracleColumnType(&ColumnTypeInfo{Name: "MEDIUMBLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})
	OracleLONGBLOB   = addOracleColumnType(&ColumnTypeInfo{Name: "LONGBLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})

	OracleSET = addOracleColumnType(&ColumnTypeInfo{Name: "SET", TypeFormat: "SET($l)", HasLength: true, IsString: true})
)

func addOracleColumnType(columnTypeInfo *ColumnTypeInfo) *ColumnTypeInfo {
	AppendColumnTypeInfo(DatabaseTypeOracle, columnTypeInfo)
	return columnTypeInfo
}
