package elasticsearch

import (
	"errors"
	"teamide/pkg/data_engine"
	"teamide/pkg/util"
	"time"
)

type ImportTask struct {
	*Task
	ImportType  string                `json:"importType,omitempty"`
	Count       int                   `json:"count,omitempty"`
	BatchNumber int                   `json:"batchNumber,omitempty"`
	Id          string                `json:"id,omitempty"`
	ColumnList  []*StrategyDataColumn `json:"columnList"`
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
	if this_.Count <= 0 {
		return
	}
	if this_.needStop() {
		return
	}

	task := &data_engine.StrategyTask{}

	this_.StrategyData = &data_engine.StrategyData{}

	task.StrategyData = this_.StrategyData

	this_.StrategyData.AddField(&data_engine.StrategyDataField{
		Name:  "id",
		Value: this_.Id,
	})
	for _, column := range this_.ColumnList {
		if column.Name == "" {
			continue
		}
		this_.StrategyData.AddField(&data_engine.StrategyDataField{
			Name:        column.Name,
			Value:       column.Value,
			ReuseNumber: column.ReuseNumber,
		})
	}

	task.OnError = func(onErr error) {
		err = onErr
	}

	this_.taskList = append(this_.taskList, task)

	var batchNumber = this_.BatchNumber
	if batchNumber <= 0 {
		batchNumber = 200
	}
	this_.StrategyData.DataNumber = this_.Count
	this_.StrategyData.BatchNumber = batchNumber

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
				this_.IncrDataErrorCount(1, 0)
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
			this_.IncrDataErrorCount(size, useTime)
			return
		}
		this_.IncrDataSuccessCount(size, useTime)

		return
	}

	task.OnEnd = func() {

	}

	task.Start()

	return
}
