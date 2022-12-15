package elasticsearch

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"teamide/pkg/data_engine"
	"teamide/pkg/util"
	"time"
)

type ImportTask struct {
	*Task
	ImportType string `json:"importType,omitempty"`

	Id         string                `json:"id,omitempty"`
	ColumnList []*StrategyDataColumn `json:"columnList"`
}

type StrategyDataColumn struct {
	Name        string `json:"name,omitempty"`
	Value       string `json:"value,omitempty"`
	ReuseNumber int    `json:"reuseNumber,omitempty"`
}

func (this_ *ImportTask) do() (err error) {
	if this_.ImportType == "strategy" {
		err = this_.doStrategy()
		if err != nil {
			return
		}
	}
	return
}
func (this_ *ImportTask) doStrategy() (err error) {

	if this_.IndexName == "" {
		err = errors.New("必须配置indexName")
		return
	}
	if this_.DataNumber <= 0 {
		return
	}
	if this_.needStop() {
		return
	}
	var threadNumber = this_.ThreadNumber
	var dataNumber = this_.DataNumber
	if threadNumber < 1 {
		threadNumber = 1
	}
	var eachThreadDataCount = dataNumber / threadNumber
	if eachThreadDataCount < 1 {
		eachThreadDataCount = 1
	}

	var batchNumber = this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 200
	}

	var waitGroupForStop = &sync.WaitGroup{}

	waitGroupForStop.Add(threadNumber)

	var startIndex int
	for threadIndex := 0; threadIndex < threadNumber; threadIndex++ {
		var threadDataCount = eachThreadDataCount
		if dataNumber < threadDataCount {
			threadDataCount = dataNumber
		}
		go this_.doThreadStrategy(waitGroupForStop, threadIndex, startIndex, threadDataCount, batchNumber)
		dataNumber = dataNumber - threadDataCount
		startIndex = startIndex + threadDataCount
	}

	waitGroupForStop.Wait()

	return
}

func (this_ *ImportTask) doThreadStrategy(waitGroupForStop *sync.WaitGroup, threadIndex int, startIndex int, threadDataCount int, batchNumber int) {

	var err error

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			util.Logger.Error("doThreadStrategy error", zap.Error(err))
		}
		waitGroupForStop.Done()
	}()
	if threadDataCount <= 0 {
		return
	}
	if this_.needStop() {
		return
	}

	task := &data_engine.StrategyTask{}

	task.StrategyData = &data_engine.StrategyData{
		DataStatistics: &data_engine.DataStatistics{},
	}

	this_.ReadyDataStatisticsList = append(this_.ReadyDataStatisticsList, task.DataStatistics)
	doDataStatistics := &data_engine.DataStatistics{}
	this_.DoDataStatisticsList = append(this_.DoDataStatisticsList, doDataStatistics)

	task.StrategyData.AddField(&data_engine.StrategyDataField{
		Name:  "id",
		Value: this_.Id,
	})
	for _, column := range this_.ColumnList {
		if column.Name == "" {
			continue
		}
		task.StrategyData.AddField(&data_engine.StrategyDataField{
			Name:        column.Name,
			Value:       column.Value,
			ReuseNumber: column.ReuseNumber,
		})
	}

	task.OnError = func(onErr error) {
		err = onErr
	}

	this_.taskList = append(this_.taskList, task)

	task.StartIndex = startIndex
	task.DataNumber = threadDataCount
	task.BatchNumber = batchNumber

	task.OnDataList = func(dataList []map[string]interface{}) (err error) {

		if this_.needStop() {
			return
		}
		var startTime = time.Now()

		var docs []*InsertDoc
		for _, data := range dataList {

			doc := &InsertDoc{
				IndexName: this_.IndexName,
			}
			doc.Id, err = util.GetStringValue(data["id"])
			if err != nil {
				doDataStatistics.IncrDataErrorCount(1, 0)
				util.Logger.Error("get id value error", zap.Any("data", data), zap.Error(err))
				if this_.ErrorContinue {
					err = nil
					continue
				}
				return
			}

			doc.Doc = data

			docs = append(docs, doc)
		}

		var size = len(docs)

		_, err = this_.Service.BatchInsertNotWait(docs)
		var endTime = time.Now()
		var useTime = util.GetTimeTime(endTime) - util.GetTimeTime(startTime)
		if err != nil {
			doDataStatistics.IncrDataErrorCount(size, useTime)
			util.Logger.Error("BatchInsertNotWait error", zap.Any("size", size), zap.Any("useTime", useTime), zap.Error(err))
			if this_.ErrorContinue {
				err = nil
			}
			return
		}
		doDataStatistics.IncrDataSuccessCount(size, useTime)

		return
	}

	task.OnEnd = func() {

	}

	task.Start()

	return
}
