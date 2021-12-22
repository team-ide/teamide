package modelcoder

import "time"

type FactoryScriptDefault struct {
}

func (this_ *FactoryScriptDefault) IsEmpty(obj interface{}) (res bool) { // 是否为空  null或空字符串
	if obj == nil || obj == "" || obj == 0 {
		res = true
	}
	return
}

func (this_ *FactoryScriptDefault) IsNotEmpty(obj interface{}) (res bool) { // 是否不为空  取 是否为空 反
	return !this_.IsEmpty(obj)
}

func (this_ *FactoryScriptDefault) IsTrue(obj interface{}) (res bool) { // 是否为真 入：true、"true"、1、"1"为真
	if obj == true || obj == "true" || obj == 1 || obj == "1" {
		res = true
	}
	return
}

func (this_ *FactoryScriptDefault) IsFalse(obj interface{}) (res bool) { // 是否为假  取 是否为真 反
	return !this_.IsTrue(obj)
}

func (this_ *FactoryScriptDefault) GetErrCode(err error) (res string) { // 获取异常中的错误码
	res = "-1"
	return
}

func (this_ *FactoryScriptDefault) GetErrMsg(err error) (res string) { // 获取异常中的信息
	res = err.Error()
	return
}

func (this_ *FactoryScriptDefault) AesECBEncrypt(data, key []byte) (origData []byte, err error) { // AES ECB模式加密

	return
}

func (this_ *FactoryScriptDefault) AesECBDecrypt(origData, key []byte) (data []byte, err error) { // AES ECB模式解密

	return
}

func (this_ *FactoryScriptDefault) AesCBCEncrypt(data, key []byte) (origData []byte, err error) { // AES CBC模式加密

	return
}
func (this_ *FactoryScriptDefault) AesCBCDecrypt(origData, key []byte) (data []byte, err error) { // AES CBC模式解密

	return
}
func (this_ *FactoryScriptDefault) Base64Encode(data string) (origData string, err error) { // Base64 加密

	return
}
func (this_ *FactoryScriptDefault) Base64Decode(origData string) (data string, err error) { // Base64 解密

	return
}

func (this_ *FactoryScriptDefault) MD5(data string) (origData string, err error) { // MD5 加密

	return
}
func (this_ *FactoryScriptDefault) Now() (res time.Time) { // 当前时间

	return
}
func (this_ *FactoryScriptDefault) NowFormat(format string) (res string) { // 当前时间格式化 默认：yyyy-MM-dd hh:mm:ss.SSS

	return
}
func (this_ *FactoryScriptDefault) NowTime() (res int64) { // 当前时间戳

	return
}
func (this_ *FactoryScriptDefault) DateFormat(date time.Time, format string) (res string) { // 获取传入时间的时间格式化 默认：yyyy-MM-dd hh:mm:ss.SSS

	return
}
func (this_ *FactoryScriptDefault) DateTime(date time.Time) (res int64) { // 获取传入时间的时间戳

	return
}
func (this_ *FactoryScriptDefault) ToDate(time string, format string) (res time.Time) { // 根据格式化的日期和格式转为日期

	return
}
func (this_ *FactoryScriptDefault) ToDateByTime(time int64) (res time.Time) { // 根据时间戳转为日期

	return
}

func (this_ *FactoryScriptDefault) GetID() (res int64) { // 生成ID，生成不重复的ID

	return
}
func (this_ *FactoryScriptDefault) GetIDByType(idType string) (res int64) { // 根据ID类型生成ID，生成不重复的ID

	return
}

func (this_ *FactoryScriptDefault) GetStringID() (res string) { // 生成字符串ID，生成不重复的ID

	return
}
func (this_ *FactoryScriptDefault) GetStringIDByType(idType string) (res string) { // 根据ID类型生成字符串ID，生成不重复的ID

	return
}
