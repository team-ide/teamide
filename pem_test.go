package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestPem(t *testing.T) {

	//rsa 密钥文件产生
	// GenRsaKey(1024)
	PrintPem()
}

func PrintPem() {
	//读取证书并解码
	pemTmp, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	certBlock, restBlock := pem.Decode(pemTmp)
	if certBlock == nil {
		fmt.Println(err)
		return
	}
	//可从剩余判断是否有证书链等，继续解析
	fmt.Println(restBlock)
	//证书解析
	certBody, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 提取公钥
	publicKeyDer, _ := x509.MarshalPKIXPublicKey(certBody.PublicKey)
	publicKeyBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDer,
	}
	publicKeyPem := string(pem.EncodeToMemory(&publicKeyBlock))
	fmt.Println(publicKeyPem)
	//可以根据证书结构解析
	fmt.Println(certBody.SignatureAlgorithm)
	fmt.Println(certBody.PublicKeyAlgorithm)
}

//RSA公钥私钥产生
func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
