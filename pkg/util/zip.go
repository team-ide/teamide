package util

import (
	"bytes"
	"compress/zlib"
	"io"
)

//ZipBytes 压缩
func ZipBytes(data []byte) []byte {
	var in bytes.Buffer
	z := zlib.NewWriter(&in)
	z.Write(data)
	defer z.Close()
	return in.Bytes()
}

//UZipBytes 解压
func UZipBytes(data []byte) []byte {
	var out bytes.Buffer
	var in bytes.Buffer
	in.Write(data)
	r, _ := zlib.NewReader(&in)
	defer r.Close()
	io.Copy(&out, r)
	return out.Bytes()
}
