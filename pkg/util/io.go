package util

import "io"

func Read(reader io.Reader, buf []byte, onRead func(n int) (err error)) (err error) {
	for {
		var n int
		n, err = reader.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		e := onRead(n)
		if e != nil {
			err = e
			break
		}
		if err == io.EOF {
			err = nil
			break
		}
	}
	return
}

func Write(writer io.Writer, buf []byte, onWrite func(n int) (err error)) (err error) {
	var bs = buf
	for {
		if len(bs) == 0 {
			break
		}
		var n int
		n, err = writer.Write(bs)
		if err != nil && err != io.EOF {
			break
		}
		if onWrite != nil {
			e := onWrite(n)
			if e != nil {
				err = e
				break
			}
		}
		if n < len(bs) {
			bs = bs[n:]
		} else {
			bs = []byte{}
		}
		if err == io.EOF {
			err = nil
			break
		}
	}
	return
}
