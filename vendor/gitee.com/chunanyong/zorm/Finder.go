package zorm

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

//Finder 查询数据库的载体,所有的sql语句都要通过Finder执行.
//Finder To query the database carrier, all SQL statements must be executed through Finder
type Finder struct {
	//拼接SQL
	//Splicing SQL.
	sqlBuilder strings.Builder
	//SQL的参数值
	//SQL parameter values.
	values []interface{}
	//注入检查,默认true 不允许SQL注入的 ' 单引号
	//Injection check, default true does not allow SQL injection  single quote
	InjectionCheck bool
	//CountFinder 自定义的查询总条数'Finder',使用指针默认为nil.主要是为了在'group by'等复杂情况下,为了性能,手动编写总条数语句
	//CountFinder The total number of custom queries is'Finder', and the pointer is nil by default. It is mainly used to manually write the total number of statements for performance in complex situations such as'group by'
	CountFinder *Finder
	//是否自动查询总条数,默认true.同时需要Page不为nil,才查询总条数
	//Whether to automatically query the total number of entries, the default is true. At the same time, the Page is not nil to query the total number of entries
	SelectTotalCount bool
	//SQL语句
	//SQL statement
	sqlstr string
}

//NewFinder  初始化一个Finder,生成一个空的Finder
//NewFinder Initialize a Finder and generate an empty Finder
func NewFinder() *Finder {
	finder := Finder{}
	finder.SelectTotalCount = true
	finder.InjectionCheck = true
	finder.values = make([]interface{}, 0)
	return &finder
}

//NewSelectFinder 根据表名初始化查询的Finder | Finder that initializes the query based on the table name
//NewSelectFinder("tableName") SELECT * FROM tableName
//NewSelectFinder("tableName", "id,name") SELECT id,name FROM tableName
func NewSelectFinder(tableName string, strs ...string) *Finder {
	finder := NewFinder()
	finder.sqlBuilder.WriteString("SELECT ")
	if len(strs) > 0 {
		for _, str := range strs {
			finder.sqlBuilder.WriteString(str)
		}
	} else {
		finder.sqlBuilder.WriteString("*")
	}
	finder.sqlBuilder.WriteString(" FROM ")
	finder.sqlBuilder.WriteString(tableName)
	return finder
}

//NewUpdateFinder 根据表名初始化更新的Finder,  UPDATE tableName SET
//NewUpdateFinder Initialize the updated Finder according to the table name, UPDATE tableName SET
func NewUpdateFinder(tableName string) *Finder {
	finder := NewFinder()
	finder.sqlBuilder.WriteString("UPDATE ")
	finder.sqlBuilder.WriteString(tableName)
	finder.sqlBuilder.WriteString(" SET ")
	return finder
}

//NewDeleteFinder 根据表名初始化删除的'Finder',  DELETE FROM tableName
//NewDeleteFinder Finder for initial deletion based on table name. DELETE FROM tableName
func NewDeleteFinder(tableName string) *Finder {
	finder := NewFinder()
	finder.sqlBuilder.WriteString("DELETE FROM ")
	finder.sqlBuilder.WriteString(tableName)
	//所有的 WHERE 都不加,规则统一,好记
	//No WHERE is added, the rules are unified, easy to remember
	//finder.sqlBuilder.WriteString(" WHERE ")
	return finder
}

//Append 添加SQL和参数的值,第一个参数是语句,后面的参数[可选]是参数的值,顺序要正确
//例如: finder.Append(" and id=? and name=? ",23123,"abc")
//只拼接SQL,例如: finder.Append(" and name=123 ")
//Append:Add SQL and parameter values, the first parameter is the statement, and the following parameter (optional) is the value of the parameter, in the correct order
//E.g:  finder.Append(" and id=? and name=? ",23123,"abc")
//Only splice SQL, E.g: finder.Append(" and name=123 ")
func (finder *Finder) Append(s string, values ...interface{}) *Finder {

	//不要自己构建finder,使用Newxxx方法
	//Don't build finder by yourself, use Newxxx method
	if finder.values == nil {
		return nil
	}

	if len(s) > 0 {
		if len(finder.sqlstr) > 0 {
			finder.sqlstr = ""
		}
		//默认加一个空格,避免手误两个字符串连接再一起
		//A space is added by default to avoid hand mistakes when connecting two strings together
		finder.sqlBuilder.WriteString(" ")

		finder.sqlBuilder.WriteString(s)

	}
	if values == nil || len(values) < 1 {
		return finder
	}
	//for _, v := range values {
	//	finder.Values = append(finder.Values, v)
	//}
	finder.values = append(finder.values, values...)
	return finder
}

//AppendFinder 添加另一个Finder finder.AppendFinder(f)
//AppendFinder Add another Finder . finder.AppendFinder(f)
func (finder *Finder) AppendFinder(f *Finder) (*Finder, error) {
	if f == nil {
		return nil, errors.New("finder-->AppendFinder参数是nil")
	}

	//不要自己构建finder,使用Newxxx方法
	//Don't build finder by yourself, use Newxxx method
	if finder.values == nil {
		return nil, errors.New("finder-->AppendFinder不要自己构建finder,使用Newxxx方法")
	}

	//添加f的SQL
	//SQL to add f。
	sqlstr, err := f.GetSQL()
	if err != nil {
		return nil, err
	}
	finder.sqlstr = ""
	finder.sqlBuilder.WriteString(sqlstr)
	//添加f的值
	//Add the value of f
	finder.values = append(finder.values, f.values...)
	return finder, nil
}

