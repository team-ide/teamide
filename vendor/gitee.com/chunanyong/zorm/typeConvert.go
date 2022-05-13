package zorm

import (
	"fmt"
	"strconv"
	"time"

	"gitee.com/chunanyong/zorm/decimal"
)

func typeConvertFloat32(i interface{}) float32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(float32); ok {
		return v
	}
	v, _ := strconv.ParseFloat(typeConvertString(i), 32)
	return float32(v)
}

func typeConvertFloat64(i interface{}) float64 {
	if i == nil {
		return 0
	}
	if v, ok := i.(float64); ok {
		return v
	}
	v, _ := strconv.ParseFloat(typeConvertString(i), 64)
	return v
}

func typeConvertDecimal(i interface{}) decimal.Decimal {
	if i == nil {
		return decimal.Zero
	}
	if v, ok := i.(decimal.Decimal); ok {
		return v
	}
	v, _ := decimal.NewFromString(typeConvertString(i))
	return v
}

func typeConvertInt64toInt(from int64) (int, error) {
	//int64 转 int
	strInt64 := strconv.FormatInt(from, 10)
	to, err := strconv.Atoi(strInt64)
	if err != nil {
		return -1, err
	}
	return to, err
}

func typeConvertTime(i interface{}, format string, TZLocation ...*time.Location) time.Time {
	s := typeConvertString(i)
	t, _ := typeConvertStrToTime(s, format, TZLocation...)
	return t
}

func typeConvertStrToTime(str string, format string, TZLocation ...*time.Location) (time.Time, error) {
	if len(TZLocation) > 0 {
		t, err := time.ParseInLocation(format, str, TZLocation[0])
		if err == nil {
			return t, nil
		}
		return time.Time{}, err

	}
	t, err := time.ParseInLocation(format, str, time.Local)
	if err == nil {
		return t, nil
	}
	return time.Time{}, err

}

func typeConvertInt64(i interface{}) int64 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int64); ok {
		return v
	}
	return int64(typeConvertInt(i))
}

