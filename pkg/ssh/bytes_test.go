package ssh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"teamide/pkg/util"
	"testing"
)

// 230 24 214 176 229 187 186 32 230 24 214 24 199 230 24 220 172 230 24 214 24 199 230 161 163 46 116 120 116

// 230 150 176 229 187 186 32 230 150 135 230 156 172 230 150 135 230 161 163 46 116 120 116

func TestBytes(t *testing.T) {
	bs := []byte{230, 24, 214, 176, 229, 187, 186, 32, 230, 24, 214, 24, 199, 230, 24, 220, 172, 230, 24, 214, 24, 199, 230, 161, 163, 46, 116, 120, 116}

	name := "新建 文本文档.txt"
	nameBS := []byte(name)

	tb := Utf8ArrayToStr(bs)

	fmt.Println("bs:", bs)
	fmt.Println("name:", name)
	fmt.Println("nameBS:", nameBS)
	fmt.Println("bs to:", tb)
	fmt.Println("bs to string:", string(tb))
	util.Logger.Info("bs logger", zap.Uint8s("Uint8s", bs))
	util.Logger.Info("bs logger", zap.Any("Any", bs))

	var s interface{}
	_ = json.Unmarshal(bs, &s)

	fmt.Println("bs json:", s)

}

func Utf8ArrayToStr(array []byte) []byte {
	var i int
	var length int
	var c byte
	var char2 byte
	var char3 byte
	var bs []byte
	length = len(array)
	i = 0
	for i < length {
		c = array[i]
		i++
		switch c >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			bs = append(bs, c)
			break
		case 12, 13:
			char2 = array[i]
			i++

			bs = append(bs, ((c&0x1F)<<6)|(char2&0x3F))
			break
		case 14:
			char2 = array[i]
			i++
			char3 = array[i]
			i++
			bs = append(bs, (c&0x0F)<<12|
				((char2&0x3F)<<6)|
				((char3&0x3F)<<0))
			break
		}
	}
	return bs
}

//GBK -> UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.HZGB2312.NewEncoder())
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return all, err
	}
	return all, nil
}

//UTF-8 -> GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return all, err
	}
	return all, nil
}
