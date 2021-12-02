package config

import "C"
import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"unsafe"
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
}

var (
	ServerInfo       *CerInfo
	ServerInfoConent string
)

func loadServerInfo() {
	var err error
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := path + "/./conf/server.info"
	ServerInfo, ServerInfoConent, err = ReadCerInfo(filePath)
	if err != nil {
		panic(err)
	}
	libPath := path + "/./lib/server"
	lib, err := syscall.LoadDLL(libPath)
	if err != nil {
		panic(err)
	}
	validate, err := lib.FindProc("CValidate")
	if err != nil {
		panic(err)
	}
	res1, res2, err := validate.Call(uintptr(unsafe.Pointer(C.CString(ToJSON(ServerInfo)))), uintptr(unsafe.Pointer(C.CString(ServerInfoConent))))
	if err != nil {
		panic(err)
	}
	fmt.Println(res1)
	fmt.Println(res2)
}

func ReadCerInfo(cerInfoPath string) (info *CerInfo, conent string, err error) {
	var exists bool
	exists, err = PathExists(cerInfoPath)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New("证书信息[" + cerInfoPath + "]不存在！")
		return
	}
	var f *os.File
	f, err = os.Open(cerInfoPath)
	if err != nil {
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var line string
	var infoMap = map[string]string{}
	for {

		line, err = r.ReadString('\n')
		if err != nil && err != io.EOF {
			return
		}
		if strings.Index(line, "  ") > 0 {
			key := line[:strings.Index(line, "  ")]
			value := getInfoValue(line, key)
			infoMap[key] = value
		}
		if err == io.EOF {
			err = nil
			break
		}
		conent += line
	}
	json := ToJSON(infoMap)
	info = &CerInfo{}
	ToBean([]byte(json), info)

	if info.No == "" {
		err = errors.New("证书信息加载失败！")
	}
	return
}

func getInfoValue(line string, key string) (value string) {
	value = line[strings.Index(line, key)+len(key):]
	value = strings.TrimSpace(value)
	return
}

func ToJSON(data interface{}) string {
	if data != nil {
		bs, _ := json.Marshal(data)
		if bs != nil {
			return string(bs)
		}
	}
	return ""
}

func ToBean(bytes []byte, req interface{}) (err error) {
	err = json.Unmarshal(bytes, req)
	return
}