func typeConvertString(i interface{}) string {
	if i == nil {
		return ""
	}
	switch value := i.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.Itoa(int(value))
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(uint64(value), 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

//false: "", 0, false, off
func typeConvertBool(i interface{}) bool {
	if i == nil {
		return false
	}
	if v, ok := i.(bool); ok {
		return v
	}
	if s := typeConvertString(i); s != "" && s != "0" && s != "false" && s != "off" {
		return true
	}
	return false
}

func typeConvertInt(i interface{}) int {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case int:
		return value
	case int8:
		return int(value)
	case int16:
		return int(value)
	case int32:
		return int(value)
	case int64:
		return int(value)
	case uint:
		return int(value)
	case uint8:
		return int(value)
	case uint16:
		return int(value)
	case uint32:
		return int(value)
	case uint64:
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	case bool:
		if value {
			return 1
		}
		return 0
	default:
		v, _ := strconv.Atoi(typeConvertString(value))
		return v
	}
}

/*
func encodeString(s string) []byte {
	return []byte(s)
}

func decodeToString(b []byte) string {
	return string(b)
}

func encodeBool(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}

}

func encodeInt(i int) []byte {
	if i <= math.MaxInt8 {
		return encodeInt8(int8(i))
	} else if i <= math.MaxInt16 {
		return encodeInt16(int16(i))
	} else if i <= math.MaxInt32 {
		return encodeInt32(int32(i))
	} else {
		return encodeInt64(int64(i))
	}
}

func encodeUint(i uint) []byte {
	if i <= math.MaxUint8 {
		return encodeUint8(uint8(i))
	} else if i <= math.MaxUint16 {
		return encodeUint16(uint16(i))
	} else if i <= math.MaxUint32 {
		return encodeUint32(uint32(i))
	} else {
		return encodeUint64(uint64(i))
	}
}

func encodeInt8(i int8) []byte {
	return []byte{byte(i)}
}

func encodeUint8(i uint8) []byte {
	return []byte{byte(i)}
}

func encodeInt16(i int16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(i))
	return bytes
}

func encodeUint16(i uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, i)
	return bytes
}

func encodeInt32(i int32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(i))
	return bytes
}

func encodeUint32(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, i)
	return bytes
}

func encodeInt64(i int64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(i))
	return bytes
}

func encodeUint64(i uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, i)
	return bytes
}

func encodeFloat32(f float32) []byte {
	bits := math.Float32bits(f)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func encodeFloat64(f float64) []byte {
	bits := math.Float64bits(f)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func encode(vs ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for i := 0; i < len(vs); i++ {
		switch value := vs[i].(type) {
		case int:
			buf.Write(encodeInt(value))
		case int8:
			buf.Write(encodeInt8(value))
		case int16:
			buf.Write(encodeInt16(value))
		case int32:
			buf.Write(encodeInt32(value))
		case int64:
			buf.Write(encodeInt64(value))
		case uint:
			buf.Write(encodeUint(value))
		case uint8:
			buf.Write(encodeUint8(value))
		case uint16:
			buf.Write(encodeUint16(value))
		case uint32:
			buf.Write(encodeUint32(value))
		case uint64:
			buf.Write(encodeUint64(value))
		case bool:
			buf.Write(encodeBool(value))
		case string:
			buf.Write(encodeString(value))
		case []byte:
			buf.Write(value)
		case float32:
			buf.Write(encodeFloat32(value))
		case float64:
			buf.Write(encodeFloat64(value))
		default:
			if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
				buf.Write(encodeString(fmt.Sprintf("%v", value)))
			}
		}
	}
	return buf.Bytes()
}

func isNumeric(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < byte('0') || s[i] > byte('9') {
			return false
		}
	}
	return true
}
func typeConvertTimeDuration(i interface{}) time.Duration {
	return time.Duration(typeConvertInt64(i))
}

func typeConvertBytes(i interface{}) []byte {
	if i == nil {
		return nil
	}
	if r, ok := i.([]byte); ok {
		return r
	}
	return encode(i)

}

func typeConvertStrings(i interface{}) []string {
	if i == nil {
		return nil
	}
	if r, ok := i.([]string); ok {
		return r
	} else if r, ok := i.([]interface{}); ok {
		strs := make([]string, len(r))
		for k, v := range r {
			strs[k] = typeConvertString(v)
		}
		return strs
	}
	return []string{fmt.Sprintf("%v", i)}
}

func typeConvertInt8(i interface{}) int8 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int8); ok {
		return v
	}
	return int8(typeConvertInt(i))
}

func typeConvertInt16(i interface{}) int16 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int16); ok {
		return v
	}
	return int16(typeConvertInt(i))
}

func typeConvertInt32(i interface{}) int32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int32); ok {
		return v
	}
	return int32(typeConvertInt(i))
}

func typeConvertUint(i interface{}) uint {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case int:
		return uint(value)
	case int8:
		return uint(value)
	case int16:
		return uint(value)
	case int32:
		return uint(value)
	case int64:
		return uint(value)
	case uint:
		return value
	case uint8:
		return uint(value)
	case uint16:
		return uint(value)
	case uint32:
		return uint(value)
	case uint64:
		return uint(value)
	case float32:
		return uint(value)
	case float64:
		return uint(value)
	case bool:
		if value {
			return 1
		}
		return 0
	default:
		v, _ := strconv.ParseUint(typeConvertString(value), 10, 64)
		return uint(v)
	}
}

func typeConvertUint8(i interface{}) uint8 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint8); ok {
		return v
	}
	return uint8(typeConvertUint(i))
}

func typeConvertUint16(i interface{}) uint16 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint16); ok {
		return v
	}
	return uint16(typeConvertUint(i))
}

func typeConvertUint32(i interface{}) uint32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint32); ok {
		return v
	}
	return uint32(typeConvertUint(i))
}

func typeConvertUint64(i interface{}) uint64 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint64); ok {
		return v
	}
	return uint64(typeConvertUint(i))
}
*/
