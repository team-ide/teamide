package module_metadata

import (
	"fmt"
	"teamide/internal/base"
	"teamide/pkg/util"
)

// MetadataField 元数据结构体字段
type MetadataField struct {
	ParentCode   int              `json:"parentCode,omitempty"`
	Code         int              `json:"code,omitempty"`
	Name         string           `json:"name,omitempty"`
	Comment      string           `json:"comment,omitempty"`
	FieldType    int              `json:"fieldType,omitempty"`
	DefaultValue string           `json:"defaultValue,omitempty"`
	Options      []base.OBean     `json:"options,omitempty"`
	Fields       []*MetadataField `json:"fields,omitempty"`
}

const (
	// 字段类型，定义的值不可删除，不可修改，不可重复

	FTString     int = 1
	FTInt8       int = 2
	FTInt        int = 3
	FTInt64      int = 4
	FTBool       int = 5
	FTFloat      int = 6
	FTStruct     int = 7
	FTListStruct int = 8
)

var (
	MetadataFieldCache     []*MetadataField
	MetadataFieldCodeCache = map[int]*MetadataField{}

	// UMPStruct 用户画像元数据
	UMPStruct = AppendMetadata("persona", "画像", 10001, UMPName, UMPAge, UMPSex, UMPPhoto)
	UMPName   = NewMetadataField("name", "姓名", 1000110001, FTString, "")
	UMPAge    = NewMetadataField("age", "年龄", 1000110002, FTInt, "")
	UMPSex    = NewMetadataField("sex", "性别", 1000110003, FTInt8, "", base.NewOBean("男", 1), base.NewOBean("女", 2))
	UMPPhoto  = NewMetadataField("photo", "照片", 1000110004, FTString, "")

	// UMEStruct 用户企业元数据
	UMEStruct = AppendMetadata("enterprise", "企业信息", 10002, UMEId, UMEName, UMEOrgList)
	UMEId     = NewMetadataField("id", "企业编号", 1000210001, FTInt64, "")
	UMEName   = NewMetadataField("name", "企业名称", 1000210002, FTString, "")

	UMEOrgList      = NewMetadataStruct("orgList", "部门", 10003, UMEOrgId, UMEOrgName, UMEOrgCodeName, UMEPositionName)
	UMEOrgId        = NewMetadataField("id", "部门编号", 1000310001, FTInt64, "")
	UMEOrgName      = NewMetadataField("name", "部门名称", 1000310002, FTString, "")
	UMEOrgCodeName  = NewMetadataField("code", "部门编码", 1000310003, FTString, "")
	UMEPositionName = NewMetadataField("position", "部门职位", 1000310004, FTString, "")
)

func AppendMetadataField(field *MetadataField) {
	checkMetadata(field, false)
	MetadataFieldCache = append(MetadataFieldCache, field)
}

func AppendMetadata(name string, comment string, code int, fields ...*MetadataField) *MetadataField {
	obj := NewMetadataStruct(name, comment, code, fields...)
	AppendMetadataField(obj)
	return obj
}

func NewMetadataStruct(name string, comment string, code int, fields ...*MetadataField) *MetadataField {
	res := &MetadataField{}
	res.Name = name
	res.Comment = comment
	res.Code = code
	if len(fields) > 0 {
		res.Fields = []*MetadataField{}
		for _, one := range fields {
			one.ParentCode = res.Code
			res.Fields = append(res.Fields, one)
		}
	}
	return res
}

func NewMetadataField(name string, comment string, code int, fieldType int, defaultValue string, options ...base.OBean) *MetadataField {
	res := &MetadataField{}
	res.Name = name
	res.Comment = comment
	res.Code = code
	res.FieldType = fieldType
	res.DefaultValue = defaultValue

	if len(options) > 0 {
		res.Options = []base.OBean{}
		res.Options = append(res.Options, options...)
	}
	return res
}

func checkMetadata(field *MetadataField, out bool) {
	if out {
		fmt.Println("----------字段（", field.Comment, "）----------")
		fmt.Println(util.ToJSON(field))
	}
	_, ok := MetadataFieldCodeCache[field.Code]
	if ok {
		panic(fmt.Sprint("元数据编码:", field.Code, "已存在!"))
	}
	MetadataFieldCodeCache[field.Code] = field
	if len(field.Fields) > 0 {
		for _, one := range field.Fields {
			checkMetadata(one, out)
		}
	}
}
