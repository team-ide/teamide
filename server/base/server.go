package base

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	CerSecret string = "@Team DE Secret "
)

var (
	serverInfo *CerInfo
)

type CerInfo struct {
	IssuedBy    string `json:"颁发者"`
	IssuedTo    string `json:"颁发给"`
	No          string `json:"证书编号"`
	Id          string `json:"ID"`
	Key         string `json:"密钥"`
	IssueDate   string `json:"颁发日期"`
	ValidPeriod string `json:"有效期"`
	MAC         string `json:"MAC"`
	Version     string `json:"版本"`
	CerType     string `json:"证书类型"`
	User        string `json:"用户"`
	ServerId    int64
}

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := path + "/./conf/server.info"
	exists, err := PathExists(filePath)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("服务信息文件[" + filePath + "]不存在")
	}
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var line string
	var code string
	for {
		line, err = r.ReadString('\n')
		if err != nil && err != io.EOF {
			return
		}
		if code != "" {
			code += line
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
	serverInfo = GetCerInfoByCode(code)

	serverInfo.ServerId, err = strconv.ParseInt(serverInfo.No, 10, 64)
	if err != nil {
		fmt.Println("服务信息错误！")
		return
	}

	str := "测试加解密字段"
	str1 := Encrypt(str)
	if str1 == "" || str1 == str {
		panic("加密异常，请确认服务器信息是否正确！")
	}
	str2 := Decrypt(str1)
	if str2 == "" || str2 != str {
		panic("解密异常，请确认服务器信息是否正确！")
	}
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
func Encrypt(origData string) (res string) {
	if serverInfo == nil || serverInfo.Key == "" {
		fmt.Println("加密失败，证书信息不存在！")
		return
	}
	bs, err := AESEncrypt([]byte(origData), []byte(serverInfo.Key))
	if err != nil {
		fmt.Println("加密失败！")
		return
	}

	// 经过一次base64 否则 直接转字符串乱码
	res = base64.URLEncoding.EncodeToString(bs)
	return
}

//AES解密
func Decrypt(crypted string) (res string) {
	if serverInfo == nil || serverInfo.Key == "" {
		fmt.Println("解密失败，证书信息不存在！")
		return
	}
	// 经过一次base64 否则 直接转字符串乱码
	bs, err := base64.URLEncoding.DecodeString(crypted)
	if err != nil {
		fmt.Println("解密失败！")
		return
	}
	bs, err = AESDecrypt(bs, []byte(serverInfo.Key))
	if err != nil {
		fmt.Println("解密失败！")
		return
	}
	res = string(bs)
	return
}

func GetCerInfoByCode(code string) (info *CerInfo) {
	code = strings.ReplaceAll(code, "\n", "")
	code = strings.TrimSpace(code)
	strs := strings.Split(code, " ")
	bs := []byte{}
	for _, str := range strs {
		n, _ := strconv.Atoi(str)
		b := byte(n)
		bs = append(bs, b)
	}

	bs, _ = AESDecrypt(bs, []byte(CerSecret))

	info = &CerInfo{}
	json.Unmarshal(bs, &info)
	return
}

func buildOrderStr(infoMap map[string]string, secret string) (returnStr string) {
	keys := []string{}

	for k := range infoMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		if infoMap[k] == "" {
			continue
		}
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}

		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(infoMap[k])
	}

	buf.WriteString("&secret=" + secret)
	returnStr = buf.String()

	return returnStr
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES加密,CBC
func AESEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pkcs7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return crypted, nil
}

//AES解密
func AESDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs7UnPadding(origData)
	return origData, nil
}
