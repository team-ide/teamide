package modelcoder

type VariableModel struct {
	Name        string       `json:"name,omitempty"`        // 名称，同一个方法模块唯一
	Comment     string       `json:"comment,omitempty"`     // 注释说明
	ValueScript string       `json:"valueScript,omitempty"` // 值脚本
	IsList      string       `json:"isList,omitempty"`      // 是否是列表
	DataType    string       `json:"dataType,omitempty"`    // 数据类型
	DataStruct  *StructModel `json:"-"`                     // 数据结构体
}

type DataType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	dataTypes = []*DataType{}

	DATA_TYPE_STRING  = newDataType("string", "字符串(String,string)")
	DATA_TYPE_INT     = newDataType("int", "整形(int,int32)")
	DATA_TYPE_LONG    = newDataType("long", "长整型(long,int64)")
	DATA_TYPE_BOOLEAN = newDataType("boolean", "布尔型(boolean,bool)")
	DATA_TYPE_BYTE    = newDataType("byte", "字节型(byte,int8)")
	DATA_TYPE_DATE    = newDataType("date", "字节型(Date,date)")
	DATA_TYPE_SHORT   = newDataType("short", "短整型(short,int16)")
	DATA_TYPE_DOUBLE  = newDataType("double", "双精度浮点型(double,float64)")
	DATA_TYPE_FLOAT   = newDataType("float", "浮点型(float,float32)")
	DATA_TYPE_MAP     = newDataType("map", "集合(Map,map)")
)

func newDataType(value, text string) *DataType {
	res := &DataType{
		Value: value,
		Text:  text,
	}
	dataTypes = append(dataTypes, res)
	return res
}

func GetDataType(value string) *DataType {
	for _, one := range dataTypes {
		if one.Value == value {
			return one
		}
	}
	return nil
}
