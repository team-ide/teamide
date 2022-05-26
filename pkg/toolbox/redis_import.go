package toolbox

import (
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
	Count     int                 `json:"count,omitempty"`
	Key       string              `json:"key,omitempty"`
	ValueType string              `json:"valueType,omitempty"`
	Value     string              `json:"value,omitempty"`
	ValueList []map[string]string `json:"valueList,omitempty"`
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

	//var dataList []map[string]interface{}
	batchNumber := this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 10
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

		if strategyData.ValueType == "string" {

		}

	}
	return
}
