package module_sync

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/thrift"
	"github.com/team-ide/go-tool/util"
	yaml "gopkg.in/yaml.v3"
	"reflect"
	"sort"
	"strings"
)

type SyncInfo struct {
	Explain               string           `json:"explain" yaml:"说明" sign:"explain"`
	CreateBy              string           `json:"createBy" yaml:"所属" sign:"createBy"`
	CreateAt              string           `json:"createAt" yaml:"时间" sign:"createAt"`
	Encrypt               string           `json:"encrypt" yaml:"加密" sign:"encrypt"`
	SignKey               string           `json:"signKey" yaml:"密串" sign:"signKey"`
	Sign                  string           `json:"sign" yaml:"签名"`
	UserSettingText       string           `json:"-" yaml:"个人设置" sign:"userSettingText"`
	ToolboxGroupListText  string           `json:"-" yaml:"工具分组" sign:"toolboxGroupListText"`
	ToolboxListText       string           `json:"-" yaml:"工具" sign:"toolboxListText"`
	ToolboxExtendListText string           `json:"-" yaml:"工具扩展" sign:"toolboxExtendListText"`
	UserSetting           map[string]any   `json:"-"`
	UserSettingSize       int              `json:"userSettingSize"`
	ToolboxGroupList      []map[string]any `json:"-"`
	ToolboxGroupSize      int              `json:"toolboxGroupSize"`
	ToolboxList           []map[string]any `json:"-"`
	ToolboxSize           int              `json:"toolboxSize"`
	ToolboxExtendList     []map[string]any `json:"-"`
	ToolboxExtendSize     int              `json:"toolboxExtendSize"`
}

func (this_ *SyncInfo) GetPassword(key string) (res []byte) {
	// 取 MD5 的 8 ~ 24 作为密码
	return []byte(strings.ToUpper(util.GetMD5(key + this_.SignKey)[8:24]))
}

func (this_ *SyncInfo) GenSign(key string) (sign string) {
	signData := GetSignData(this_, true)
	signStr := JoinStringsInASCII(signData, "&", false, true, "sign")
	signStr += "&key=" + key
	signStr = strings.ReplaceAll(signStr, " ", "")
	signStr = strings.ReplaceAll(signStr, "\n", "")
	//fmt.Println("signStr:", signStr)
	sign = strings.ToUpper(util.GetMD5(signStr))
	return
}

func Read(key string, content string) (info *SyncInfo, err error) {
	info = &SyncInfo{}
	err = yaml.Unmarshal([]byte(content), info)
	if err != nil {
		return
	}

	sign := info.GenSign(key)
	if sign != info.Sign {
		err = errors.New("签名验证失败，请检查密钥或文件内容是否正确")
		return
	}

	password := info.GetPassword(key)
	lines := strings.Split(strings.TrimSpace(info.UserSettingText), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		bs, e := Decrypt(line, password)
		if e != nil {
			continue
		}
		if e = thrift.Deserialize(UserSetting, bs); e == nil {
			info.UserSetting = UserSetting.GetData()
		}
	}
	if info.UserSetting != nil && info.UserSetting["option"] != nil {
		option, ok := info.UserSetting["option"].(map[string]any)
		if ok {
			info.UserSettingSize = len(option)
		}
	}

	lines = strings.Split(strings.TrimSpace(info.ToolboxGroupListText), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		bs, e := Decrypt(line, password)
		if e != nil {
			continue
		}
		if e = thrift.Deserialize(ToolboxGroup, bs); e == nil {
			info.ToolboxGroupList = append(info.ToolboxGroupList, ToolboxGroup.GetData())
		}
	}
	info.ToolboxGroupSize = len(info.ToolboxGroupList)

	lines = strings.Split(strings.TrimSpace(info.ToolboxListText), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		bs, e := Decrypt(line, password)
		if e != nil {
			continue
		}
		if e = thrift.Deserialize(Toolbox, bs); e == nil {
			info.ToolboxList = append(info.ToolboxList, Toolbox.GetData())
		}
	}
	info.ToolboxSize = len(info.ToolboxList)

	lines = strings.Split(strings.TrimSpace(info.ToolboxExtendListText), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		bs, e := Decrypt(line, password)
		if e != nil {
			continue
		}
		if e = thrift.Deserialize(ToolboxExtend, bs); e == nil {
			info.ToolboxExtendList = append(info.ToolboxExtendList, ToolboxExtend.GetData())
		}
	}
	info.ToolboxExtendSize = len(info.ToolboxExtendList)
	return
}

var (
	lineS1 = "######################"
)

