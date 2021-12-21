package modelcoder

type ParamModel struct {
	Name        string `json:"name,omitempty"`        // 名称，同一个方法模块唯一
	Comment     string `json:"comment,omitempty"`     // 注释说明
	ValueScript string `json:"valueScript,omitempty"` // 值脚本
	IsList      string `json:"isList,omitempty"`      // 是否是列表
	DataType    string `json:"dataType,omitempty"`    // 数据类型
}

type HttpParamModel struct {
	Name        string `json:"name,omitempty"`        // 名称，同一个方法模块唯一
	Comment     string `json:"comment,omitempty"`     // 注释说明
	ValueScript string `json:"valueScript,omitempty"` // 值脚本
	IsList      string `json:"isList,omitempty"`      // 是否是列表
	DataType    string `json:"dataType,omitempty"`    // 数据类型
	DataPlace   string `json:"dataPlace,omitempty"`   // 数据位置
}

type ParamData struct {
	Name string      `json:"name,omitempty"` // 参数名称
	Data interface{} `json:"data,omitempty"` // 参数值
}

type ParamModelDataType struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	paramModelDataTypes []*ParamModelDataType

	PARAM_STRING  = newParamModelDataType("string", "字符串(String,string)")
	PARAM_INT     = newParamModelDataType("int", "整形(int,int32)")
	PARAM_LONG    = newParamModelDataType("long", "长整型(long,int64)")
	PARAM_BOOLEAN = newParamModelDataType("boolean", "布尔型(boolean,bool)")
	PARAM_BYTE    = newParamModelDataType("byte", "字节型(byte,int8)")
	PARAM_DATE    = newParamModelDataType("date", "字节型(Date,date)")
	PARAM_SHORT   = newParamModelDataType("short", "短整型(short,int16)")
	PARAM_DOUBLE  = newParamModelDataType("double", "双精度浮点型(double,float64)")
	PARAM_FLOAT   = newParamModelDataType("float", "浮点型(float,float32)")
	PARAM_MAP     = newParamModelDataType("map", "集合(Map,map)")
)

func newParamModelDataType(value, text string) *ParamModelDataType {
	res := &ParamModelDataType{
		Value: value,
		Text:  text,
	}
	paramModelDataTypes = append(paramModelDataTypes, res)
	return res
}

func GetParamModelDataType(value string) *ParamModelDataType {
	for _, one := range paramModelDataTypes {
		if one.Value == value {
			return one
		}
	}
	return nil
}
