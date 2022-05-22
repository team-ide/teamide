package toolbox

import (
	"fmt"
	"github.com/dop251/goja"
	"go.uber.org/zap"
	"teamide/pkg/db"
	"teamide/pkg/javascript"
	"teamide/pkg/util"
	"time"
)

var (
	databaseImportTaskCache = map[string]*databaseImportTask{}
)

func addDatabaseImportTask(task *databaseImportTask) {
	databaseImportTaskCache[task.Key] = task
	go task.Start()
}

type databaseImportTask struct {
	request          *DatabaseBaseRequest
	generateParam    *db.GenerateParam
	Key              string                   `json:"key,omitempty"`
	ImportType       string                   `json:"importType,omitempty"`
	StrategyDataList []map[string]interface{} `json:"strategyDataList,omitempty"`
	BatchNumber      int                      `json:"batchNumber,omitempty"`
	DataCount        int                      `json:"dataCount"`
	ReadyDataCount   int                      `json:"readyDataCount"`
	SuccessCount     int                      `json:"successCount"`
	ErrorCount       int                      `json:"errorCount"`
	IsEnd            bool                     `json:"isEnd,omitempty"`
	StartTime        time.Time                `json:"startTime,omitempty"`
	EndTime          time.Time                `json:"endTime,omitempty"`
	Error            string                   `json:"error,omitempty"`
	UseTime          int64                    `json:"useTime"`
	IsStop           bool                     `json:"isStop"`
	service          DatabaseService
}

func (this_ *databaseImportTask) Stop() {
	this_.IsStop = true
}
func (this_ *databaseImportTask) Start() {
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

func (this_ *databaseImportTask) doStrategy() (err error) {
	for _, importData := range this_.StrategyDataList {
		importCount := 0
		if importData["_$importCount"] != nil {
			importCount = int(importData["_$importCount"].(float64))
		}
		if importCount <= 0 {
			importCount = 0
		}
		importData["_$importCount"] = importCount
		this_.DataCount += importCount
	}

	for _, importData := range this_.StrategyDataList {
		if this_.IsStop {
			break
		}
		err = this_.doStrategyData(this_.request.Database, this_.request.Table, this_.request.ColumnList, importData)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *databaseImportTask) doStrategyData(database, table string, columnList []*db.TableColumnModel, importData map[string]interface{}) (err error) {
	importCount := importData["_$importCount"].(int)
	if importCount <= 0 {
		return
	}
	if this_.IsStop {
		return
	}

	var dataList []map[string]interface{}
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

	for i := 0; i < importCount; i++ {
		data := map[string]interface{}{}
		err = vm.Set("_$index", i)
		if err != nil {
			return
		}

		for _, column := range columnList {

			if this_.IsStop {
				return
			}

			value, valueOk := importData[column.Name]
			if !valueOk {
				continue
			}
			valueString, valueStringOk := value.(string)
			if valueStringOk && valueString != "" {
				var scriptValue goja.Value
				scriptValue, err = vm.RunString(valueString)
				if err != nil {
					util.Logger.Error("表达式执行异常", zap.Any("script", valueString), zap.Error(err))
					return
				}
				value = scriptValue.Export()
			}
			data[column.Name] = value

			err = vm.Set(column.Name, value)
			if err != nil {
				return
			}
		}
		this_.ReadyDataCount++
		dataList = append(dataList, data)
		if len(dataList) >= batchNumber {

			if this_.IsStop {
				return
			}
			err = this_.doImportData(database, table, columnList, dataList)
			if err != nil {
				this_.ErrorCount += len(dataList)
				return
			} else {
				this_.SuccessCount += len(dataList)
			}
			dataList = []map[string]interface{}{}
		}
	}
	err = this_.doImportData(database, table, columnList, dataList)
	if err != nil {
		this_.ErrorCount += len(dataList)
		return
	} else {
		this_.SuccessCount += len(dataList)
	}
	return
}

func (this_ *databaseImportTask) doImportData(database, table string, columnList []*db.TableColumnModel, dataList []map[string]interface{}) (err error) {

	if len(dataList) == 0 {
		return
	}
	var sqlList []string
	var paramsList [][]interface{}

	sqlList, paramsList, err = db.DataListInsertSql(this_.generateParam, database, table, columnList, dataList)
	if err != nil {
		return
	}

	_, err = this_.service.Execs(sqlList, paramsList)
	if err != nil {
		return
	}
	return
}
