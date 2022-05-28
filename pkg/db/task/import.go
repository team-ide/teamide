package task

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"teamide/pkg/data"
	"teamide/pkg/db"
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
	Count       int               `json:"count,omitempty"`
	BatchNumber int               `json:"batchNumber,omitempty"`
	ColumnList  []*StrategyColumn `json:"columnList,omitempty"`
}
type StrategyColumn struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type ImportTask struct {
	Database         string                 `json:"database,omitempty"`
	Table            string                 `json:"table,omitempty"`
	ColumnList       []*db.TableColumnModel `json:"columnList,omitempty"`
	Key              string                 `json:"key,omitempty"`
	ImportType       string                 `json:"importType,omitempty"`
	StrategyDataList []*StrategyData        `json:"strategyDataList,omitempty"`
	DataCount        int                    `json:"dataCount"`
	ReadyDataCount   int                    `json:"readyDataCount"`
	SuccessCount     int                    `json:"successCount"`
	ErrorCount       int                    `json:"errorCount"`
	IsEnd            bool                   `json:"isEnd,omitempty"`
	StartTime        time.Time              `json:"startTime,omitempty"`
	EndTime          time.Time              `json:"endTime,omitempty"`
	Error            string                 `json:"error,omitempty"`
	UseTime          int64                  `json:"useTime"`
	IsStop           bool                   `json:"isStop"`
	GenerateParam    *db.GenerateParam      `json:"-"`
	Service          *db.Service            `json:"-"`
	taskList         []*data.StrategyTask
}

func (this_ *ImportTask) Stop() {
	this_.IsStop = true
	for _, t := range this_.taskList {
		t.Stop()
	}
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
		if this_.needStop() {
			break
		}
		err = this_.doStrategyData(this_.Database, this_.Table, this_.ColumnList, strategyData)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *ImportTask) doStrategyData(database, table string, columnList []*db.TableColumnModel, strategyData *StrategyData) (err error) {
	if strategyData.Count <= 0 {
		return
	}
	if this_.needStop() {
		return
	}

	batchNumber := strategyData.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 100
	}

	task := &data.StrategyTask{}

	taskStrategyData := &data.StrategyData{}

	task.StrategyDataList = append(task.StrategyDataList, taskStrategyData)

	taskStrategyData.Count = strategyData.Count
	for _, strategyColumn := range strategyData.ColumnList {
		taskStrategyData.FieldList = append(taskStrategyData.FieldList, &data.StrategyDataField{
			Name:  strategyColumn.Name,
			Value: strategyColumn.Value,
		})
	}

	task.OnError = func(onErr error) {
		err = onErr
	}

	this_.taskList = append(this_.taskList, task)

	var dataList []map[string]interface{}
	task.OnData = func(onData map[string]interface{}) (err error) {

		if this_.needStop() {
			return
		}
		this_.ReadyDataCount++

		dataList = append(dataList, onData)
		if len(dataList) >= batchNumber {

			err = this_.doImportData(database, table, columnList, dataList)
			if err != nil {
				return
			}
			dataList = []map[string]interface{}{}
		}
		return
	}

	task.OnEnd = func() {

	}

	task.Start()

	if len(dataList) > 0 {
		err = this_.doImportData(database, table, columnList, dataList)
		if err != nil {
			return
		}
	}

	return
}

func (this_ *ImportTask) doImportData(database, table string, columnList []*db.TableColumnModel, dataList []map[string]interface{}) (err error) {

	if len(dataList) == 0 {
		return
	}

	insertColumns := ""
	for _, column := range columnList {
		insertColumns += this_.GenerateParam.PackingCharacterColumn(column.Name) + ", "
	}

	insertColumns = strings.TrimSuffix(insertColumns, ", ")

	sql := "INSERT INTO "

	if this_.GenerateParam.AppendDatabase && database != "" {
		sql += this_.GenerateParam.PackingCharacterDatabase(database) + "."
	}
	sql += this_.GenerateParam.PackingCharacterTable(table)
	if insertColumns != "" {
		sql += "(" + insertColumns + ")"
	}

	sql += "VALUES"
	var values []interface{}

	for _, data := range dataList {

		insertValues := ""
		for _, column := range columnList {
			value, valueOk := data[column.Name]
			if !valueOk {
				insertValues += "NULL, "
			} else {
				insertValues += "?, "
				values = append(values, this_.GenerateParam.FormatColumnValue(column, value))
			}

		}
		insertValues = strings.TrimSuffix(insertValues, ", ")

		sql += "(" + insertValues + "), "

	}
	sql = strings.TrimSuffix(sql, ", ")

	_, err = this_.Service.Execs([]string{sql}, [][]interface{}{values})
	if err != nil {
		this_.ErrorCount += len(dataList)
		return
	}
	this_.SuccessCount += len(dataList)
	return
}

func (this_ *ImportTask) needStop() bool {
	if this_.IsStop || this_.IsEnd {
		return true
	}
	return false
}
