package component

import (
	"errors"
	"fmt"
	"os"
	base2 "teamide/internal/server/base"
)

var (
	rsa_public_key  = ""
	rsa_private_key = ""
)

func initServer() {
	if base2.IS_STAND_ALONE {
		rsa_public_key = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC28eh+rmkHWEfFSjRIpQzdBaCr
t32/EYH39D4DuCEwQgRYBEynPxr6yhS0ozdEujSX0vENk2YQ6RdnOfkVCLzh/huN
6aguW94DFmU5Xc0AdtglSekCDE8Alk4MmhH6p+nN2Z22FiSIZY63rw0026613rD6
y0QLQ1GgtBeAVaNdhwIDAQAB
-----END PUBLIC KEY-----`

		rsa_private_key = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQC28eh+rmkHWEfFSjRIpQzdBaCrt32/EYH39D4DuCEwQgRYBEyn
Pxr6yhS0ozdEujSX0vENk2YQ6RdnOfkVCLzh/huN6aguW94DFmU5Xc0AdtglSekC
DE8Alk4MmhH6p+nN2Z22FiSIZY63rw0026613rD6y0QLQ1GgtBeAVaNdhwIDAQAB
AoGAFp0edRJQD0VdUcjTX6tvRJ7edntvAsBCRYkeZU1MZO+0I8EcTIwjZJ64IoAO
Y+N0ftPnUhtHQY3eg7cJ0AzNdBgcGZqmiYO5ky6BPoE10VsCWHMi5CzciXVIkrBH
2/mlRDCGSeuXmyoXdrnvoAUTrZLNWWiGtyIdTkLR85SRYFkCQQDAkK+9pl/u6aQR
yEDvwlQTHPi7ODsjLxjzzi5ieGjfLdS/RlxrCCOnQEdUTwn8u7qSPfVkpOmxV7tE
UBKC012dAkEA8zXu7adb5Sl+1nZScmR3NeMytFEMMBwSLOk6De6xz8YPGpQzJocm
RxZiZ+C1c2udJNsLiRAVgQYtW0t6IaiQcwJAckotlCEYFSOklk1FhUfQQJvUYMIK
D2LXq3R3AUi37aY0++WV2oy1JII5E6fppJADNuMBL1/Vt8T7R5tCsVUj3QJAGc9w
csIfC3vS3RmjeEZXLF3XJLGxNG3WM/PwWEgrkJw5QB3YK8+N7V9fxBxhxUT3YVDp
sXsGfTHVoGmrJWVJJwJANJVpuXMFMdmoUQqfoCZrBUlcsZmldi3E3AvwT6WCmryC
Mr8hw5UEgHMbPnSA8H96ppLBTMOh9sgNp3ryDFE6Mw==
-----END RSA PRIVATE KEY-----`
		return
	}
	filePath := base2.BaseDir + "conf/public.pem"
	exists, err := base2.PathExists(filePath)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic(errors.New(fmt.Sprint("public key[", filePath, "] not defind")))
	}
	bs, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	rsa_public_key = string(bs)

	filePath = base2.BaseDir + "conf/private.pem"
	exists, err = base2.PathExists(filePath)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic(errors.New(fmt.Sprint("private key[", filePath, "] not defind")))
	}
	bs, err = os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	rsa_private_key = string(bs)

}

func init() {
	initServer()
	str := "测试加解密字段"
	str1 := RSAEncrypt(str)
	if str1 == "" || str1 == str {
		Logger.Error(LogStr("加密异常，请确认服务器信息是否正确!"))
		panic("加密异常，请确认服务器信息是否正确!")
	}
	Logger.Info("服务器加密成功!")
	str2 := RSADecrypt(str1)
	if str2 == "" || str2 != str {
		Logger.Error(LogStr("解密异常，请确认服务器信息是否正确!"))
		panic("解密异常，请确认服务器信息是否正确!")
	}
	Logger.Info("服务器解密成功!")

	str1 = RSAEncrypt(str)
	if str1 == "" || str1 == str {
		Logger.Error(LogStr("加密异常，请确认服务器信息是否正确!"))
		panic("加密异常，请确认服务器信息是否正确!")
	}
	Logger.Info("服务器加密验证成功!")
	str2 = RSADecrypt(str1)
	if str2 == "" || str2 != str {
		Logger.Error(LogStr("解密异常，请确认服务器信息是否正确!"))
		panic("解密异常，请确认服务器信息是否正确!")
	}
	Logger.Info("服务器解密验证成功!")
}

func GetBaseID() int64 {
	return 1000
}

func RSAEncrypt(origData string) (res string) {
	if rsa_public_key == "" {
		Logger.Error("加密失败，加密密钥不存在!")
		return
	}
	res, err := base2.RSAEncryptByKey(origData, rsa_public_key)
	if err != nil {
		Logger.Error(LogStr("加密失败:", err))
		return
	}
	return
}

func RSADecrypt(crypted string) (res string) {
	if rsa_private_key == "" {
		Logger.Error("解密失败，解密密钥不存在!")
		return
	}
	res, err := base2.RSADecryptByKey(crypted, rsa_private_key)
	if err != nil {
		Logger.Error(LogStr("解密失败:", err))
		return
	}
	return
}
