package util

import (
	"bytes"
	"compress/gzip"
	"io"
)

//GZipBytes 压缩
func GZipBytes(data []byte) ([]byte, error) {

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//UGZipBytes 解压
func UGZipBytes(data []byte) ([]byte, error) {
	var in bytes.Buffer
	in.Write(data)
	r, err := gzip.NewReader(&in)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	bs, err := io.ReadAll(r)
	return bs, err
}
