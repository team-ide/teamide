package base

import "fmt"

// 元数据结构体
type MStruct struct {
	StructCode int        `json:"structCode" column:"structCode"`
	Name       string     `json:"name" column:"name"`
	Comment    string     `json:"comment" column:"comment"`
	Fields     []*MSField `json:"fields" column:"fields"`
}

// 元数据结构体字段
type MSField struct {
	StructCode      int      `json:"structCode" column:"structCode"`
	StructFieldCode int      `json:"structFieldCode" column:"structFieldCode"`
	Struct          *MStruct `json:"struct" column:"struct"`
	Name            string   `json:"name" column:"name"`
	Comment         string   `json:"comment" column:"comment"`
	FieldType       int      `json:"fieldType" column:"fieldType"`
	DefaultValue    string   `json:"defaultValue" column:"defaultValue"`
	Options         []OBean  `json:"options" column:"options"`
}

const (
	// 字段类型，定义的值不可删除，不可修改，不可重复
	F_T_STRING int = (1 + iota)
	F_T_INT8
	F_T_INT
	F_T_INT64
	F_T_BOOL
	F_T_FLOAT
	F_T_STRUCT
	F_T_LIST_STRUCT
)

const (
	// 用户元数据Code定义，定义的值不可删除，不可修改，不可重复
	// 用户画像
	U_M_P_C      int = 10001
	U_M_P_NAME_C int = U_M_P_C*10000 + (1 + iota)
	U_M_P_AGE_C
	U_M_P_SEX_C
	U_M_P_PHOTO_C
	U_M_P_HEIGHT_C
	U_M_P_WEIGHT_C

	// 用户企业信息
	U_M_E_C      int = 10002
	U_M_E_NAME_C int = U_M_E_C*10000 + (1 + iota)
	U_M_E_SALARY_C
	U_M_E_ORG_C
	U_M_E_ORG_NAME_C
	U_M_E_ORG_CODE_C
	U_M_E_ORG_POSITION_C
)

var (
	// 用户画像元数据
	U_M_P        *MStruct = &MStruct{U_M_P_C, "persona", "画像", []*MSField{U_M_P_NAME, U_M_P_AGE, U_M_P_SEX, U_M_P_HEIGHT, U_M_P_WEIGHT}}
	U_M_P_NAME   *MSField = &MSField{U_M_P_C, U_M_P_NAME_C, nil, "name", "姓名", F_T_STRING, "", nil}
	U_M_P_AGE    *MSField = &MSField{U_M_P_C, U_M_P_AGE_C, nil, "age", "年龄", F_T_INT, "", nil}
	U_M_P_SEX    *MSField = &MSField{U_M_P_C, U_M_P_SEX_C, nil, "sex", "性别", F_T_INT8, "", []OBean{NewOBean("男", 1), NewOBean("女", 2)}}
	U_M_P_PHOTO  *MSField = &MSField{U_M_P_C, U_M_P_PHOTO_C, nil, "photo", "照片", F_T_STRING, "", nil}
	U_M_P_HEIGHT *MSField = &MSField{U_M_P_C, U_M_P_HEIGHT_C, nil, "height", "身高", F_T_FLOAT, "", nil}
	U_M_P_WEIGHT *MSField = &MSField{U_M_P_C, U_M_P_WEIGHT_C, nil, "weight", "体重", F_T_FLOAT, "", nil}

	// 用户企业元数据
	U_M_E        *MStruct = &MStruct{U_M_E_C, "enterprise", "企业信息", []*MSField{U_M_E_NAME, U_M_E_SALARY, U_M_E_ORG_}}
	U_M_E_NAME   *MSField = &MSField{U_M_E_C, U_M_E_NAME_C, nil, "name", "企业名称", F_T_STRING, "", nil}
	U_M_E_SALARY *MSField = &MSField{U_M_E_C, U_M_E_SALARY_C, nil, "salary", "薪资", F_T_FLOAT, "", nil}
	U_M_E_ORG_   *MSField = &MSField{U_M_E_C, U_M_E_ORG_C, U_M_E_ORG, "orgs", "部门", F_T_LIST_STRUCT, "", nil}

	// 用户企业部门
	U_M_E_ORG           *MStruct = &MStruct{U_M_E_C, "orgs", "部门", []*MSField{U_M_E_ORG_NAME, U_M_E_CODE_NAME, U_M_E_POSITION_NAME}}
	U_M_E_ORG_NAME      *MSField = &MSField{U_M_E_C, U_M_E_ORG_NAME_C, nil, "name", "部门名称", F_T_STRING, "", nil}
	U_M_E_CODE_NAME     *MSField = &MSField{U_M_E_C, U_M_E_ORG_CODE_C, nil, "code", "部门编码", F_T_STRING, "", nil}
	U_M_E_POSITION_NAME *MSField = &MSField{U_M_E_C, U_M_E_ORG_POSITION_C, nil, "position", "部门职位", F_T_STRING, "", nil}

	U_M []*MStruct = []*MStruct{U_M_P, U_M_E}
)

func init() {

	checkMetadata(true)
}
func checkMetadata(out bool) {
	codes := map[int]bool{}
	for _, obj := range U_M {
		if out {
			fmt.Println("----------结构体（", obj.Comment, "）----------")
			fmt.Println(ToJSON(obj))
		}
		_, ok := codes[obj.StructCode]
		if ok {
			panic(fmt.Sprint("元数据编码：", obj.StructCode, "已存在！"))
		}
		codes[obj.StructCode] = true
		if out {
			fmt.Println("----------结构体（", obj.Comment, "）属性----------")
		}
		for _, field := range obj.Fields {
			if out {
				fmt.Println(ToJSON(field))
			}
			_, ok := codes[field.StructFieldCode]
			if ok {
				panic(fmt.Sprint("元数据编码：", field.StructFieldCode, "已存在！"))
			}
			codes[field.StructFieldCode] = true
		}
	}
}
