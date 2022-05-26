package toolbox

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
	redisImportTaskCache = map[string]*redisImportTask{}
)

func addRedisImportTask(task *redisImportTask) {
	redisImportTaskCache[task.Key] = task
	go task.Start()
}

type RedisStrategyData struct {
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

type redisImportTask struct {
	request          *RedisBaseRequest
	generateParam    *db.GenerateParam
	Key              string               `json:"key,omitempty"`
	ImportType       string               `json:"importType,omitempty"`
	StrategyDataList []*RedisStrategyData `json:"strategyDataList,omitempty"`
	BatchNumber      int                  `json:"batchNumber,omitempty"`
	DataCount        int                  `json:"dataCount"`
	ReadyDataCount   int                  `json:"readyDataCount"`
	SuccessCount     int                  `json:"successCount"`
	ErrorCount       int                  `json:"errorCount"`
	IsEnd            bool                 `json:"isEnd,omitempty"`
	StartTime        time.Time            `json:"startTime,omitempty"`
	EndTime          time.Time            `json:"endTime,omitempty"`
	Error            string               `json:"error,omitempty"`
	UseTime          int64                `json:"useTime"`
	IsStop           bool                 `json:"isStop"`
	service          RedisService
	ctx              context.Context
}

func (this_ *redisImportTask) Stop() {
	this_.IsStop = true
}
func (this_ *redisImportTask) Start() {
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

	this_.ctx = context.TODO()
	if this_.ImportType == "strategy" {
		err := this_.doStrategy()
		if err != nil {
			panic(err)
		}
	}

}

func (this_ *redisImportTask) doStrategy() (err error) {
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
		err = this_.doStrategyData(this_.request.Database, strategyData)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *redisImportTask) doStrategyData(database int, strategyData *RedisStrategyData) (err error) {
	if strategyData.Count <= 0 {
		return
	}
	if this_.IsStop {
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
			err = this_.service.Set(this_.ctx, database, key, value)
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
				err = this_.service.LPush(this_.ctx, database, key, value)
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
				err = this_.service.SAdd(this_.ctx, database, key, value)
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
				err = this_.service.HSet(this_.ctx, database, key, hashKey, value)
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
