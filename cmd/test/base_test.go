package test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	//var err error
	urlStr := `thrift://192.168.6.81:11321/com.vrv.im.service.GroupMemberService$Iface?application=serverGroupMember&deprecated=false&dubbo=2.0.2&release=3.0.7&side=provider&thrift.protocol=vrv-binary&version=1.0&weight=1`
	fmt.Println("编码前:", urlStr)
	// 编码
	urlData := url.QueryEscape(urlStr)
	fmt.Println("编码后:", urlData)

	urlData = `thrift%3A%2F%2F192.168.6.81%3A11321%2Fcom.vrv.im.service.GroupMemberService%24Iface%3Fapplication%3DserverGroupMember%26deprecated%3Dfalse%26dubbo%3D2.0.2%26release%3D3.0.7%26side%3Dprovider%26thrift.protocol%3Dvrv-binary%26version%3D1.0%26weight%3D1`
	fmt.Println("解码前:", urlData)
	// 解码
	urlStr, err := url.QueryUnescape(urlData)
	if err != nil {
		panic(err)
	}
	fmt.Println("解码后:", urlStr)

	urlBean, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	bs, _ := json.Marshal(urlBean)
	fmt.Println("Bean：", string(bs))
	fmt.Println("协议：", urlBean.Scheme)
	fmt.Println("Host：", urlBean.Hostname())
	fmt.Println("Port：", urlBean.Port())

}

func TestIDMysql(t *testing.T) {
	//var err error
	str := "teamide:event:"
	fmt.Println(str)
	bs := []byte(str)
	fmt.Println(len(bs))

}
