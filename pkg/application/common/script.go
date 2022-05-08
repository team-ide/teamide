package common

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"regexp"
	"strings"
	"teamide/pkg/application/base"
	"teamide/pkg/application/model"
	"teamide/pkg/util"
	"time"
)

type IScript interface {
	TrimSuffix(s, suffix string) (res string)
	TrimPrefix(s, prefix string) (res string)
	GetLock(key string) (res Locker, err error)

	Length(obj interface{}) int                       //
	NotMatch(pattern string, value interface{}) bool  //
	Match(pattern string, value interface{}) bool     //
	ThrowError(code string, msg string) (bool, error) //

	IsEmpty(obj interface{}) bool    // 是否为空  null或空字符串
	IsNotEmpty(obj interface{}) bool // 是否不为空  取 是否为空 反

	IsTrue(obj interface{}) bool  // 是否为真 入：true、"true"、1、"1"为真
	IsFalse(obj interface{}) bool // 是否为假  取 是否为真 反

	GetErrCode(err error) string // 获取异常中的错误码
	GetErrMsg(err error) string  // 获取异常中的信息

	WrapTableName(database string, table string) string // 包装表名
	WrapColumnName(alias string, column string) string  // 包装列名

	DataToJSON(data interface{}) string
	JSONToData(json string) (interface{}, error)

	AesECBEncrypt(data, key string) (string, error) // AES ECB模式加密
	AesECBDecrypt(data, key string) (string, error) // AES ECB模式解密

	AesCBCEncrypt(origData, key string) (string, error) // AES CBC模式加密
	AesCBCDecrypt(origData, key string) (string, error) // AES CBC模式解密

	Base64Encode(data string) (string, error)     // Base64 加密
	Base64Decode(origData string) (string, error) // Base64 解密

	MD5(data string) (string, error) // MD5 加密

	UUID() string // MD5 加密
	GetReader(obj interface{}) (res io.Reader, err error)
	GetFileTypeSuffix(obj interface{}) (res string)

	Now() time.Time                                  // 当前时间
	NowFormat(format string) string                  // 当前时间格式化 默认：yyyy-MM-dd hh:mm:ss.SSS
	NowTime() int64                                  // 当前时间戳
	DateFormat(date time.Time, format string) string // 获取传入时间的时间格式化 默认：yyyy-MM-dd hh:mm:ss.SSS
	DateTime(ate time.Time) int64                    // 获取传入时间的时间戳
	ToDate(time string, format string) time.Time     // 根据格式化的日期和格式转为日期
	ToDateByTime(time int64) time.Time               // 根据时间戳转为日期

	GetID() int64                    // 生成ID，生成不重复的ID
	GetIDByType(idType string) int64 // 根据ID类型生成ID，生成不重复的ID

	GetStringID() string                    // 生成字符串ID，生成不重复的ID
	GetStringIDByType(idType string) string // 根据ID类型生成字符串ID，生成不重复的ID

	RandomInt(min int, max int) (int, error)               // 生成随机字符串
	RandomString(minLen int, maxLen int) (string, error)   // 生成随机字符串
	RandomUserName(minLen int, maxLen int) (string, error) // 生成随机用户名
	ToPinYin(name string) (string, error)                  // 转换为拼音
}

type ScriptDefault struct {
}

func (this_ *ScriptDefault) Length(obj interface{}) (res int) { //
	if obj == nil {
		return 0
	}
	switch v := obj.(type) {
	case string:
		return len(v)
	case map[interface{}]interface{}:
		return len(v)
	case []interface{}:
		return len(v)
	default:
		return len(fmt.Sprint(v))
	}
}
func (this_ *ScriptDefault) TrimPrefix(s, prefix string) (res string) { //
	res = strings.TrimPrefix(s, prefix)
	return
}

func (this_ *ScriptDefault) TrimSuffix(s, suffix string) (res string) { //
	res = strings.TrimSuffix(s, suffix)
	return
}

//
func (this_ *ScriptDefault) GetLock(key string) (res Locker, err error) { //

	return
}

