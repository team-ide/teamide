package component

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"server/base"
	"strconv"
	"strings"
)

const (
	CerSecret string = "@Team DE Secret "
)

var (
	serverInfo = initServer()
)

type CerInfo struct {
	IssuedBy    string `json:"颁发者"`
	IssuedTo    string `json:"颁发给"`
	No          string `json:"证书编号"`
	Key         string `json:"密钥"`
	IssueDate   string `json:"颁发日期"`
	ValidPeriod string `json:"有效期"`
	MAC         string `json:"MAC"`
	Version     string `json:"版本"`
	CerType     string `json:"证书类型"`
	User        string `json:"用户"`
	ServerId    int64
}

func initServer() *CerInfo {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := path + "/./conf/server.info"
	exists, err := base.PathExists(filePath)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("服务信息文件[" + filePath + "]不存在")
	}
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var line string
	var code string
	for {
		line, err = r.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		if strings.IndexByte(line, '\n') > 0 {
			line = line[0:strings.IndexByte(line, '\n')]
		}
		if strings.IndexByte(line, '\r') > 0 {
			line = line[0:strings.IndexByte(line, '\r')]
		}
		if code != "" {
			code = fmt.Sprint(code, line)
		} else {
			if len(line) > 0 {
				if IsNum(line[0:1]) {
					code = line
				}
			}
		}
		if err == io.EOF {
			err = nil
			break
		}
	}
	info := GetCerInfoByCode(code)

	info.ServerId, err = strconv.ParseInt(info.No, 10, 64)
	if err != nil {
		panic("服务信息错误!")
	}
	Logger.Info("服务器信息加载成功!")

	return info

}

func init() {
	str := "测试加解密字段"
	str1 := AesEncryptCBC(str)
	if str1 == "" || str1 == str {
		panic("加密异常，请确认服务器信息是否正确!")
	}
	Logger.Info("服务器加密验证成功!")
	str2 := AesDecryptCBC(str1)
	if str2 == "" || str2 != str {
		panic("解密异常，请确认服务器信息是否正确!")
	}

	str1 = AesEncryptECB(str)
	if str1 == "" || str1 == str {
		panic("加密异常，请确认服务器信息是否正确!")
	}
	Logger.Info("服务器加密验证成功!")
	str2 = AesDecryptECB(str1)
	if str2 == "" || str2 != str {
		panic("解密异常，请确认服务器信息是否正确!")
	}
	Logger.Info("服务器解密验证成功!")
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func GetServerId() (serverId int64) {
	if serverInfo == nil {
		return
	}
	serverId = serverInfo.ServerId
	return
}
func GetBaseID() (id int64) {
	if serverInfo == nil {
		return
	}
	id = 164317281
	return
}

//AES加密,CBC
func AesEncryptCBC(origData string) (res string) {
	if serverInfo == nil || serverInfo.Key == "" {
		Logger.Error("加密失败，证书信息不存在!")
		return
	}
	return AesEncryptCBCByKey(origData, serverInfo.Key)
}

//AES解密,CBC
func AesDecryptCBC(crypted string) (res string) {
	if serverInfo == nil || serverInfo.Key == "" {
		Logger.Error("解密失败，证书信息不存在!")
		return
	}
	return AesDecryptCBCByKey(crypted, serverInfo.Key)
}

//AES加密,ECB
func AesEncryptECB(origData string) (res string) {
	if serverInfo == nil || serverInfo.Key == "" {
		Logger.Error("加密失败，证书信息不存在!")
		return
	}
	return AesEncryptECBByKey(origData, serverInfo.Key)
}

//AES解密,ECB
func AesDecryptECB(crypted string) (res string) {
	if serverInfo == nil || serverInfo.Key == "" {
		Logger.Error("解密失败，证书信息不存在!")
		return
	}
	return AesDecryptECBByKey(crypted, serverInfo.Key)
}

//AES加密,CBC
func AesEncryptCBCByKey(origData string, key string) (res string) {
	bs, err := base.AesCBCEncrypt([]byte(origData), []byte(key))
	if err != nil {
		Logger.Error("加密失败!")
		return
	}
	// 经过一次base64 否则 直接转字符串乱码
	res = base64.StdEncoding.EncodeToString(bs)
	return
}

//AES解密,CBC
func AesDecryptCBCByKey(crypted string, key string) (res string) {
	// 经过一次base64 否则 直接转字符串乱码
	bs, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		Logger.Error("解密失败!")
		return
	}
	bs, err = base.AesCBCDecrypt(bs, []byte(key))
	if err != nil {
		Logger.Error("解密失败!")
		return
	}
	res = string(bs)
	return
}

//AES加密,ECB
func AesEncryptECBByKey(origData string, key string) (res string) {
	bs, err := base.AesECBEncrypt([]byte(origData), []byte(key))
	if err != nil {
		Logger.Error("加密失败!")
		return
	}
	// 经过一次base64 否则 直接转字符串乱码
	res = base64.StdEncoding.EncodeToString(bs)
	return
}

//AES解密,ECB
func AesDecryptECBByKey(crypted string, key string) (res string) {
	// 经过一次base64 否则 直接转字符串乱码
	bs, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		Logger.Error("解密失败!")
		return
	}
	bs, err = base.AesECBDecrypt(bs, []byte(key))
	if err != nil {
		Logger.Error("解密失败!")
		return
	}
	res = string(bs)
	return
}

func GetCerInfoByCode(code string) (info *CerInfo) {
	code = strings.ReplaceAll(code, "\n", "")
	code = strings.ReplaceAll(code, "\r", "")
	code = strings.TrimSpace(code)
	strs := strings.Split(code, " ")
	bs := []byte{}
	for _, str := range strs {
		n, _ := strconv.Atoi(str)
		b := byte(n)
		bs = append(bs, b)
	}

	bs, _ = base.AesCBCDecrypt(bs, []byte(CerSecret))
	info = &CerInfo{}
	json.Unmarshal(bs, info)
	return
}
