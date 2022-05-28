package data_engine

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"teamide/pkg/javascript"
	"teamide/pkg/util"
	"time"
)

type StrategyData struct {
	Count     int                  `json:"count,omitempty"`
	IndexName string               `json:"indexName,omitempty"`
	DataName  string               `json:"dataName,omitempty"`
	FieldList []*StrategyDataField `json:"fieldList,omitempty"`
}

type StrategyDataField struct {
	Name        string               `json:"name,omitempty"`
	IndexName   string               `json:"indexName,omitempty"`
	DataName    string               `json:"dataName,omitempty"`
	Value       string               `json:"value,omitempty"`
	ValueIsList bool                 `json:"valueIsList,omitempty"`
	ValueCount  int                  `json:"valueCount,omitempty"`
	FieldList   []*StrategyDataField `json:"fieldList,omitempty"`
}

type StrategyTask struct {
	StrategyDataList []*StrategyData `json:"strategyDataList,omitempty"`
	DataCount        int             `json:"dataCount"`
	ReadyDataCount   int             `json:"readyDataCount"`
	IsEnd            bool            `json:"isEnd,omitempty"`
	StartTime        time.Time       `json:"startTime,omitempty"`
	EndTime          time.Time       `json:"endTime,omitempty"`
	Error            string          `json:"error,omitempty"`
	UseTime          int64           `json:"useTime"`
	IsStop           bool            `json:"isStop"`
	IsError          bool            `json:"isError"`
	OnData           func(data map[string]interface{}) (err error)
	OnError          func(err error)
	OnEnd            func()
}

func (this_ *StrategyTask) Stop() {
	this_.IsStop = true
}

func (this_ *StrategyTask) Start() {
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
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
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
	if len(this_.StrategyDataList) == 0 {
		return
	}
	for _, strategyData := range this_.StrategyDataList {
		if strategyData.Count <= 0 {
			strategyData.Count = 0
		}
		this_.DataCount += strategyData.Count
	}

	for _, strategyData := range this_.StrategyDataList {
		if this_.needStop() {
			break
		}
		err = this_.doStrategyData(strategyData)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *StrategyTask) doStrategyData(strategyData *StrategyData) (err error) {
	if strategyData.Count <= 0 {
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
	for index := 0; index < strategyData.Count; index++ {

		if this_.needStop() {
			return
		}
		if strategyData.IndexName != "" {
			err = script.Set(strategyData.IndexName, index)
		} else {
			err = script.Set("_$index", index)
		}
		if err != nil {
			return
		}
		var data map[string]interface{}
		data, err = this_.doStrategyDataFieldList(script, strategyData.FieldList)
		if err != nil {
			return
		}
		this_.ReadyDataCount++
		err = this_.OnData(data)
		if err != nil {
			return
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
	return
}

func (this_ *StrategyTask) doStrategyDataFieldList(script *javascript.Script, fieldList []*StrategyDataField) (data map[string]interface{}, err error) {

	if this_.needStop() {
		return
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

	name = field.Name

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
		if field.ValueIsList {
			if field.IndexName != "" {
				err = fieldScript.Set(field.IndexName, valueIndex)
			} else {
				err = fieldScript.Set("_$"+name+"_index", valueIndex)
			}
			if err != nil {
				return
			}
		}

		var data interface{}
		if len(field.FieldList) > 0 {
			data, err = this_.doStrategyDataFieldList(fieldScript, field.FieldList)
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