//
func (this_ *ScriptDefault) ThrowError(code string, msg string) (res bool, err error) { //
	err = base.NewError(code, msg)
	return
}
func (this_ *ScriptDefault) Match(pattern string, value interface{}) (res bool) { //
	var re *regexp.Regexp = regexp.MustCompile(pattern)
	res = re.MatchString(fmt.Sprint(value))
	return
}
func (this_ *ScriptDefault) NotMatch(pattern string, value interface{}) (res bool) { //
	res = !this_.Match(pattern, value)
	return
}

func (this_ *ScriptDefault) IsEmpty(obj interface{}) (res bool) { // 是否为空  null或空字符串
	if obj == nil || obj == "" || obj == 0 {
		res = true
	}
	return
}

func (this_ *ScriptDefault) IsNotEmpty(obj interface{}) (res bool) { // 是否不为空  取 是否为空 反
	return !this_.IsEmpty(obj)
}

func (this_ *ScriptDefault) IsTrue(obj interface{}) (res bool) { // 是否为真 入：true、"true"、1、"1"为真
	if obj == true || obj == "true" || obj == 1 || obj == "1" {
		res = true
	}
	return
}

func (this_ *ScriptDefault) IsFalse(obj interface{}) (res bool) { // 是否为假  取 是否为真 反
	return !this_.IsTrue(obj)
}

func (this_ *ScriptDefault) GetErrCode(err error) (res string) { // 获取异常中的错误码
	res = "-1"
	return
}

func (this_ *ScriptDefault) GetErrMsg(err error) (res string) { // 获取异常中的信息
	res = err.Error()
	return
}

func (this_ *ScriptDefault) WrapTableName(database string, table string) (res string) { // 包装表名
	if this_.IsNotEmpty(database) {
		res = "" + database + "."
	}
	res += "" + table + ""
	return
}
func (this_ *ScriptDefault) WrapColumnName(tableAlias string, column string) (res string) { // 包装列名
	if this_.IsNotEmpty(tableAlias) {
		res = "" + tableAlias + "."
	}
	res += "" + column + ""
	return
}

func (this_ *ScriptDefault) DataToJSON(data interface{}) (res string) {
	if data == nil {
		return
	}
	res = base.ToJSON(data)
	return
}
func (this_ *ScriptDefault) JSONToData(json string) (data interface{}, err error) {
	if json == "" {
		return
	}
	data = make(map[string]interface{})
	err = base.ToBean([]byte(json), &data)
	return
}

func (this_ *ScriptDefault) AesECBEncrypt(data, key string) (origData string, err error) { // AES ECB模式加密
	var bs []byte
	bs, err = base.AesECBEncrypt([]byte(data), []byte(key))
	if err != nil {
		return
	}
	// 经过一次base64 否则 直接转字符串乱码
	origData = base64.StdEncoding.EncodeToString(bs)
	return
}

func (this_ *ScriptDefault) AesECBDecrypt(origData, key string) (data string, err error) { // AES ECB模式解密
	var bs []byte
	// 经过一次base64 否则 直接转字符串乱码
	bs, err = base64.StdEncoding.DecodeString(origData)
	if err != nil {
		return
	}
	bs, err = base.AesECBDecrypt(bs, []byte(key))
	if err != nil {
		return
	}
	data = string(bs)
	return
}

func (this_ *ScriptDefault) AesCBCEncrypt(data, key string) (origData string, err error) { // AES CBC模式加密
	var bs []byte
	bs, err = base.AesCBCEncrypt([]byte(data), []byte(key))
	if err != nil {
		return
	}
	// 经过一次base64 否则 直接转字符串乱码
	origData = base64.StdEncoding.EncodeToString(bs)
	return
}
func (this_ *ScriptDefault) AesCBCDecrypt(origData, key string) (data string, err error) { // AES CBC模式解密
	var bs []byte
	// 经过一次base64 否则 直接转字符串乱码
	bs, err = base64.StdEncoding.DecodeString(origData)
	if err != nil {
		return
	}
	bs, err = base.AesCBCDecrypt(bs, []byte(key))
	if err != nil {
		return
	}
	data = string(bs)
	return
}
func (this_ *ScriptDefault) Base64Encode(data string) (origData string, err error) { // Base64 加密

	return
}
func (this_ *ScriptDefault) Base64Decode(origData string) (data string, err error) { // Base64 解密

	return
}

func (this_ *ScriptDefault) MD5(data string) (origData string, err error) { // MD5 加密
	m := md5.New()
	_, err = io.WriteString(m, data)
	if err != nil {
		return
	}
	arr := m.Sum(nil)
	origData = fmt.Sprintf("%x", arr)
	return
}

