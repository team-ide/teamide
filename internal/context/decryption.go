package context

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"strings"
)

func NewDefaultDecryption(logger *zap.Logger) (res *Decryption, err error) {
	res = &Decryption{
		Logger: logger,
	}

	res.rsaPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC28eh+rmkHWEfFSjRIpQzdBaCr
t32/EYH39D4DuCEwQgRYBEynPxr6yhS0ozdEujSX0vENk2YQ6RdnOfkVCLzh/huN
6aguW94DFmU5Xc0AdtglSekCDE8Alk4MmhH6p+nN2Z22FiSIZY63rw0026613rD6
y0QLQ1GgtBeAVaNdhwIDAQAB
-----END PUBLIC KEY-----`

	res.rsaPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
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

	err = res.init()
	if err != nil {
		return
	}

	return
}

func NewDecryption(publicPath string, privatePath string, logger *zap.Logger) (res *Decryption, err error) {
	res = &Decryption{
		Logger: logger,
	}
	exists, err := util.PathExists(publicPath)
	if err != nil {
		panic(err)
	}
	if !exists {
		err = errors.New(fmt.Sprint("public key[", publicPath, "] not defined"))
		return
	}
	bs, err := os.ReadFile(publicPath)
	if err != nil {
		return
	}
	res.rsaPublicKey = strings.TrimSpace(string(bs))

	exists, err = util.PathExists(privatePath)
	if err != nil {
		panic(err)
	}
	if !exists {
		err = errors.New(fmt.Sprint("private key[", privatePath, "] not defined"))
		return
	}
	bs, err = os.ReadFile(privatePath)
	if err != nil {
		return
	}
	res.rsaPrivateKey = strings.TrimSpace(string(bs))

	err = res.init()
	if err != nil {
		return
	}
	return
}

type Decryption struct {
	rsaPublicKey  string
	rsaPrivateKey string
	Logger        *zap.Logger
}

// init 初始化加解密
func (this_ *Decryption) init() (err error) {

	str := "测试加解密字段"
	str1, err := this_.Encrypt(str)
	if err != nil {
		this_.Logger.Error("加密异常!", zap.Error(err))
		return
	}
	if str1 == "" || str1 == str {
		err = errors.New("加密异常，请确认服务器信息是否正确")
		this_.Logger.Error("加密异常!", zap.Error(err))
		return
	}
	this_.Logger.Info("服务器加密成功!")
	str2, err := this_.Decrypt(str1)
	if err != nil {
		this_.Logger.Error("解密异常!", zap.Error(err))
		return
	}
	if str2 != str {
		err = errors.New("解密异常，请确认服务器信息是否正确")
		this_.Logger.Error("解密异常!", zap.Error(err))
		return
	}
	this_.Logger.Info("服务器解密成功!")

	str1, err = this_.Encrypt(str)
	if err != nil {
		this_.Logger.Error("加密异常!", zap.Error(err))
		return
	}
	if str1 == "" || str1 == str {
		err = errors.New("加密异常，请确认服务器信息是否正确")
		this_.Logger.Error("加密异常!", zap.Error(err))
		return
	}
	this_.Logger.Info("服务器加密验证成功!")
	str2, err = this_.Decrypt(str1)
	if err != nil {
		this_.Logger.Error("解密异常!", zap.Error(err))
		return
	}
	if str2 != str {
		err = errors.New("解密异常，请确认服务器信息是否正确")
		this_.Logger.Error("解密异常!", zap.Error(err))
		return
	}
	this_.Logger.Info("服务器解密验证成功!")
	return
}

// Encrypt 加密
func (this_ *Decryption) Encrypt(data string) (ciphertext string, err error) {
	ciphertext, err = util.RsaEncryptByKey(data, this_.rsaPublicKey)
	if err != nil {
		this_.Logger.Error("加密失败", zap.Any("dataLength", len(data)), zap.Error(err))
		return
	}
	return
}

// Decrypt 解密
func (this_ *Decryption) Decrypt(ciphertext string) (data string, err error) {
	data, err = util.RsaDecryptByKey(ciphertext, this_.rsaPrivateKey)
	if err != nil {
		this_.Logger.Error("解密失败", zap.Any("ciphertextLength", len(ciphertext)), zap.Error(err))
		return
	}
	return
}

// IsEncrypt 解密
func (this_ *Decryption) IsEncrypt(ciphertext string) (res bool) {
	_, err := util.RsaDecryptByKey(ciphertext, this_.rsaPrivateKey)
	if err != nil {
		return false
	}
	return true
}
