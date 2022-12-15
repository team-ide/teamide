package data_engine

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"teamide/pkg/javascript"
	"teamide/pkg/util"
	"time"
)

type DataStatistics struct {
	DataCount        int   `json:"dataCount"`
	DataSuccessCount int   `json:"dataSuccessCount"`
	DataErrorCount   int   `json:"dataErrorCount"`
	UseTime          int64 `json:"useTime"`

	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`

	countLock sync.Mutex
}

func (this_ *DataStatistics) IncrDataSuccessCount(num int, useTime int64) {
	this_.countLock.Lock()
	defer this_.countLock.Unlock()
	this_.DataSuccessCount += num
	this_.DataCount += num
	this_.UseTime += useTime
	return
}

func (this_ *DataStatistics) IncrDataErrorCount(num int, useTime int64) {
	this_.countLock.Lock()
	defer this_.countLock.Unlock()
	this_.DataErrorCount += num
	this_.UseTime += useTime
	return
}

type StrategyData struct {
	DataNumber  int                  `json:"dataNumber,omitempty"`
	StartIndex  int                  `json:"startIndex,omitempty"`
	BatchNumber int                  `json:"batchNumber,omitempty"`
	IndexName   string               `json:"indexName,omitempty"`
	DataName    string               `json:"dataName,omitempty"`
	FieldList   []*StrategyDataField `json:"fieldList,omitempty"`
	*DataStatistics
}

func (this_ *StrategyData) AddField(field *StrategyDataField) {
	this_.FieldList = append(this_.FieldList, field)
}

type StrategyDataField struct {
	Name        string               `json:"name,omitempty"`
	IndexName   string               `json:"indexName,omitempty"`
	DataName    string               `json:"dataName,omitempty"`
	Value       string               `json:"value,omitempty"`
	ReuseNumber int                  `json:"reuseNumber,omitempty"`
	ValueIsList bool                 `json:"valueIsList,omitempty"`
	ValueCount  int                  `json:"valueCount,omitempty"`
	FieldList   []*StrategyDataField `json:"fieldList,omitempty"`

	reuseValue      interface{}
	reuseValueCount int
}

type StrategyTask struct {
	*StrategyData
	IsEnd      bool   `json:"isEnd,omitempty"`
	Error      string `json:"error,omitempty"`
	IsStop     bool   `json:"isStop"`
	IsError    bool   `json:"isError"`
	OnDataList func(dataList []map[string]interface{}) (err error)
	OnError    func(err error)
	OnEnd      func()
}

func (this_ *StrategyTask) Stop() {
	this_.IsStop = true
}

func (this_ *StrategyTask) Start() {
	if this_.DataStatistics == nil {
		this_.DataStatistics = &DataStatistics{}
	}
	this_.StartTime = time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			err, ok := rec.(error)
			if ok {
				util.Logger.Error("数据生成异常", zap.Any("error", err))
				this_.Error = fmt.Sprint(err)
				this_.IsError = true
				this_.OnError(err)
			}
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		util.Logger.Info("数据生成结束")
		this_.OnEnd()
	}()
	util.Logger.Info("数据生成开始")
	err := this_.do()

	if err != nil {
		util.Logger.Error("数据生成异常", zap.Any("error", err))
		panic(err)
	}

	return
}

func (this_ *StrategyTask) do() (err error) {
	err = this_.doStrategyData(this_.StrategyData)
	if err != nil {
		return
	}
	return
}

func (this_ *StrategyTask) doStrategyData(strategyData *StrategyData) (err error) {
	if strategyData.DataNumber <= 0 {
		return
	}
	if this_.needStop() {
		return
	}
	if len(strategyData.FieldList) == 0 {
		err = errors.New("策略数据未配置字段")
		return
	}

	script, err := javascript.NewScript()
	if err != nil {
		return
	}
	dataNumber := strategyData.DataNumber
	batchNumber := strategyData.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 200
	}
	var data map[string]interface{}
	var dataList []map[string]interface{}
	for index := 0; index < dataNumber; index++ {

		if this_.needStop() {
			return
		}

		var dataIndex = index
		if this_.StartIndex > 0 {
			dataIndex = dataIndex + this_.StartIndex
		}

		var startTime = time.Now()
		data, err = this_.doStrategyDataFieldList(dataIndex, strategyData.IndexName, script, strategyData.FieldList)
		var endTime = time.Now()
		var useTime = util.GetTimeTime(endTime) - util.GetTimeTime(startTime)
		if err != nil {
			strategyData.IncrDataErrorCount(1, useTime)
			return
		}
		strategyData.IncrDataSuccessCount(1, useTime)
		dataList = append(dataList, data)

		if len(dataList) >= batchNumber {
			err = this_.OnDataList(dataList)
			dataList = []map[string]interface{}{}
			if err != nil {
				return
			}
		}
		if strategyData.DataName != "" {
			err = script.Set(strategyData.DataName, data)
		} else {
			err = script.Set("_$data", data)
		}
		if err != nil {
			return
		}
	}
	if len(dataList) > 0 {
		err = this_.OnDataList(dataList)
		dataList = []map[string]interface{}{}
		if err != nil {
			return
		}
	}
	return
}

func (this_ *StrategyTask) doStrategyDataFieldList(index int, indexName string, script *javascript.Script, fieldList []*StrategyDataField) (data map[string]interface{}, err error) {

	if this_.needStop() {
		return
	}
	if index >= 0 {
		if indexName != "" {
			err = script.Set(indexName, index)
		} else {
			err = script.Set("_$index", index)
		}
		if err != nil {
			return
		}
	}

	data = map[string]interface{}{}
	for _, field := range fieldList {
		var name string
		var value interface{}
		name, value, err = this_.doStrategyDataField(script, field)
		if err != nil {
			return
		}
		data[name] = value
		err = script.Set(name, value)
		if err != nil {
			return
		}

	}
	return
}

func (this_ *StrategyTask) doStrategyDataField(script *javascript.Script, field *StrategyDataField) (name string, value interface{}, err error) {

	defer func() {
		field.reuseValueCount++
		field.reuseValue = value
	}()
	name = field.Name
	if field.ReuseNumber <= 0 {
		field.ReuseNumber = 1
	}
	if field.reuseValueCount >= field.ReuseNumber {
		field.reuseValue = nil
		field.reuseValueCount = 0
	}

	if field.reuseValueCount > 0 {
		value = field.reuseValue
		return
	}

	valueCount := 1
	if field.ValueIsList {
		valueCount = field.ValueCount
	}

	var fieldScript *javascript.Script
	fieldScript, err = javascript.NewScriptByParent(script)
	if err != nil {
		return
	}

	var valueList []interface{}
	for valueIndex := 0; valueIndex < valueCount; valueIndex++ {

		if this_.needStop() {
			return
		}
		var setIndex = -1
		if field.ValueIsList {
			setIndex = valueIndex
		}

		var data interface{}
		if len(field.FieldList) > 0 {
			data, err = this_.doStrategyDataFieldList(setIndex, field.IndexName, fieldScript, field.FieldList)
		} else {
			data, err = fieldScript.GetScriptValue(field.Value)
		}
		if err != nil {
			return
		}
		if field.DataName != "" {
			err = fieldScript.Set(field.DataName, data)
		} else {
			err = fieldScript.Set("_$"+name+"_data", data)
		}
		if err != nil {
			return
		}
		valueList = append(valueList, data)
	}
	if field.ValueIsList {
		value = valueList
	} else {
		value = valueList[0]
	}
	return
}

func (this_ *StrategyTask) needStop() bool {
	if this_.IsStop || this_.IsEnd {
		return true
	}
	return false
}
