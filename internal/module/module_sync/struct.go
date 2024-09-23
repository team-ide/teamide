package module_sync

import "github.com/team-ide/go-tool/thrift"
import thriftlib "github.com/apache/thrift/lib/go/thrift"

var (
	UserSetting = &thrift.Struct{
		Fields: []*thrift.Field{
			{Name: "option", Num: 1, Type: &thrift.FieldType{
				TypeId:       thriftlib.MAP,
				MapKeyType:   thrift.GetFieldType(thriftlib.STRING),
				MapValueType: thrift.GetFieldType(thriftlib.STRING),
			}},
		},
	}

	ToolboxGroup = &thrift.Struct{
		Fields: []*thrift.Field{
			{Name: "groupId", Num: 1, Type: thrift.GetFieldType(thriftlib.I64)},
			{Name: "name", Num: 2, Type: thrift.GetFieldType(thriftlib.STRING)},
			{Name: "comment", Num: 3, Type: thrift.GetFieldType(thriftlib.STRING)},
			{Name: "option", Num: 4, Type: thrift.GetFieldType(thriftlib.STRING)},
			{Name: "sequence", Num: 5, Type: thrift.GetFieldType(thriftlib.I32)},
		},
	}

	Toolbox = &thrift.Struct{
		Fields: []*thrift.Field{
			{Name: "toolboxId", Num: 1, Type: thrift.GetFieldType(thriftlib.I64)},
			{Name: "toolboxType", Num: 2, Type: thrift.GetFieldType(thriftlib.STRING)},
			{Name: "groupId", Num: 3, Type: thrift.GetFieldType(thriftlib.I64)},
			{Name: "name", Num: 4, Type: thrift.GetFieldType(thriftlib.STRING)},
			{Name: "comment", Num: 5, Type: thrift.GetFieldType(thriftlib.STRING)},
			{Name: "option", Num: 6, Type: thrift.GetFieldType(thriftlib.STRING)},
			{Name: "visibility", Num: 7, Type: thrift.GetFieldType(thriftlib.I32)},
			{Name: "sequence", Num: 8, Type: thrift.GetFieldType(thriftlib.I32)},
		},
	}

	ToolboxExtend = &thrift.Struct{
		Fields: []*thrift.Field{
			{Name: "extendId", Num: 1, Type: thrift.GetFieldType(thriftlib.I64)},
			{Name: "toolboxId", Num: 2, Type: thrift.GetFieldType(thriftlib.I64)},
			{Name: "extendType", Num: 3, Type: thrift.GetFieldType(thriftlib.STRING)},
			{Name: "name", Num: 4, Type: thrift.GetFieldType(thriftlib.STRING)},
			{Name: "value", Num: 5, Type: thrift.GetFieldType(thriftlib.STRING)},
		},
	}
)