//GetSQL 返回Finder封装的SQL语句
//GetSQL Return the SQL statement encapsulated by the Finder
func (finder *Finder) GetSQL() (string, error) {
	//不要自己构建finder,使用Newxxx方法
	//Don't build finder by yourself, use Newxxx method
	if finder.values == nil {
		return "", errors.New("finder-->GetSQL不要自己构建finder,使用Newxxx方法")
	}
	if len(finder.sqlstr) > 0 {
		return finder.sqlstr, nil
	}
	sqlstr := finder.sqlBuilder.String()
	finder.sqlstr = sqlstr
	//包含单引号,属于非法字符串
	//Contains single quotes, which are illegal strings
	if finder.InjectionCheck && (strings.Contains(sqlstr, "'")) {
		return "", errors.New("finder-->GetSQL SQL语句请不要直接拼接字符串参数!!!使用标准的占位符实现,例如  finder.Append(' and id=? and name=? ','123','abc')")
	}

	//处理sql语句中的in,实际就是把数组变量展开,例如 id in(?) ["1","2","3"] 语句变更为 id in (?,?,?) 参数也展开到参数数组里
	//这里认为 slice类型的参数就是in
	//Processing the in in the SQL statement is actually expanding the array variables,
	//for example, id in(?) ["1","2","3"] The statement is changed to id in (?,?,?)
	//The parameters are also expanded to the parameters In the array
	//It is considered that the parameter of the slice type is in
	if finder.values == nil || len(finder.values) < 1 { //如果没有参数
		return sqlstr, nil
	}

	//?问号切割的数组
	//Question mark cut array
	questions := strings.Split(sqlstr, "?")

	//语句中没有?问号
	//No in the sentence Question mark
	if len(questions) < 1 {
		return sqlstr, nil
	}

	//重新记录参数值
	//Re-record the parameter value
	newValues := make([]interface{}, 0)
	//新的sql
	//new sql
	var newSQLStr strings.Builder
	//问号切割的语句实际长度比问号号个数多1个,先把第一个语句片段加上,后面就是比参数的索引大1
	//The actual length of the sentence cut by the question mark is one more than the number of question marks. First, add the first sentence fragment, and the latter is one greater than the index of the parameter
	newSQLStr.WriteString(questions[0])

	//遍历所有的参数
	//Traverse all parameters
	for i, v := range finder.values {
		//先拼接问号,问号切割之后,问号就丢失了,先补充上
		//First splicing the question mark, after the question mark is cut, the question mark is lost, add it first
		newSQLStr.WriteString("?")

		valueOf := reflect.ValueOf(v)
		typeOf := reflect.TypeOf(v)
		kind := valueOf.Kind()
		//如果参数是个指针类型
		//If the parameter is a pointer type
		if kind == reflect.Ptr { //如果是指针 ｜ If it is a pointer
			valueOf = valueOf.Elem()
			typeOf = typeOf.Elem()
			kind = valueOf.Kind()
		}

		//如果不是数组或者slice
		//If it is not an array or slice
		if !(kind == reflect.Array || kind == reflect.Slice) {
			//记录新值
			//Record new value.
			newValues = append(newValues, v)
			//记录SQL
			//Log SQL。
			newSQLStr.WriteString(questions[i+1])
			continue
		}
		//字节数组是特殊的情况
		//Byte array is a special case
		if typeOf == reflect.TypeOf([]byte{}) {
			//记录新值
			//Record new value
			newValues = append(newValues, v)
			//记录SQL
			//Log SQL
			newSQLStr.WriteString(questions[i+1])
			continue
		}

		//如果不是字符串类型的值,无法取长度,这个是个bug,先注释了
		//获取数组类型参数值的长度
		//If it is not a string type value, the length cannot be taken, this is a bug, first comment
		//Get the length of the array type parameter value
		sliceLen := valueOf.Len()
		//数组类型的参数长度小于1,认为是有异常的参数
		//The parameter length of the array type is less than 1, which is considered to be an abnormal parameter
		if sliceLen < 1 {
			return sqlstr, errors.New("finder-->GetSQL语句:" + sqlstr + ",第" + strconv.Itoa(i+1) + "个参数,类型是Array或者Slice,值的长度为0,请检查sql参数有效性")
		}

		for j := 0; j < sliceLen; j++ {
			//每多一个参数,对应",?" 两个符号.增加的问号长度总计是(sliceLen-1)*2
			//Every additional parameter, correspond ",?" ,The total length of the increased question mark is (sliceLen-1)*2
			if j >= 1 {
				//记录SQL
				//Log SQL.
				newSQLStr.WriteString(",?")
			}
			//记录新值
			//Record new value
			sliceValue := valueOf.Index(j).Interface()
			newValues = append(newValues, sliceValue)
		}
		//记录SQL
		//Log SQL
		newSQLStr.WriteString(questions[i+1])
	}
	//重新赋值
	//Reassign
	finder.sqlstr = newSQLStr.String()
	finder.values = newValues
	return finder.sqlstr, nil
}
