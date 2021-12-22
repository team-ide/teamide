package modelcoder

import "time"

type FactoryScript interface {
	IsEmpty(obj interface{}) bool    // 是否为空  null或空字符串
	IsNotEmpty(obj interface{}) bool // 是否不为空  取 是否为空 反

	IsTrue(obj interface{}) bool  // 是否为真 入：true、"true"、1、"1"为真
	IsFalse(obj interface{}) bool // 是否为假  取 是否为真 反

	GetErrCode(err error) string // 获取异常中的错误码
	GetErrMsg(err error) string  // 获取异常中的信息

	AesECBEncrypt(data, key []byte) ([]byte, error) // AES ECB模式加密
	AesECBDecrypt(data, key []byte) ([]byte, error) // AES ECB模式解密

	AesCBCEncrypt(origData, key []byte) ([]byte, error) // AES CBC模式加密
	AesCBCDecrypt(origData, key []byte) ([]byte, error) // AES CBC模式解密

	Base64Encode(data string) (string, error)     // Base64 加密
	Base64Decode(origData string) (string, error) // Base64 解密

	MD5(data string) (string, error) // MD5 加密

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
}
