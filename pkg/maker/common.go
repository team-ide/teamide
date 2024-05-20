package maker

import (
	"fmt"
	"strings"
	"time"
)

func SetBySlash(data map[string]interface{}, name string, value any) {
	//fmt.Println("SetBySlash:", name)
	index := strings.Index(name, "/")
	if index < 0 {
		data[name] = value
		return
	}
	pName := name[:index]
	cName := name[index+1:]
	//fmt.Println("SetBySlash pName:", pName, "cName:", cName)
	parent := data[pName]
	if parent == nil {
		parent = map[string]interface{}{}
		data[pName] = parent
	}
	SetBySlash(parent.(map[string]interface{}), cName, value)
}

func GetDurationFormatByMillisecond(millisecond int64) (formatString string) {
	if millisecond == 0 {
		return fmt.Sprintf("%d毫秒", 0)
	}

	duration := time.Duration(millisecond) * time.Millisecond
	h := int(duration.Hours())
	m := int(duration.Minutes()) % 60
	s := int(duration.Seconds()) % 60
	ms := int(duration.Milliseconds()) % 1000
	if h > 0 {
		formatString = fmt.Sprintf("%d小时", h)
	}
	if m > 0 {
		formatString += fmt.Sprintf("%d分钟", m)
	}
	if s > 0 {
		formatString += fmt.Sprintf("%d秒", s)
	}
	if ms > 0 {
		formatString += fmt.Sprintf("%d毫秒", ms)
	}
	return
}

func invokeStart(name string, invokeData *InvokeData) (funcInvoke *FuncInvoke) {
	funcInvoke = &FuncInvoke{
		name:       name,
		invokeData: invokeData,
	}
	funcInvoke.start()
	return
}

type FuncInvoke struct {
	name       string
	startTime  time.Time
	endTime    time.Time
	err        error
	invokeData *InvokeData
}

func (this_ *FuncInvoke) start() {
	this_.startTime = time.Now()
}

func (this_ *FuncInvoke) end(err error) {
	this_.err = err
	this_.endTime = time.Now()
}

func (this_ *FuncInvoke) use() string {
	return GetDurationFormatByMillisecond(this_.endTime.UnixMilli() - this_.startTime.UnixMilli())
}
