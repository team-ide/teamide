package toolbox

import (
	"time"
)

//GetNowTime 获取当前时间戳
func GetNowTime() int64 {
	return time.Now().UnixNano() / 1e6
}
