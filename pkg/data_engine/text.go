package data_engine

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"time"
)

type TextTask struct {
	Path           string    `json:"path,omitempty"`
	CellSeparator  string    `json:"cellSeparator,omitempty"`
	SkipRow        int       `json:"skipRow,omitempty"`
	NameList       []string  `json:"nameList,omitempty"`
	DataCount      int       `json:"dataCount"`
	ReadyDataCount int       `json:"readyDataCount"`
	IsEnd          bool      `json:"isEnd,omitempty"`
	StartTime      time.Time `json:"startTime,omitempty"`
	EndTime        time.Time `json:"endTime,omitempty"`
	Error          string    `json:"error,omitempty"`
	UseTime        int64     `json:"useTime"`
	IsStop         bool      `json:"isStop"`
	IsError        bool      `json:"isError"`
	OnData         func(data map[string]interface{}) (err error)
	OnError        func(err error)
	OnEnd          func()
}

func (this_ *TextTask) Stop() {
	this_.IsStop = true
}

func (this_ *TextTask) Start() {
	this_.StartTime = time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			err, ok := rec.(error)
			if ok {
				util.Logger.Error("数据读取异常", zap.Any("error", err))
				this_.Error = fmt.Sprint(err)
				this_.IsError = true
				this_.OnError(err)
			}
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		this_.UseTime = util.GetTimeByTime(this_.EndTime) - util.GetTimeByTime(this_.StartTime)
		util.Logger.Info("数据读取结束")
		this_.OnEnd()
	}()
	util.Logger.Info("数据读取开始")
	err := this_.do()

	if err != nil {
		panic(err)
	}

	return
}

func (this_ *TextTask) do() (err error) {
	if this_.Path == "" {
		err = errors.New("文件地址不能为空")
		return
	}
	if this_.CellSeparator == "" {
		err = errors.New("列分隔符不能为空")
		return
	}
	if len(this_.NameList) == 0 {
		return
	}
	util.Logger.Info("读取文件", zap.Any("path", this_.Path))
	textF, err := os.Open(this_.Path)
	if err != nil {
		return
	}
	buf := bufio.NewReader(textF)
	var line string
	for {
		line, err = buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				err = nil
				break
			} else {
				return
			}
		}
		separator := this_.CellSeparator
		values := strings.Split(line, separator)

		var data = map[string]interface{}{}

		hasValue := false
		for cellIndex, name := range this_.NameList {
			if cellIndex >= len(values) {
				break
			}
			var value = values[cellIndex]
			data[name] = value
			if value != "" {
				hasValue = true
			}
		}

		if !hasValue {
			continue
		}

		this_.DataCount++
		this_.ReadyDataCount++
		err = this_.OnData(data)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *TextTask) needStop() bool {
	if this_.IsStop || this_.IsEnd {
		return true
	}
	return false
}
