package db

var (

	/** 数值类型 **/

	SqliteBIT       = addSqliteColumnType(&ColumnTypeInfo{Name: "BIT", TypeFormat: "NUMBER($l, $d)", HasLength: false, IsNumber: true})
	SqliteTINYINT   = addSqliteColumnType(&ColumnTypeInfo{Name: "TINYINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	SqliteSMALLINT  = addSqliteColumnType(&ColumnTypeInfo{Name: "SMALLINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	SqliteMEDIUMINT = addSqliteColumnType(&ColumnTypeInfo{Name: "MEDIUMINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	SqliteINT       = addSqliteColumnType(&ColumnTypeInfo{Name: "INT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	SqliteINTEGER   = addSqliteColumnType(&ColumnTypeInfo{Name: "INTEGER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	SqliteBIGINT    = addSqliteColumnType(&ColumnTypeInfo{Name: "BIGINT", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 小数 **/

	SqliteFLOAT   = addSqliteColumnType(&ColumnTypeInfo{Name: "FLOAT", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	SqliteDOUBLE  = addSqliteColumnType(&ColumnTypeInfo{Name: "DOUBLE", TypeFormat: "NUMBER($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	SqliteDEC     = addSqliteColumnType(&ColumnTypeInfo{Name: "DEC", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	SqliteDECIMAL = addSqliteColumnType(&ColumnTypeInfo{Name: "DECIMAL", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})
	SqliteNUMBER  = addSqliteColumnType(&ColumnTypeInfo{Name: "NUMBER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 日期/时间类型 **/

	SqliteYEAR      = addSqliteColumnType(&ColumnTypeInfo{Name: "YEAR", TypeFormat: "YEAR", IsDateTime: true})
	SqliteTIME      = addSqliteColumnType(&ColumnTypeInfo{Name: "TIME", TypeFormat: "TIME", IsDateTime: true})
	SqliteDATE      = addSqliteColumnType(&ColumnTypeInfo{Name: "DATE", TypeFormat: "DATE", IsDateTime: true})
	SqliteDATETIME  = addSqliteColumnType(&ColumnTypeInfo{Name: "DATETIME", TypeFormat: "DATETIME", IsDateTime: true})
	SqliteTIMESTAMP = addSqliteColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", TypeFormat: "TIMESTAMP", IsDateTime: true})

	/** 字符串类型 **/

	SqliteCHAR       = addSqliteColumnType(&ColumnTypeInfo{Name: "CHAR", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	SqliteVARCHAR    = addSqliteColumnType(&ColumnTypeInfo{Name: "VARCHAR", TypeFormat: "VARCHAR($l)", HasLength: true, IsString: true})
	SqliteTINYTEXT   = addSqliteColumnType(&ColumnTypeInfo{Name: "TINYTEXT", TypeFormat: "TINYTEXT($l)", HasLength: true, IsString: true})
	SqliteTEXT       = addSqliteColumnType(&ColumnTypeInfo{Name: "TEXT", TypeFormat: "TEXT($l)", HasLength: true, IsString: true})
	SqliteMEDIUMTEXT = addSqliteColumnType(&ColumnTypeInfo{Name: "MEDIUMTEXT", TypeFormat: "MEDIUMTEXT($l)", HasLength: true, IsString: true})
	SqliteLONGTEXT   = addSqliteColumnType(&ColumnTypeInfo{Name: "LONGTEXT", TypeFormat: "LONGTEXT($l)", HasLength: true, IsString: true})
	SqliteENUM       = addSqliteColumnType(&ColumnTypeInfo{Name: "ENUM", TypeFormat: "ENUM($l)", HasLength: true, IsString: true})
	SqliteTINYBLOB   = addSqliteColumnType(&ColumnTypeInfo{Name: "TINYBLOB", TypeFormat: "TINYBLOB($l)", HasLength: true, IsString: true})
	SqliteBLOB       = addSqliteColumnType(&ColumnTypeInfo{Name: "BLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})
	SqliteMEDIUMBLOB = addSqliteColumnType(&ColumnTypeInfo{Name: "MEDIUMBLOB", TypeFormat: "MEDIUMBLOB($l)", HasLength: true, IsString: true})
	SqliteLONGBLOB   = addSqliteColumnType(&ColumnTypeInfo{Name: "LONGBLOB", TypeFormat: "LONGBLOB($l)", HasLength: true, IsString: true})

	SqliteSET = addSqliteColumnType(&ColumnTypeInfo{Name: "SET", TypeFormat: "SET($l)", HasLength: true, IsString: true})
)

func addSqliteColumnType(columnTypeInfo *ColumnTypeInfo) *ColumnTypeInfo {
	AppendColumnTypeInfo(DatabaseTypeSqlite, columnTypeInfo)
	return columnTypeInfo
}
