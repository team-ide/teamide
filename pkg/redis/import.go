package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"go.uber.org/zap"
	"teamide/pkg/db"
	"teamide/pkg/javascript"
	"teamide/pkg/util"
	"time"
)

var (
	ImportTaskCache = map[string]*ImportTask{}
)

func StartImportTask(task *ImportTask) {
	ImportTaskCache[task.Key] = task
	go task.Start()
}

func GetImportTask(taskKey string) *ImportTask {
	task := ImportTaskCache[taskKey]
	return task
}

func StopImportTask(taskKey string) *ImportTask {
	task := ImportTaskCache[taskKey]
	if task != nil {
		task.Start()
	}
	return task
}

func CleanImportTask(taskKey string) *ImportTask {
	task := ImportTaskCache[taskKey]
	if task != nil {
		delete(ImportTaskCache, taskKey)
	}
	return task
}

type StrategyData struct {
	Count      int    `json:"count,omitempty"`
	Key        string `json:"key,omitempty"`
	ValueType  string `json:"valueType,omitempty"`
	Value      string `json:"value,omitempty"`
	ValueCount int    `json:"valueCount,omitempty"`
	ListValue  string `json:"listValue,omitempty"`
	SetValue   string `json:"setValue,omitempty"`
	HashKey    string `json:"hashKey,omitempty"`
	HashValue  string `json:"hashValue,omitempty"`
}

type ImportTask struct {
	Database         int             `json:"database,omitempty"`
	Key              string          `json:"key,omitempty"`
	ImportType       string          `json:"importType,omitempty"`
	StrategyDataList []*StrategyData `json:"strategyDataList,omitempty"`
	BatchNumber      int             `json:"batchNumber,omitempty"`
	DataCount        int             `json:"dataCount"`
	ReadyDataCount   int             `json:"readyDataCount"`
	SuccessCount     int             `json:"successCount"`
	ErrorCount       int             `json:"errorCount"`
	IsEnd            bool            `json:"isEnd,omitempty"`
	StartTime        time.Time       `json:"startTime,omitempty"`
	EndTime          time.Time       `json:"endTime,omitempty"`
	Error            string          `json:"error,omitempty"`
	UseTime          int64           `json:"useTime"`
	IsStop           bool            `json:"isStop"`
	Service          Service         `json:"-"`
}

func (this_ *ImportTask) Stop() {
	this_.IsStop = true
}

func (this_ *ImportTask) Start() {
	this_.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("导入数据异常", zap.Any("error", err))
			this_.Error = fmt.Sprint(err)
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
	}()

	if this_.ImportType == "strategy" {
		err := this_.doStrategy()
		if err != nil {
			panic(err)
		}
	}

}

func (this_ *ImportTask) doStrategy() (err error) {
	for _, strategyData := range this_.StrategyDataList {
		if strategyData.Count <= 0 {
			strategyData.Count = 0
		}
		this_.DataCount += strategyData.Count
	}

	for _, strategyData := range this_.StrategyDataList {
		if this_.IsStop {
			break
		}
		err = this_.doStrategyData(this_.Database, strategyData)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *ImportTask) doStrategyData(database int, strategyData *StrategyData) (err error) {
	if strategyData.Count <= 0 {
		return
	}
	if this_.IsStop {
		return
	}

	ctx := context.TODO()
	client, err := this_.Service.GetClient(ctx, database)
	if err != nil {
		return
	}

	scriptContext := javascript.GetContext()

	vm := goja.New()

	for key, value := range scriptContext {
		err = vm.Set(key, value)
		if err != nil {
			return
		}
	}

	for i := 0; i < strategyData.Count; i++ {
		//data := map[string]interface{}{}
		err = vm.Set("_$index", i)
		if err != nil {
			return
		}

		if strategyData.Key == "" {
			err = errors.New("必须配置Key")
			return
		}
		if strategyData.ValueType == "" {
			err = errors.New("必须配置值类型")
			return
		}
		var key = strategyData.Key

		var scriptValue goja.Value
		if scriptValue, err = vm.RunString(key); err != nil {
			util.Logger.Error("表达式执行异常", zap.Any("script", key), zap.Error(err))
			return
		}
		key = db.GetStringValue(scriptValue.Export())

		switch strategyData.ValueType {
		case "string":
			var value = strategyData.Value
			if value != "" {
				if scriptValue, err = vm.RunString(value); err != nil {
					util.Logger.Error("表达式执行异常", zap.Any("script", value), zap.Error(err))
					return
				}
				value = db.GetStringValue(scriptValue.Export())
			}
			this_.ReadyDataCount++
			err = Set(ctx, client, key, value)
			if err != nil {
				this_.ErrorCount++
				return
			}
			this_.SuccessCount++

		case "list":
			for valueIndex := 0; valueIndex < strategyData.ValueCount; valueIndex++ {
				err = vm.Set("_$value_index", valueIndex)
				if err != nil {
					return
				}
				var value = strategyData.ListValue
				if value != "" {
					if scriptValue, err = vm.RunString(value); err != nil {
						util.Logger.Error("表达式执行异常", zap.Any("script", value), zap.Error(err))
						return
					}
					value = db.GetStringValue(scriptValue.Export())
				}
				this_.ReadyDataCount++
				err = LPush(ctx, client, key, value)
				if err != nil {
					this_.ErrorCount++
					return
				}
				this_.SuccessCount++
			}

		case "set":
			for valueIndex := 0; valueIndex < strategyData.ValueCount; valueIndex++ {
				err = vm.Set("_$value_index", valueIndex)
				if err != nil {
					return
				}
				var value = strategyData.SetValue
				if value != "" {
					if scriptValue, err = vm.RunString(value); err != nil {
						util.Logger.Error("表达式执行异常", zap.Any("script", value), zap.Error(err))
						return
					}
					value = db.GetStringValue(scriptValue.Export())
				}
				this_.ReadyDataCount++
				err = SAdd(ctx, client, key, value)
				if err != nil {
					this_.ErrorCount++
					return
				}
				this_.SuccessCount++
			}

		case "hash":
			for valueIndex := 0; valueIndex < strategyData.ValueCount; valueIndex++ {
				err = vm.Set("_$value_index", valueIndex)
				if err != nil {
					return
				}
				var hashKey = strategyData.HashKey
				if hashKey != "" {
					if scriptValue, err = vm.RunString(hashKey); err != nil {
						util.Logger.Error("表达式执行异常", zap.Any("script", hashKey), zap.Error(err))
						return
					}
					hashKey = db.GetStringValue(scriptValue.Export())
				}
				var value = strategyData.HashValue
				if value != "" {
					if scriptValue, err = vm.RunString(value); err != nil {
						util.Logger.Error("表达式执行异常", zap.Any("script", value), zap.Error(err))
						return
					}
					value = db.GetStringValue(scriptValue.Export())
				}
				this_.ReadyDataCount++
				err = HSet(ctx, client, key, hashKey, value)
				if err != nil {
					this_.ErrorCount++
					return
				}
				this_.SuccessCount++
			}
		default:
			err = errors.New("不支持的值类型[" + strategyData.ValueType + "]")
		}

	}
	return
}
