package component

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"teamide/server/base"
	"teamide/util"
)

const (
	CerSecret string = "@Team DE Secret "
)

var (
	serverInfo = initServer()
)
var (
	defaultServerConfig = `
Team IDE 证书
证书编号        567910
颁发者          Team IDE
颁发给          Team IDE 免费服务器
颁发日期        2021/10/25
有效期          2021/10/25 - 永久
MAC             ALL
版本            0.0.1
证书类型        免费版
密钥            JUbLWZUhy0ZcYJfU
用户            100
签名            8c128f24d5f0aba4d92935929b753a96
103 190 104 203 150 145 192 138
 132 133 177 163 123 76 161 25
2 101 251 233 217 44 60 254 18
5 15 250 39 182 43 205 21 16 2
6 108 22 6 94 169 179 249 224 
109 102 100 227 88 49 95 70 27
 89 40 96 23 155 58 129 194 89
 114 81 18 101 209 93 183 43 1
12 161 240 149 117 19 23 200 4
6 52 131 167 128 225 187 227 1
90 201 45 101 205 45 62 79 177
 19 78 166 149 195 17 65 43 13
9 73 238 228 200 66 106 10 35 
113 122 100 47 166 175 160 10 
85 247 81 137 147 191 124 180 
172 21 142 16 98 43 194 205 19
9 201 210 22 131 137 90 69 161
 139 29 39 65 3 17 116 228 187
 233 197 178 237 206 29 154 97
 198 73 134 13 8 186 242 231 1
61 60 7 88 202 84 106 108 212 
88 141 84 57 36 59 238 141 217
 73 250 237 27 20 227 214 175 
140 241 238 244 243 213 32 21 
72 166 234 244 16 196 68 105 1
58 123 64 238 193 159 238 198 
34 80 13 66 146 246 247 139 17
9 57 214 171 220 171 198 175 2
51 19 18 166 242 23 17 2 135 5
9 127 196 67 133 60 200 48 75 
8 197 115 115 158 235 0 182 66
 184 233 209 65 156 127 77 57 
184 131 7 125 74 50 181 72 125
 74 198 99 179 168 206 39 148 
119 225 219 194 178 184 210 11
9 129 226 162 9 39 190 201 35 
110 246 79`
)

type CerInfo struct {
	IssuedBy    string `json:"颁发者,omitempty"`
	IssuedTo    string `json:"颁发给,omitempty"`
	No          string `json:"证书编号,omitempty"`
	Key         string `json:"密钥,omitempty"`
	IssueDate   string `json:"颁发日期,omitempty"`
	ValidPeriod string `json:"有效期,omitempty"`
	MAC         string `json:"MAC,omitempty"`
	Version     string `json:"版本,omitempty"`
	CerType     string `json:"证书类型,omitempty"`
	User        string `json:"用户,omitempty"`
	ServerId    int64
}

func initServer() *CerInfo {
	filePath := util.BaseDir + "conf/server.info"
	exists, err := base.PathExists(filePath)
	if err != nil {
		panic(err)
	}
	serverConfig := defaultServerConfig
	if exists {

		var f *os.File
		f, err = os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		var bs []byte
		bs, err = io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		if len(bs) > 0 {
			serverConfig = string(bs)
		}
	}
	serverConfigLines := strings.Split(serverConfig, "\n")

	var code string
	for _, line := range serverConfigLines {
		line = strings.TrimPrefix(line, "\r")
		line = strings.TrimPrefix(line, "\t")
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
		Logger.Error(LogStr("服务信息错误!"))
		panic("服务信息错误!")
	}
	Logger.Info("服务器信息加载成功!")

	return info

}

func init() {
	str := "测试加解密字段"
	str1 := AesEncryptCBC(str)
	if str1 == "" || str1 == str {
		Logger.Error(LogStr("加密异常，请确认服务器信息是否正确!"))
		panic("加密异常，请确认服务器信息是否正确!")
	}
	Logger.Info("服务器加密成功!")
	str2 := AesDecryptCBC(str1)
	if str2 == "" || str2 != str {
		Logger.Error(LogStr("解密异常，请确认服务器信息是否正确!"))
		panic("解密异常，请确认服务器信息是否正确!")
	}
	Logger.Info("服务器解密成功!")

	str1 = AesEncryptECB(str)
	if str1 == "" || str1 == str {
		Logger.Error(LogStr("加密异常，请确认服务器信息是否正确!"))
		panic("加密异常，请确认服务器信息是否正确!")
	}
	Logger.Info("服务器加密验证成功!")
	str2 = AesDecryptECB(str1)
	if str2 == "" || str2 != str {
		Logger.Error(LogStr("解密异常，请确认服务器信息是否正确!"))
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
func GetKey() (key string) {
	if serverInfo == nil {
		return
	}
	key = serverInfo.Key
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