func (this_ *ScriptDefault) GetReader(obj interface{}) (res io.Reader, err error) { // MD5 加密
	if obj == nil {
		return
	}
	switch v := obj.(type) {
	case io.Reader:
		res = v
	case *model.FileInfo:
		res, err = v.GetReader()
	case map[string]interface{}:
		fileInfo := &model.FileInfo{}
		if v["name"] != nil {
			fileInfo.Name = v["name"].(string)
		}
		if v["path"] != nil {
			fileInfo.Path = v["path"].(string)
		}
		if v["absolutePath"] != nil {
			fileInfo.AbsolutePath = v["absolutePath"].(string)
		}
		if v["fileHeader"] != nil {
			fileInfo.FileHeader = v["fileHeader"].(*multipart.FileHeader)
		}
		if v["file"] != nil {
			fileInfo.File = v["file"].(*os.File)
		}
		res, err = fileInfo.GetReader()
	}
	return
}

func (this_ *ScriptDefault) GetFileTypeSuffix(obj interface{}) (res string) { // MD5 加密
	if obj == nil {
		return
	}

	var fileName string
	var fileType string
	switch v := obj.(type) {
	case *model.FileInfo:
		fileName = v.Name
		fileType = v.Type
	case map[string]interface{}:
		if v["name"] != nil {
			fileName = v["name"].(string)
		}
		if v["type"] != nil {
			fileType = v["type"].(string)
		}
	}

	if fileType == "" && strings.Contains(fileName, ".") {
		fileType = fileName[strings.LastIndex(fileName, ".")+1:]
	}
	if fileType != "" {
		res = "." + fileType
	}
	return
}

func (this_ *ScriptDefault) UUID() (res string) {
	return util.GenerateUUID()
}

func (this_ *ScriptDefault) Now() (res time.Time) { // 当前时间
	res = time.Now()
	return
}
func (this_ *ScriptDefault) NowFormat(format string) (res string) { // 当前时间格式化 默认：yyyy-MM-dd hh:mm:ss.SSS
	if format == "" {
		format = "2006-01-02 15:04:05.000"
	}
	res = this_.Now().Format(format)
	return
}
func (this_ *ScriptDefault) NowTime() (res int64) { // 当前时间戳
	res = time.Now().UnixNano() / 1e6
	return
}
func (this_ *ScriptDefault) DateFormat(date time.Time, format string) (res string) { // 获取传入时间的时间格式化 默认：yyyy-MM-dd hh:mm:ss.SSS

	return
}
func (this_ *ScriptDefault) DateTime(date time.Time) (res int64) { // 获取传入时间的时间戳

	return
}
func (this_ *ScriptDefault) ToDate(time string, format string) (res time.Time) { // 根据格式化的日期和格式转为日期

	return
}
func (this_ *ScriptDefault) ToDateByTime(time int64) (res time.Time) { // 根据时间戳转为日期

	return
}

func (this_ *ScriptDefault) GetID() (res int64) { // 生成ID，生成不重复的ID
	num, _ := this_.RandomInt(1, 999999999)
	res = int64(num)
	return
}
func (this_ *ScriptDefault) GetIDByType(idType string) (res int64) { // 根据ID类型生成ID，生成不重复的ID

	return
}

func (this_ *ScriptDefault) GetStringID() (res string) { // 生成字符串ID，生成不重复的ID

	return
}
func (this_ *ScriptDefault) GetStringIDByType(idType string) (res string) { // 根据ID类型生成字符串ID，生成不重复的ID

	return
}

func (this_ *ScriptDefault) RandomInt(minLen int, maxLen int) (int, error) { // 生成随机字符串
	return util.RandomInt(minLen, maxLen)
}

func (this_ *ScriptDefault) RandomString(minLen int, maxLen int) (string, error) { // 生成随机字符串
	return util.RandomString(minLen, maxLen)
}

func (this_ *ScriptDefault) RandomUserName(minLen int, maxLen int) (string, error) { // 生成随机用户名
	return util.RandomUserName(minLen, maxLen)
}

func (this_ *ScriptDefault) ToPinYin(name string) (string, error) { // 转换为拼音
	return util.ToPinYin(name)
}
