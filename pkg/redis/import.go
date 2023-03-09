package redis

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"go.uber.org/zap"
//	"teamide/pkg/data_engine"
//	"github.com/team-ide/go-tool/util"
//	"time"
//)
//
//var (
//	ImportTaskCache = map[string]*ImportTask{}
//)
//
//func StartImportTask(task *ImportTask) {
//	ImportTaskCache[task.Key] = task
//	go task.Start()
//}
//
//func GetImportTask(taskKey string) *ImportTask {
//	task := ImportTaskCache[taskKey]
//	return task
//}
//
//func StopImportTask(taskKey string) *ImportTask {
//	task := ImportTaskCache[taskKey]
//	if task != nil {
//		task.Start()
//	}
//	return task
//}
//
//func CleanImportTask(taskKey string) *ImportTask {
//	task := ImportTaskCache[taskKey]
//	if task != nil {
//		delete(ImportTaskCache, taskKey)
//	}
//	return task
//}
//
//type StrategyData struct {
//	Count      int    `json:"count,omitempty"`
//	Key        string `json:"key,omitempty"`
//	ValueType  string `json:"valueType,omitempty"`
//	Value      string `json:"value,omitempty"`
//	ValueCount int    `json:"valueCount,omitempty"`
//	ListValue  string `json:"listValue,omitempty"`
//	SetValue   string `json:"setValue,omitempty"`
//	HashKey    string `json:"hashKey,omitempty"`
//	HashValue  string `json:"hashValue,omitempty"`
//}
//
//type ImportTask struct {
//	Database         int             `json:"database,omitempty"`
//	Key              string          `json:"key,omitempty"`
//	ImportType       string          `json:"importType,omitempty"`
//	StrategyDataList []*StrategyData `json:"strategyDataList,omitempty"`
//	BatchNumber      int             `json:"batchNumber,omitempty"`
//	DataCount        int             `json:"dataCount"`
//	ReadyDataCount   int             `json:"readyDataCount"`
//	SuccessCount     int             `json:"successCount"`
//	ErrorCount       int             `json:"errorCount"`
//	IsEnd            bool            `json:"isEnd,omitempty"`
//	StartTime        time.Time       `json:"startTime,omitempty"`
//	EndTime          time.Time       `json:"endTime,omitempty"`
//	Error            string          `json:"error,omitempty"`
//	UseTime          int64           `json:"useTime"`
//	IsStop           bool            `json:"isStop"`
//	Service          Service         `json:"-"`
//	taskList         []*data_engine.StrategyTask
//}
//
//func (this_ *ImportTask) Stop() {
//	this_.IsStop = true
//	for _, t := range this_.taskList {
//		t.Stop()
//	}
//}
//
//func (this_ *ImportTask) Start() {
//	this_.StartTime = time.Now()
//	defer func() {
//		if err := recover(); err != nil {
//			util.Logger.Error("导入数据异常", zap.Any("error", err))
//			this_.Error = fmt.Sprint(err)
//		}
//		this_.EndTime = time.Now()
//		this_.IsEnd = true
//		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
//	}()
//
//	if this_.ImportType == "strategy" {
//		err := this_.doStrategy()
//		if err != nil {
//			panic(err)
//		}
//	}
//
//}
//
//func (this_ *ImportTask) doStrategy() (err error) {
//
//	for _, strategyData := range this_.StrategyDataList {
//		if strategyData.Count <= 0 {
//			strategyData.Count = 0
//		}
//
//		if strategyData.Key == "" {
//			err = errors.New("必须配置Key")
//			return
//		}
//		if strategyData.ValueType == "" {
//			err = errors.New("必须配置值类型")
//			return
//		}
//		switch strategyData.ValueType {
//		case "string":
//			this_.DataCount += strategyData.Count
//		case "list":
//			this_.DataCount += strategyData.Count * strategyData.ValueCount
//		case "set":
//			this_.DataCount += strategyData.Count * strategyData.ValueCount
//		case "hash":
//			this_.DataCount += strategyData.Count * strategyData.ValueCount
//		default:
//			err = errors.New("不支持的值类型[" + strategyData.ValueType + "]")
//		}
//
//	}
//
//	for _, strategyData := range this_.StrategyDataList {
//		if this_.needStop() {
//			break
//		}
//		err = this_.doStrategyData(this_.Database, strategyData)
//		if err != nil {
//			return
//		}
//	}
//	return
//}
//
//func (this_ *ImportTask) doStrategyData(database int, strategyData *StrategyData) (err error) {
//	if strategyData.Count <= 0 {
//		return
//	}
//	if this_.needStop() {
//		return
//	}
//
//	ctx := context.TODO()
//	client, err := this_.Service.GetClient(ctx, database)
//	if err != nil {
//		return
//	}
//
//	task := &data_engine.StrategyTask{}
//
//	taskStrategyData := &data_engine.StrategyData{}
//
//	task.StrategyDataList = append(task.StrategyDataList, taskStrategyData)
//
//	taskStrategyData.Count = strategyData.Count
//	taskStrategyData.FieldList = append(taskStrategyData.FieldList, &data_engine.StrategyDataField{
//		Name:  "key",
//		Value: strategyData.Key,
//	})
//	taskStrategyData.FieldList = append(taskStrategyData.FieldList, &data_engine.StrategyDataField{
//		Name:  "valueType",
//		Value: `"` + strategyData.ValueType + `"`,
//	})
//	switch strategyData.ValueType {
//	case "string":
//		taskStrategyData.FieldList = append(taskStrategyData.FieldList, &data_engine.StrategyDataField{
//			Name:  "value",
//			Value: strategyData.Value,
//		})
//	}
//	var newValueTaskStrategyData = func(key string) (valueTaskStrategyData *data_engine.StrategyData) {
//		valueTaskStrategyData = &data_engine.StrategyData{}
//		valueTaskStrategyData.IndexName = "_$value_index"
//		valueTaskStrategyData.Count = strategyData.ValueCount
//
//		valueTaskStrategyData.FieldList = append(valueTaskStrategyData.FieldList, &data_engine.StrategyDataField{
//			Name:  "key",
//			Value: `"` + key + `"`,
//		})
//
//		valueTaskStrategyData.FieldList = append(valueTaskStrategyData.FieldList, &data_engine.StrategyDataField{
//			Name:  "valueType",
//			Value: `"` + strategyData.ValueType + `"`,
//		})
//
//		switch strategyData.ValueType {
//		case "list":
//			valueTaskStrategyData.FieldList = append(valueTaskStrategyData.FieldList, &data_engine.StrategyDataField{
//				Name:  "value",
//				Value: strategyData.ListValue,
//			})
//		case "set":
//			valueTaskStrategyData.FieldList = append(valueTaskStrategyData.FieldList, &data_engine.StrategyDataField{
//				Name:  "value",
//				Value: strategyData.SetValue,
//			})
//		case "hash":
//			valueTaskStrategyData.FieldList = append(valueTaskStrategyData.FieldList, &data_engine.StrategyDataField{
//				Name:  "hashKey",
//				Value: strategyData.HashKey,
//			})
//			valueTaskStrategyData.FieldList = append(valueTaskStrategyData.FieldList, &data_engine.StrategyDataField{
//				Name:  "value",
//				Value: strategyData.HashValue,
//			})
//		}
//		return
//	}
//
//	task.OnError = func(onErr error) {
//		err = onErr
//	}
//
//	this_.taskList = append(this_.taskList, task)
//
//	task.OnData = func(onData map[string]interface{}) (err error) {
//
//		if this_.needStop() {
//			return
//		}
//
//		valueType := onData["valueType"].(string)
//		var key string
//		key, err = util.GetStringValue(onData["key"])
//		if err != nil {
//			return
//		}
//		var value interface{}
//		var valueOK bool
//		value, valueOK = onData["value"]
//		var valueString string
//		if valueOK {
//			valueString, err = util.GetStringValue(value)
//			if err != nil {
//				return
//			}
//		}
//		switch valueType {
//		case "string":
//			this_.ReadyDataCount++
//			err = Set(ctx, client, key, valueString)
//			if err != nil {
//				this_.ErrorCount++
//				return
//			}
//			this_.SuccessCount++
//
//		case "list":
//			if valueOK {
//				this_.ReadyDataCount++
//				err = LPush(ctx, client, key, valueString)
//				if err != nil {
//					this_.ErrorCount++
//					return
//				}
//				this_.SuccessCount++
//			} else {
//
//				valueTask := &data_engine.StrategyTask{}
//				this_.taskList = append(this_.taskList, valueTask)
//				valueTask.StrategyDataList = append(valueTask.StrategyDataList, newValueTaskStrategyData(key))
//				valueTask.OnData = task.OnData
//				valueTask.OnError = task.OnError
//				valueTask.OnEnd = task.OnEnd
//				valueTask.Start()
//			}
//		case "set":
//			if valueOK {
//				this_.ReadyDataCount++
//				err = SAdd(ctx, client, key, valueString)
//				if err != nil {
//					this_.ErrorCount++
//					return
//				}
//				this_.SuccessCount++
//
//			} else {
//				valueTask := &data_engine.StrategyTask{}
//				this_.taskList = append(this_.taskList, valueTask)
//				valueTask.StrategyDataList = append(valueTask.StrategyDataList, newValueTaskStrategyData(key))
//				valueTask.OnData = task.OnData
//				valueTask.OnError = task.OnError
//				valueTask.OnEnd = task.OnEnd
//				valueTask.Start()
//			}
//		case "hash":
//			if valueOK {
//				var hashKey string
//				hashKey, err = util.GetStringValue(onData["hashKey"])
//				if err != nil {
//					return
//				}
//				this_.ReadyDataCount++
//				err = HSet(ctx, client, key, hashKey, valueString)
//				if err != nil {
//					this_.ErrorCount++
//					return
//				}
//				this_.SuccessCount++
//			} else {
//				valueTask := &data_engine.StrategyTask{}
//				this_.taskList = append(this_.taskList, valueTask)
//				valueTask.StrategyDataList = append(valueTask.StrategyDataList, newValueTaskStrategyData(key))
//				valueTask.OnData = task.OnData
//				valueTask.OnError = task.OnError
//				valueTask.OnEnd = task.OnEnd
//				valueTask.Start()
//			}
//		}
//
//		return
//	}
//
//	task.OnEnd = func() {
//
//	}
//
//	task.Start()
//
//	return
//}
//
//func (this_ *ImportTask) needStop() bool {
//	if this_.IsStop || this_.IsEnd {
//		return true
//	}
//	return false
//}
