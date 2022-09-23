package util

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"strings"
	"sync"
)

var fileTypeMap sync.Map

func init() {
	fileTypeMap.Store("ffd8ffe000104a464946", "jpg") //JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "png") //PNG (png)
	fileTypeMap.Store("47494638396126026f01", "gif") //GIF (gif)
	fileTypeMap.Store("49492a00227105008037", "tif") //TIFF (tif)
	fileTypeMap.Store("424d228c010000000000", "bmp") //16色位图(bmp)
	fileTypeMap.Store("424d8240090000000000", "bmp") //24位位图(bmp)
	fileTypeMap.Store("424d8e1b030000000000", "bmp") //256色位图(bmp)

	//fileTypeMap.Store("41433130313500000000", "dwg") //CAD (dwg)
	//fileTypeMap.Store("255044462d312e350d0a", "pdf") //Adobe Acrobat (pdf)
	//
	//fileTypeMap.Store("00000020667479706d70", "mp4")
	//fileTypeMap.Store("49443303000000002176", "mp3")
	//fileTypeMap.Store("000001ba210001000180", "mpg") //
	//fileTypeMap.Store("3026b2758e66cf11a6d9", "wmv") //wmv与asf相同
	//fileTypeMap.Store("52494646e27807005741", "wav") //Wave (wav)
	//fileTypeMap.Store("52494646d07d60074156", "avi")
	//fileTypeMap.Store("4d546864000000060001", "mid") //MIDI (mid)
	//
	//fileTypeMap.Store("504b0304140000000800", "zip")
	//fileTypeMap.Store("526172211a0700cf9073", "rar")
	//fileTypeMap.Store("504b03040a0000000000", "jar")
	//fileTypeMap.Store("1f8b0800000000000000", "gz") //gz文件
}

// bytesToHexString 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// GetFileType 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
func GetFileType(fSrc []byte) string {
	var fileType string
	fileCode := bytesToHexString(fSrc)

	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}
