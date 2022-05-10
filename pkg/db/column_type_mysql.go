package db

var (

	/** 数值类型 **/
	/**
	id INT(3)括号内的3不是限制存储数据的大小,而是指示显示宽度.显示宽度和数据类型的取值范围是无关的
	显示宽度只用于显示，并不能限制取值范围和占用空间。例如：INT(3)会占用4字节的存储空间，并且允许的最大值不会是999，而是INT整型所允许的最大值。显示宽度只是指明MySQL最大可能显示的数字个数，数值的位数小于指定的宽度时会由空格填充
	*/

	MySqlBIT       = addMySqlColumnType(&ColumnTypeInfo{Name: "BIT", HasLength: false, DataByte: "1", NumberSymbolRangeMin: 0, NumberSymbolRangeMax: 1})
	MySqlTINYINT   = addMySqlColumnType(&ColumnTypeInfo{Name: "TINYINT", HasLength: true, DataByte: "1", NumberSymbolRangeMin: -128, NumberSymbolRangeMax: 127, NumberNoSymbolRangeMin: 0, NumberNoSymbolRangeMax: 255})
	MySqlSMALLINT  = addMySqlColumnType(&ColumnTypeInfo{Name: "SMALLINT", HasLength: true, DataByte: "2", NumberSymbolRangeMin: -32768, NumberSymbolRangeMax: 32767, NumberNoSymbolRangeMin: 0, NumberNoSymbolRangeMax: 65535})
	MySqlMEDIUMINT = addMySqlColumnType(&ColumnTypeInfo{Name: "MEDIUMINT", HasLength: true, DataByte: "3", NumberSymbolRangeMin: -8388608, NumberSymbolRangeMax: 8388607, NumberNoSymbolRangeMin: 0, NumberNoSymbolRangeMax: 16777215})
	MySqlINT       = addMySqlColumnType(&ColumnTypeInfo{Name: "INT", HasLength: true, DataByte: "4", NumberSymbolRangeMin: -2147483648, NumberSymbolRangeMax: 2147483647, NumberNoSymbolRangeMin: 0, NumberNoSymbolRangeMax: 4294967295})
	MySqlBIGINT    = addMySqlColumnType(&ColumnTypeInfo{Name: "BIGINT", HasLength: true, DataByte: "8"})
	MySqlFLOAT     = addMySqlColumnType(&ColumnTypeInfo{Name: "FLOAT", HasLength: true, HasDecimal: true, DataByte: "4"})
	MySqlDOUBLE    = addMySqlColumnType(&ColumnTypeInfo{Name: "DOUBLE", HasLength: true, HasDecimal: true, DataByte: "8"})

	/**
	DECIMAL。浮点数类型和定点数类型都可以用（M，N）来表示。其中，M称为精度，表示总共的位数；N称为标度，表示小数的位数.DECIMAL若不指定精度则默认为(10,0)
	不论是定点数还是浮点数类型，如果用户指定的精度超出精度范围，则会四舍五入
	*/

	MySqlDECIMAL = addMySqlColumnType(&ColumnTypeInfo{Name: "DECIMAL", HasLength: true, DataByte: "length+2"})

	/** 日期/时间类型 **/
	/**
	注意: TIMESTAMP的范围是1970年到2038年
	TIMESTAMP与DATETIME除了存储字节和支持的范围不同外，还有一个最大的区别就是：DATETIME在存储日期数据时，按实际输入的格式存储，即输入什么就存储什么，与时区无关；而TIMESTAMP值的存储是以UTC（世界标准时间）格式保存的，存储时对当前时区进行转换，检索时再转换回当前时区。查询时，不同时区显示的时间值是不同的。

	DATE:
	（1）以‘YYYY-MM-DD’或者‘YYYYMMDD’字符串格式表示的日期，取值范围为‘1000-01-01’～‘9999-12-3’。例如，输入‘2012-12-31’或者‘20121231’，插入数据库的日期都为2012-12-31。
	（2）以‘YY-MM-DD’或者‘YYMMDD’字符串格式表示的日期，在这里YY表示两位的年值。包含两位年值的日期会令人模糊，因为不知道世纪。MySQL使用以下规则解释两位年值：‘00～69’范围的年值转换为‘2000～2069’；‘70～99’范围的年值转换为‘1970～1999’。例如，输入‘12-12-31’，插入数据库的日期为2012-12-31；输入‘981231’，插入数据的日期为1998-12-31。
	（3）以YY-MM-DD或者YYMMDD数字格式表示的日期，与前面相似，00~69范围的年值转换为2000～2069，70～99范围的年值转换为1970～1999。例如，输入12-12-31插入数据库的日期为2012-12-31；输入981231，插入数据的日期为1998-12-31
	*/

	MySqlYEAR      = addMySqlColumnType(&ColumnTypeInfo{Name: "YEAR", DataByte: "1", NumberNoSymbolRangeMin: 1901, NumberNoSymbolRangeMax: 2155})
	MySqlTIME      = addMySqlColumnType(&ColumnTypeInfo{Name: "TIME", DataByte: "3"})
	MySqlDATE      = addMySqlColumnType(&ColumnTypeInfo{Name: "DATE", DataByte: "3"})
	MySqlDATETIME  = addMySqlColumnType(&ColumnTypeInfo{Name: "DATETIME", DataByte: "8"})
	MySqlTIMESTAMP = addMySqlColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", DataByte: "4"})

	/** 字符串类型 **/
	/**
	varchar(M)说明 括号内的M和INT(4)类型的限制不一样,这里M对插入数据的长度有限制,超长就会报错
	varchar字段长度直接按字符计算不区分中英文字符
	*/

	MySqlCHAR       = addMySqlColumnType(&ColumnTypeInfo{Name: "CHAR", HasLength: true, MaxByte: 255})
	MySqlVARCHAR    = addMySqlColumnType(&ColumnTypeInfo{Name: "VARCHAR", HasLength: true, MaxByte: 65535})
	MySqlTINYTEXT   = addMySqlColumnType(&ColumnTypeInfo{Name: "TINYTEXT", HasLength: true, MaxByte: 255})
	MySqlTEXT       = addMySqlColumnType(&ColumnTypeInfo{Name: "TEXT", HasLength: true, MaxByte: 65535})
	MySqlMEDIUMTEXT = addMySqlColumnType(&ColumnTypeInfo{Name: "MEDIUMTEXT", HasLength: true, MaxByte: 16777215})
	MySqlLONGTEXT   = addMySqlColumnType(&ColumnTypeInfo{Name: "LONGTEXT", HasLength: true, MaxByte: 4294967295})
	MySqlENUM       = addMySqlColumnType(&ColumnTypeInfo{Name: "ENUM", HasLength: true, MaxByte: 2})
	MySqlTINYBLOB   = addMySqlColumnType(&ColumnTypeInfo{Name: "TINYBLOB", HasLength: true, MaxByte: 255})
	MySqlBLOB       = addMySqlColumnType(&ColumnTypeInfo{Name: "BLOB", HasLength: true, MaxByte: 65535})
	MySqlMEDIUMBLOB = addMySqlColumnType(&ColumnTypeInfo{Name: "MEDIUMBLOB", HasLength: true, MaxByte: 16777215})
	MySqlLONGBLOB   = addMySqlColumnType(&ColumnTypeInfo{Name: "LONGBLOB", HasLength: true, MaxByte: 4294967295})
)

func addMySqlColumnType(columnTypeInfo *ColumnTypeInfo) *ColumnTypeInfo {
	AppendColumnTypeInfo(DatabaseTypeMySql, columnTypeInfo)
	return columnTypeInfo
}
