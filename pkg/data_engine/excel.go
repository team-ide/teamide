package data_engine

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
	"teamide/pkg/util"
	"time"
)

type ExcelSheet struct {
	Name     int      `json:"name,omitempty"`
	SkipRow  int      `json:"skipRow,omitempty"`
	NameList []string `json:"nameList,omitempty"`
}

type ExcelTask struct {
	Path           string        `json:"path,omitempty"`
	SheetList      []*ExcelSheet `json:"sheetList,omitempty"`
	DataCount      int           `json:"dataCount"`
	ReadyDataCount int           `json:"readyDataCount"`
	IsEnd          bool          `json:"isEnd,omitempty"`
	StartTime      time.Time     `json:"startTime,omitempty"`
	EndTime        time.Time     `json:"endTime,omitempty"`
	Error          string        `json:"error,omitempty"`
	UseTime        int64         `json:"useTime"`
	IsStop         bool          `json:"isStop"`
	IsError        bool          `json:"isError"`
	OnData         func(data map[string]interface{}) (err error)
	OnError        func(err error)
	OnEnd          func()
}

func (this_ *ExcelTask) Stop() {
	this_.IsStop = true
}

func (this_ *ExcelTask) Start() {
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
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
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

func (this_ *ExcelTask) do() (err error) {
	if this_.Path == "" {
		err = errors.New("文件地址不能为空")
		return
	}
	if len(this_.SheetList) == 0 {
		return
	}
	util.Logger.Info("读取文件", zap.Any("path", this_.Path))
	xlsxF, err := xlsx.OpenFile(this_.Path)
	if err != nil {
		return
	}

	for sheetIndex, sheetConfig := range this_.SheetList {

		if this_.needStop() {
			return
		}
		if sheetIndex >= len(xlsxF.Sheets) {
			break
		}
		sheet := xlsxF.Sheets[sheetIndex]
		maxRow := sheet.MaxRow
		util.Logger.Info("读取Sheet", zap.Any("Sheet", sheetIndex), zap.Any("行数", maxRow))
		nameList := sheetConfig.NameList

		skipRow := sheetConfig.SkipRow
		if skipRow < 0 {
			skipRow = 0
		}
		for rowIndex := skipRow; rowIndex < maxRow; rowIndex++ {

			if this_.needStop() {
				return
			}

			row := sheet.Rows[rowIndex]

			var data = map[string]interface{}{}

			hasValue := false
			for cellIndex, name := range nameList {
				if cellIndex >= len(row.Cells) {
					break
				}
				cell := row.Cells[cellIndex]
				var value = cell.String()
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
	}
	return
}

func (this_ *ExcelTask) needStop() bool {
	if this_.IsStop || this_.IsEnd {
		return true
	}
	return false
}
