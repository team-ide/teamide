package node

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

var (
	LengthError = errors.New("读取流长度错误")
)

type Message struct {
	Token  string `json:"token,omitempty"`
	Id     string `json:"id,omitempty"`
	Method string `json:"method,omitempty"`
	Error  string `json:"error,omitempty"`
}

func ReadMessage(reader io.Reader) (message *Message, err error) {
	var bytes []byte

	bytes, err = ReadBytes(reader)
	if err != nil {
		return
	}
	message = &Message{}

	err = json.Unmarshal(bytes, &message)
	if err != nil {
		return
	}
	return
}

func WriteMessage(writer io.Writer, message *Message) (err error) {
	var bytes []byte

	bytes, err = json.Marshal(message)
	if err != nil {
		return
	}

	err = WriteBytes(writer, bytes)
	return
}

func ReadData(reader io.Reader) (data map[string]interface{}, err error) {
	var bytes []byte

	bytes, err = ReadBytes(reader)
	if err != nil {
		return
	}
	data = map[string]interface{}{}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return
	}
	return
}

func WriteData(writer io.Writer, data map[string]interface{}) (err error) {
	var bytes []byte

	bytes, err = json.Marshal(data)
	if err != nil {
		return
	}

	err = WriteBytes(writer, bytes)
	return
}

func ReadBytes(reader io.Reader) (bytes []byte, err error) {
	var buf []byte
	var n int

	buf = make([]byte, 4)
	n, err = reader.Read(buf)
	if err != nil {
		return
	}
	if n < 4 {
		err = LengthError
		return
	}

	length := int(binary.LittleEndian.Uint32(buf))
	if length < 0 {
		err = LengthError
		return
	}

	if length > 0 {
		bytes = make([]byte, length)
		n, err = reader.Read(bytes)
		if err != nil {
			return
		}
		if n < length {
			err = LengthError
			return
		}
	}
	return
}

func WriteBytes(writer io.Writer, bytes []byte) (err error) {
	var n int
	var length = len(bytes)

	lengthBytes := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(lengthBytes, uint32(length))

	n, err = writer.Write(lengthBytes)
	if err != nil {
		return
	}
	if n < 4 {
		err = LengthError
		return
	}
	if length == 0 {
		return
	}
	n, err = writer.Write(bytes)
	if err != nil {
		return
	}
	if n < length {
		err = LengthError
		return
	}
	return
}