func Gen(key string, info *SyncInfo) (content string, err error) {
	var bs []byte
	var str string
	info.SignKey = strings.ToUpper(util.GetUUID())
	if key != "" {
		info.Encrypt = "是"
	} else {
		info.Encrypt = "否"
	}
	password := info.GetPassword(key)

	// 生成 用户 设置
	UserSetting.SetData(info.UserSetting)
	bs, err = thrift.Serialize(UserSetting)
	if err != nil {
		return
	}
	str, err = Encrypt(bs, password)
	if err != nil {
		return
	}
	info.UserSettingText += "  " + str + "\n"

	// 生成 工具组 设置
	for _, one := range info.ToolboxGroupList {
		ToolboxGroup.SetData(one)
		bs, err = thrift.Serialize(ToolboxGroup)
		if err != nil {
			return
		}
		str, err = Encrypt(bs, password)
		if err != nil {
			return
		}
		info.ToolboxGroupListText += "  " + str + "\n"
	}
	// 生成 工具 设置
	for _, one := range info.ToolboxList {
		Toolbox.SetData(one)
		bs, err = thrift.Serialize(Toolbox)
		if err != nil {
			return
		}
		str, err = Encrypt(bs, password)
		if err != nil {
			return
		}
		info.ToolboxListText += "  " + str + "\n"
	}
	// 生成 工具扩展 设置
	for _, one := range info.ToolboxExtendList {
		ToolboxExtend.SetData(one)
		bs, err = thrift.Serialize(ToolboxExtend)
		if err != nil {
			return
		}
		str, err = Encrypt(bs, password)
		if err != nil {
			return
		}
		info.ToolboxExtendListText += "  " + str + "\n"
	}

	info.Sign = info.GenSign(key)

	content += lineS1 + "\n"
	content += "说明: " + info.Explain + "\n"
	content += "所属: " + info.CreateBy + "\n"
	content += "加密: " + info.Encrypt + "\n"
	content += "时间: " + info.CreateAt + "\n"
	content += "密串: " + info.SignKey + "\n"
	content += "签名: " + info.Sign + "\n"
	content += lineS1 + "\n\n"

	content += lineS1 + "\n"
	content += "个人设置: |" + "\n"
	content += info.UserSettingText + "\n"
	content += lineS1 + "\n\n"

	content += lineS1 + "\n"
	content += "工具分组: |" + "\n"
	content += info.ToolboxGroupListText + "\n"
	content += lineS1 + "\n\n"

	content += lineS1 + "\n"
	content += "工具: |" + "\n"
	content += info.ToolboxListText + "\n"
	content += lineS1 + "\n\n"

	content += lineS1 + "\n"
	content += "工具扩展: |" + "\n"
	content += info.ToolboxExtendListText + "\n"
	content += lineS1 + "\n\n"

	return
}

func Encrypt(origData []byte, key []byte) (res string, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	bs, err := util.AesCBCEncrypt(origData, key)
	if err != nil {
		return
	}
	// 经过一次base64 否则 直接转字符串乱码
	res = base64.StdEncoding.EncodeToString(bs)
	return
}

func Decrypt(encrypt string, key []byte) (res []byte, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	// 经过一次base64 否则 直接转字符串乱码
	bs, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return
	}
	res, err = util.AesCBCDecrypt(bs, key)
	if err != nil {
		return
	}
	return
}

func GetSignData(obj any, includeEmpty bool) (data map[string]any) {
	data = make(map[string]any)
	objV := reflect.ValueOf(obj)
	for objV.Kind() == reflect.Ptr {
		objV = objV.Elem()
	}
	objT := reflect.TypeOf(obj)
	for objT.Kind() == reflect.Ptr {
		objT = objT.Elem()
	}

	for i := 0; i < objT.NumField(); i++ {
		field := objT.Field(i)
		fieldV := objV.Field(i)
		if fieldV.Interface() == nil {
			continue
		}

		fTKind := fieldV.Kind()
		if fTKind == reflect.Ptr && fieldV.IsNil() {
			continue
		}
		for fTKind == reflect.Ptr {
			fieldV = fieldV.Elem()
			fTKind = fieldV.Kind()
		}
		str := field.Tag.Get("sign")
		if str == "" {
			continue
		}
		if str != "-" {
			ss := strings.Split(str, ",")
			str = ss[0]
		} else {
			continue
		}
		v := fieldV.Interface()
		if includeEmpty && (v == "" || v == 0) {
			continue
		}
		data[str] = v
	}
	return
}

// JoinStringsInASCII 按照规则，参数名ASCII码从小到大排序后拼接
// data 待拼接的数据
// sep 连接符
// onlyValues 是否只包含参数值，true则不包含参数名，否则参数名和参数值均有
// includeEmpty 是否包含空值，true则包含空值，否则不包含，注意此参数不影响参数名的存在
// exceptKeys 被排除的参数名，不参与排序及拼接
func JoinStringsInASCII(data map[string]any, sep string, onlyValues, includeEmpty bool, exceptKeys ...string) string {
	var list []string
	var keyList []string
	m := make(map[string]int)
	if len(exceptKeys) > 0 {
		for _, except := range exceptKeys {
			m[except] = 1
		}
	}
	for k := range data {
		if _, ok := m[k]; ok {
			continue
		}
		value := util.GetStringValue(data[k])
		if !includeEmpty && (value == "" || value == "0") {
			continue
		}
		if onlyValues {
			keyList = append(keyList, k)
		} else {
			list = append(list, fmt.Sprintf("%s=%s", k, value))
		}
	}
	if onlyValues {
		asciiSort(keyList)
		for _, v := range keyList {
			list = append(list, util.GetStringValue(data[v]))
		}
	} else {
		asciiSort(list)
	}
	return strings.Join(list, sep)
}
func asciiSort(list []string) {
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	return
}
