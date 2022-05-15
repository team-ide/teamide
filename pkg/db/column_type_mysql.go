package db

var (

	/** 数值类型 **/
	/**
	MySQL 支持所有标准 SQL 数值数据类型。
	这些类型包括严格数值数据类型(INTEGER、SMALLINT、DECIMAL 和 NUMERIC)，以及近似数值数据类型(FLOAT、REAL 和 DOUBLE PRECISION)。
	关键字INT是INTEGER的同义词，关键字DEC是DECIMAL的同义词。
	BIT数据类型保存位字段值，并且支持 MyISAM、MEMORY、InnoDB 和 BDB表。
	作为 SQL 标准的扩展，MySQL 也支持整数类型 TINYINT、MEDIUMINT 和 BIGINT。下面的表显示了需要的每个整数类型的存储和范围。

	如果不设置长度，会有默认的长度
	长度代表了显示的最大宽度，如果不够会用0在左边填充，但必须搭配zerofill 使用！
	例如：
	INT(7) 括号中7不是指范围，范围是由数据类型决定的，只是代表显示结果的宽度
	*/

	MySqlBIT       = addMySqlColumnType(&ColumnTypeInfo{Name: "BIT", TypeFormat: "BIT($l)", HasLength: false, IsNumber: true})
	MySqlTINYINT   = addMySqlColumnType(&ColumnTypeInfo{Name: "TINYINT", TypeFormat: "TINYINT($l)", HasLength: true, IsNumber: true})
	MySqlSMALLINT  = addMySqlColumnType(&ColumnTypeInfo{Name: "SMALLINT", TypeFormat: "SMALLINT($l)", HasLength: true, IsNumber: true})
	MySqlMEDIUMINT = addMySqlColumnType(&ColumnTypeInfo{Name: "MEDIUMINT", TypeFormat: "MEDIUMINT($l)", HasLength: true, IsNumber: true})
	MySqlINT       = addMySqlColumnType(&ColumnTypeInfo{Name: "INT", TypeFormat: "INT($l)", HasLength: true, IsNumber: true})
	MySqlINTEGER   = addMySqlColumnType(&ColumnTypeInfo{Name: "INTEGER", TypeFormat: "INTEGER($l)", HasLength: true, IsNumber: true})
	MySqlBIGINT    = addMySqlColumnType(&ColumnTypeInfo{Name: "BIGINT", TypeFormat: "BIGINT($l)", HasLength: true, IsNumber: true})

	/** 小数 **/

	/**
	M：整数部位+小数部位
	D：小数部位
	如果超过范围，则插入临界值
	M和D都可以省略
	如果是DECIMAL，则M默认为10，D默认为0
	如果是FLOAT和DOUBLE，则会根据插入的数值的精度来决定精度
	定点型的精确度较高，如果要求插入数值的精度较高如货币运算等则考虑使用
	原则：所选择的类型越简单越好，能保存数值的类型越小越好
	*/

	MySqlFLOAT  = addMySqlColumnType(&ColumnTypeInfo{Name: "FLOAT", TypeFormat: "FLOAT($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})
	MySqlDOUBLE = addMySqlColumnType(&ColumnTypeInfo{Name: "DOUBLE", TypeFormat: "DOUBLE($l, $d)", HasLength: true, HasDecimal: true, IsNumber: true})

	/**
	DECIMAL。浮点数类型和定点数类型都可以用（M，N）来表示。其中，M称为精度，表示总共的位数；N称为标度，表示小数的位数.DECIMAL若不指定精度则默认为(10,0)
	不论是定点数还是浮点数类型，如果用户指定的精度超出精度范围，则会四舍五入
	*/

	MySqlDEC     = addMySqlColumnType(&ColumnTypeInfo{Name: "DEC", TypeFormat: "DEC($l, $d)", HasLength: true, IsNumber: true})
	MySqlDECIMAL = addMySqlColumnType(&ColumnTypeInfo{Name: "DECIMAL", TypeFormat: "DOUBLE($l, $d)", HasLength: true, IsNumber: true})

	MySqlNUMBER = addMySqlColumnType(&ColumnTypeInfo{Name: "NUMBER", TypeFormat: "NUMBER($l, $d)", HasLength: true, IsNumber: true})

	/** 日期/时间类型 **/
	/**
	表示时间值的日期和时间类型为DATETIME、DATE、TIMESTAMP、TIME和YEAR。
	每个时间类型有一个有效值范围和一个"零"值，当指定不合法的MySQL不能表示的值时使用"零"值。
	TIMESTAMP类型有专有的自动更新特性，将在后面描述。
	DATE:
	（1）以‘YYYY-MM-DD’或者‘YYYYMMDD’字符串格式表示的日期，取值范围为‘1000-01-01’～‘9999-12-3’。例如，输入‘2012-12-31’或者‘20121231’，插入数据库的日期都为2012-12-31。
	（2）以‘YY-MM-DD’或者‘YYMMDD’字符串格式表示的日期，在这里YY表示两位的年值。包含两位年值的日期会令人模糊，因为不知道世纪。MySQL使用以下规则解释两位年值：‘00～69’范围的年值转换为‘2000～2069’；‘70～99’范围的年值转换为‘1970～1999’。例如，输入‘12-12-31’，插入数据库的日期为2012-12-31；输入‘981231’，插入数据的日期为1998-12-31。
	（3）以YY-MM-DD或者YYMMDD数字格式表示的日期，与前面相似，00~69范围的年值转换为2000～2069，70～99范围的年值转换为1970～1999。例如，输入12-12-31插入数据库的日期为2012-12-31；输入981231，插入数据的日期为1998-12-31
	*/

	MySqlYEAR      = addMySqlColumnType(&ColumnTypeInfo{Name: "YEAR", TypeFormat: "YEAR", IsDateTime: true})
	MySqlTIME      = addMySqlColumnType(&ColumnTypeInfo{Name: "TIME", TypeFormat: "TIME", IsDateTime: true})
	MySqlDATE      = addMySqlColumnType(&ColumnTypeInfo{Name: "DATE", TypeFormat: "DATE", IsDateTime: true})
	MySqlDATETIME  = addMySqlColumnType(&ColumnTypeInfo{Name: "DATETIME", TypeFormat: "DATETIME", IsDateTime: true})
	MySqlTIMESTAMP = addMySqlColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", TypeFormat: "TIMESTAMP", IsDateTime: true})

	/** 字符串类型 **/
	/**
	字符串类型指CHAR、VARCHAR、BINARY、VARBINARY、BLOB、TEXT、ENUM和SET。该节描述了这些类型如何工作以及如何在查询中使用这些类型

	注意：char(n) 和 varchar(n) 中括号中 n 代表字符的个数，并不代表字节个数，比如 CHAR(30) 就可以存储 30 个字符。
	CHAR 和 VARCHAR 类型类似，但它们保存和检索的方式不同。它们的最大长度和是否尾部空格被保留等方面也不同。在存储或检索过程中不进行大小写转换。
	BINARY 和 VARBINARY 类似于 CHAR 和 VARCHAR，不同的是它们包含二进制字符串而不要非二进制字符串。也就是说，它们包含字节字符串而不是字符字符串。这说明它们没有字符集，并且排序和比较基于列值字节的数值值。
	BLOB 是一个二进制大对象，可以容纳可变数量的数据。有 4 种 BLOB 类型：TINYBLOB、BLOB、MEDIUMBLOB 和 LONGBLOB。它们区别在于可容纳存储范围不同。
	有 4 种 TEXT 类型：TINYTEXT、TEXT、MEDIUMTEXT 和 LONGTEXT。对应的这 4 种 BLOB 类型，可存储的最大长度不同，可根据实际情况选择。
	*/

	MySqlCHAR       = addMySqlColumnType(&ColumnTypeInfo{Name: "CHAR", TypeFormat: "CHAR($l)", HasLength: true, IsString: true})
	MySqlVARCHAR    = addMySqlColumnType(&ColumnTypeInfo{Name: "VARCHAR", TypeFormat: "VARCHAR($l)", HasLength: true, IsString: true})
	MySqlTINYTEXT   = addMySqlColumnType(&ColumnTypeInfo{Name: "TINYTEXT", TypeFormat: "TINYTEXT($l)", HasLength: true, IsString: true})
	MySqlTEXT       = addMySqlColumnType(&ColumnTypeInfo{Name: "TEXT", TypeFormat: "TEXT($l)", HasLength: true, IsString: true})
	MySqlMEDIUMTEXT = addMySqlColumnType(&ColumnTypeInfo{Name: "MEDIUMTEXT", TypeFormat: "MEDIUMTEXT($l)", HasLength: true, IsString: true})
	MySqlLONGTEXT   = addMySqlColumnType(&ColumnTypeInfo{Name: "LONGTEXT", TypeFormat: "LONGTEXT($l)", HasLength: true, IsString: true})
	MySqlENUM       = addMySqlColumnType(&ColumnTypeInfo{Name: "ENUM", TypeFormat: "ENUM($l)", HasLength: true, IsString: true})
	MySqlTINYBLOB   = addMySqlColumnType(&ColumnTypeInfo{Name: "TINYBLOB", TypeFormat: "TINYBLOB($l)", HasLength: true, IsString: true})
	MySqlBLOB       = addMySqlColumnType(&ColumnTypeInfo{Name: "BLOB", TypeFormat: "BLOB($l)", HasLength: true, IsString: true})
	MySqlMEDIUMBLOB = addMySqlColumnType(&ColumnTypeInfo{Name: "MEDIUMBLOB", TypeFormat: "MEDIUMBLOB($l)", HasLength: true, IsString: true})
	MySqlLONGBLOB   = addMySqlColumnType(&ColumnTypeInfo{Name: "LONGBLOB", TypeFormat: "LONGBLOB($l)", HasLength: true, IsString: true})

	MySqlSET = addMySqlColumnType(&ColumnTypeInfo{Name: "SET", TypeFormat: "SET($l)", HasLength: true, IsString: true})
)

var (
	MySqlColumnTypeInfos []*ColumnTypeInfo
)

func addMySqlColumnType(columnTypeInfo *ColumnTypeInfo) *ColumnTypeInfo {
	MySqlColumnTypeInfos = append(MySqlColumnTypeInfos, columnTypeInfo)
	AppendColumnTypeInfo(DatabaseTypeMySql, columnTypeInfo)
	return columnTypeInfo
}
